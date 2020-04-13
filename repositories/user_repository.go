package repositories

import (
	"Iris_product/common"
	"Iris_product/datamodels"
	"database/sql"
	"errors"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userID int64, err error)
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table: table, mysqlConn: db}
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}

		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "iris_product_user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}

	sql := "Select * from " + u.table + " where userName=?"
	row, err := u.mysqlConn.Query(sql, userName)
	defer row.Close()
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,password=?"
	//fmt.Println(sql)
	stmt, err := u.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return userID, err
	}

	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return userID, err
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	sql := "Select * from " + u.table + " where ID=" + strconv.FormatInt(userID, 10)
	row, err := u.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在! ")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
