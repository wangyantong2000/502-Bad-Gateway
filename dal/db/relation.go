package db

import (
	"douyin/dal/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

//func ChangeRelation(fromId int64, toId int64, action int64) error {
//	if action == 1 {
//		err := InsertRelationByUserId(fromId, toId)
//		if err != nil {
//			return err
//		}
//	} else {
//		err := DeleteRelationByUserId(fromId, toId)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func DeleteRelationByUserId(tx *gorm.DB, fromId int64, toId int64) error {
	var record model.Follow
	err := tx.Debug().Model(&model.Follow{}).Where("user_id=? AND to_user_id=?", fromId, toId).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("取消关注失败")
		return err
	}
	err = tx.Delete(&record).Error // 取消关注
	if err != nil {
		log.Printf("取消关注失败")
		return err
	}
	return err
}

func InsertRelationByUserId(tx *gorm.DB, fromId int64, toId int64) error {
	var record model.Follow
	record = model.Follow{UserID: fromId, ToUserID: toId}
	err := tx.Debug().Create(&record).Error
	if err != nil {
		log.Printf("关注失败")
		return err
	}
	return err

}

func GetFollowingUserById(tx *gorm.DB, id int64) ([]*model.User, error) {
	followingIds, err := GetFollowingIdsByUserId(tx, id)
	if err != nil {
		log.Printf("读取失败")
		return nil, err
	}
	if followingIds == nil {
		return nil, err
	}
	var users []*model.User
	err = tx.Debug().Where("id IN (?)", followingIds).Find(&users).Error
	if err != nil {
		log.Printf("查询失败")
		return nil, err
	}
	return users, nil
}

func GetFollowingIdsByUserId(tx *gorm.DB, id int64) ([]int64, error) {
	var followingIds []int64
	err := tx.Model(&model.Follow{}).Debug().Select("to_user_id").Where("user_id=?", id).Find(&followingIds).Error
	return followingIds, err
}

func GetFollowerListById(tx *gorm.DB, id int64) ([]*model.User, error) {
	var followingIds []int64
	err := tx.Model(&model.Follow{}).Debug().Select("user_id").Where("to_user_id=?", id).Find(&followingIds).Error
	if err != nil {
		log.Printf("读取失败")
		return nil, err
	}
	if followingIds == nil {
		return nil, err
	}
	var users []*model.User
	err = tx.Debug().Where("id IN (?)", followingIds).Find(&users).Error
	if err != nil {
		log.Printf("查询用户失败")
		return nil, err
	}
	return users, nil
}

func GetFriendListByUserId(tx *gorm.DB, id int64) ([]*model.User, error) {
	var friendList []*model.Follow
	err := tx.Debug().Raw("SELECT user_id,to_user_id FROM follows WHERE user_id=? AND to_user_id IN (SELECT user_id FROM follows f where f.to_user_id = follows.user_id)", id).Scan(&friendList).Error
	if err != nil {
		log.Printf("查询好友列表失败")
		return nil, err
	}
	var friendIDs []int64
	for _, follow := range friendList {
		friendIDs = append(friendIDs, follow.ToUserID)
	}
	var friends []*model.User
	// 查询与好友关系匹配的用户
	err = tx.Debug().Where("id IN (?)", friendIDs).Find(&friends).Error
	if err != nil {
		log.Printf("查询好友信息失败")
		return nil, err
	}
	return friends, nil
}

func GetLatestMessageByUserId(tx *gorm.DB, friendId, userId int64) (string, int64, error) {
	var message model.Message
	result := tx.Where("user_id = ? AND to_user_id = ? OR user_id = ? AND to_user_id = ?", userId, friendId, friendId, userId).
		Order("created_date DESC").First(&message)
	if result.Error != nil {
		return "", 0, result.Error
	}

	msgType := int64(0)
	if message.UserID == userId {
		msgType = 1
	}

	return message.Content, msgType, nil
}

func IncreaseFollowerCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Increment the FollowerCount field
	user.FollowerCount++

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func IncreaseFollowingCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Increase the FollowCount field
	user.FollowCount++

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// DecreaseFollowingCountByUserId decreases the FollowCount of a user by their ID.
func DecreaseFollowingCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Decrease the FollowCount field
	if user.FollowCount > 0 {
		user.FollowCount--
	}

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func DecreaseFollowerCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Decrease the FollowCount field
	if user.FollowerCount > 0 {
		user.FollowerCount--
	}

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
