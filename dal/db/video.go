package db

import (
	"douyin/dal/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

func CreateVideo(video *model.Video) error {
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
	if err := tx.Create(&video).Error; err != nil {
		log.Printf(err.Error())
		tx.Rollback()
		return err
	}
	res := tx.Model(&model.User{}).Where("id=?", video.AuthorID).Update("work_count", gorm.Expr("work_count+?", 1))
	if res.Error != nil {
		log.Printf(res.Error.Error())
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected != 1 {
		log.Printf("错误")
		tx.Rollback()
		return errors.New("错误")
	}
	return tx.Commit().Error
}
func GetVideosByAuthorID(userId int64) ([]*model.Video, error) {
	var videos []*model.Video
	if err := DB.Where("author_id=?", userId).Find(&videos).Error; err != nil {
		log.Printf("读取失败")
		return videos, err
	}
	return videos, nil
}
func DeleteVideoByID(id int64, authorId int64) error {
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
	if err := tx.Delete(&model.Video{}, id).Error; err != nil {
		log.Printf(err.Error())
		tx.Rollback()
		return err
	}

	res := tx.Model(&model.User{}).Where("id=?", authorId).Update("work_count", gorm.Expr("work_count-?", 1))
	if res.Error != nil {
		log.Printf(res.Error.Error())
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected != 1 {
		log.Printf("错误")
		tx.Rollback()
		return errors.New("错误")
	}
	return tx.Commit().Error
}
func GetVideosByPubilshTime(limit int, lastTime int64) ([]*model.Video, error) {
	var videos []*model.Video
	if err := DB.Limit(limit).Order("publish_time desc").Where("publish_time < ?", lastTime).Find(&videos).Error; err != nil {
		log.Printf("读取失败")
		return videos, err
	}
	return videos, nil
}
func IsFavoriteById(userId int64, videoId int64) (*model.Favorite, error) {
	favorite := new(model.Favorite)
	err := DB.Where("user_id=? AND video_id=?", userId, videoId).First(&favorite).Error
	if err != nil {
		log.Printf("读取失败")
		return favorite, err
	}
	return favorite, nil
}
