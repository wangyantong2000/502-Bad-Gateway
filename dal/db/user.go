package db

import (
	"douyin/dal/model"
	"gorm.io/gorm"
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

//	func CreateUser(user *model.User) error {
//		if err := DB.Create(&user).Error; err != nil {
//			log.Printf("创建失败")
//			return err
//		}
//		return nil
//	}
func CreateUser(user *model.User) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		log.Printf(err.Error())
		return err
	}
	if err := tx.Create(&user).Error; err != nil {
		log.Printf(err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func IsFollowById(user_id int64, to_user_id int64) (*model.Follow, error) {
	follow := new(model.Follow)
	err := DB.Where("user_id=? AND to_user_id=?", user_id, to_user_id).First(&follow).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("没有关注")
		} else {
			log.Printf("读取失败")
		}
		return follow, err
	}
	return follow, nil
}
