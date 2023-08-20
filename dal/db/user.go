package db

import (
	"douyin/dal/model"
	"log"
)

func GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("user_name=?", name).First(&user).Error
	if err != nil {
		log.Printf("读取失败")
		return nil, err
	}
	return user, nil
}
func GetUserById(id int64) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("id=?", id).First(&user).Error
	if err != nil {
		log.Printf("读取失败")
		return user, err
	}
	return user, nil
}
func CreateUser(user *model.User) error {
	if err := DB.Create(&user).Error; err != nil {
		log.Printf("创建失败")
		return err
	}
	return nil
}
func IsFollowById(user_id int64, to_user_id int64) (*model.Follow, error) {
	follow := new(model.Follow)
	err := DB.Where("user_id=? AND to_user_id=?", user_id, to_user_id).First(&follow).Error
	if err != nil {
		log.Printf("读取失败")
		return follow, err
	}
	return follow, nil
}
