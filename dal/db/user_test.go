package db

import (
	"douyin/dal/model"
	"log"
	"testing"
)

func TestGetUserByName(t *testing.T) {
	Init()
	res, err := GetUserByName("123456889")
	if err == nil {
		log.Println("成功", res)
	}
}
func TestCreateUser(t *testing.T) {
	Init()
	user := new(model.User)
	user.UserName = "123456"
	user.Password = "123456"
	err := CreateUser(user)
	if err == nil {
		log.Println("成功")
	}
}
func TestGetUserById(t *testing.T) {
	Init()
	res, err := GetUserById(1)
	if err == nil {
		log.Println("成功", res)
	}
}
