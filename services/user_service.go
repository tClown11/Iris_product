package services

import (
	"Iris_product/datamodels"
	"Iris_product/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOK bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOK bool) {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}

	isOK, _ = Validatepassword(pwd, user.HashPassword)
	if !isOK {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

// 对用户的明文密码进行加密
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// 对比用户输入的密码和数据库内存的密码进行比对(byte类型的比较)
func Validatepassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}

	return true, nil
}
