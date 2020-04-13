package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

//根据结构体中sql标签将从数据库里查询的数据映射到结构体中并且转换类型
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem() // 获取obj的值，Elem返回接口v包含的值或指针v指向的值
	// NumField()   ⽤来获取结构体字段（成员）数量
	for i := 0; i < objValue.NumField(); i++ {
		//获取map中对应的值    Type()将value类型变回type类型，Filed(i)返回结构体类型的第 i 个字段，
		// 只能是结构体类型调用,如果 i 超过了总字段数，就会 panic
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//获取struct中对应第 i 个字段的名称
		name := objValue.Type().Field(i).Name
		//获取strcut中对应第 i 字段类型
		structFieldType := objValue.Field(i).Type()
		//获取map中变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		//对比struct的类型与val的类型是否一致
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {

			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
