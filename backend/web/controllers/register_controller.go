package controllers

import (
	"github.com/kataras/iris"
)

type IBController interface{
	//每个控制器都执行统一
    InitController(app *iris.Application, controllerName string)
    //注册控制器名称
    RegisterValue() string
    //注册控制器
    RegisterController()
}

