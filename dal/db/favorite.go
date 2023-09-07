package db

import (
	"douyin/dal/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

func DeleteFavoriteByVideoId(tx *gorm.DB, fromId int64, toId int64) error {
	var record model.Favorite
	err := tx.Debug().Model(&model.Favorite{}).Where("user_id=? AND video_id=?", fromId, toId).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("取消点赞失败")
		return err
	}
	err = tx.Delete(&record).Error // 取消点赞
	if err != nil {
		log.Printf("取消点赞失败")
		return err
	}
	return err
}

func InsertFavoriteByVideoId(tx *gorm.DB, fromId int64, videoId int64) error {
	var record model.Favorite
	record = model.Favorite{UserID: fromId, VideoID: videoId}
	err := tx.Debug().Create(&record).Error
	if err != nil {
		log.Printf("点赞失败")
		return err
	}
	return err

}

func GetFavoritingVideoById(tx *gorm.DB, id int64) ([]*model.Video, error) {
	FavoritingIds, err := GetFavoritingIdsByVideoId(tx, id)
	if err != nil {
		log.Printf("读取失败")
		return nil, err
	}
	if FavoritingIds == nil {
		return nil, err
	}
	var videos []*model.Video
	err = tx.Debug().Where("author_id=?)", FavoritingIds).Find(&videos).Error
	if err != nil {
		log.Printf("查询失败")
		return nil, err
	}
	return videos, nil
}

func GetFavoritingIdsByVideoId(tx *gorm.DB, id int64) ([]int64, error) {
	var FavoritingIds []int64
	err := tx.Model(&model.Favorite{}).Debug().Select("to_video_id").Where("video_id=?", id).Find(&FavoritingIds).Error
	return FavoritingIds, err
}

func GetFavoriteListById(tx *gorm.DB, id int64) ([]*model.Video, error) {
	var FavoritingIds []int64
	err := tx.Model(&model.Favorite{}).Debug().Select("user_id").Where("video_id=?", id).Find(&FavoritingIds).Error
	if err != nil {
		log.Printf("读取失败")
		return nil, err
	}
	if FavoritingIds == nil {
		return nil, err
	}
	var videos []*model.Video
	err = tx.Debug().Where("id IN (?)", FavoritingIds).Find(&videos).Error
	if err != nil {
		log.Printf("查询视频失败")
		return nil, err
	}
	return videos, nil
}

func GetFavoriteListByVideoId(tx *gorm.DB, id int64) ([]*model.Video, error) {
	var Favoritelist []*model.Favorite
	err := tx.Debug().Raw("SELECT user_id,video_id FROM favorites WHERE user_id=? AND video_id IN (SELECT user_id FROM favoritelist f where f.video_id = favorite.video_id)", id).Scan(&Favoritelist).Error
	if err != nil {
		log.Printf("查询喜欢列表失败")
		return nil, err
	}
	var favoriteIDs []int64
	for _, Favorite := range Favoritelist {
		favoriteIDs = append(favoriteIDs, Favorite.VideoID)
	}
	var Favorites []*model.Video
	err = tx.Debug().Where("id IN (?)", favoriteIDs).Find(&Favorites).Error
	if err != nil {
		log.Printf("查询视频信息失败")
		return nil, err
	}
	return Favorites, nil
}

func IncreaseFavoriteCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID
	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Increment the FollowerCount field
	user.FavoriteCount++

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func IncreaseFavoritingCountByVideoId(tx *gorm.DB, videoId int64) error {
	// Retrieve the user by their ID
	var video model.Video
	if err := tx.First(&video, videoId).Error; err != nil {
		return err
	}

	// Increase the FollowCount field
	video.FavoriteCount++

	// Update the user in the database
	if err := tx.Save(&video).Error; err != nil {
		return err
	}

	return nil
}

// DecreaseFollowingCountByUserId decreases the FollowCount of a user by their ID.
func DecreaseFavoritingCountByVideoId(tx *gorm.DB, videoId int64) error {
	// Retrieve the user by their ID
	var video model.Video
	if err := tx.First(&video, videoId).Error; err != nil {
		return err
	}

	// Decrease the FollowCount field
	if video.FavoriteCount > 0 {
		video.FavoriteCount--
	}

	// Update the user in the database
	if err := tx.Save(&video).Error; err != nil {
		return err
	}

	return nil
}

func DecreaseFavoriteCountByUserId(tx *gorm.DB, userId int64) error {
	// Retrieve the user by their ID

	var user model.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// Decrease the FollowCount field
	if user.FavoriteCount > 0 {
		user.FavoriteCount--
	}

	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// Favorite new favorite data.
/*func Favorite(ctx context.Context, uid int64, vid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1. 新增点赞数据
		user := new(model.User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}

		video := new(model.Video)
		if err := tx.WithContext(ctx).First(video, vid).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&user).Association("FavoriteVideos").Append(video); err != nil {
			return err
		}
		//2.改变 video 表中的 favorite count
		res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// DisFavorite deletes the specified favorite from the database
func DisFavorite(ctx context.Context, uid int64, vid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1. 删除点赞数据
		user := new(model.User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}

		video, err := GetFavoriteRelation(ctx, uid, vid)
		if err != nil {
			return err
		}

		err = tx.Unscoped().WithContext(ctx).Model(&user).Association("FavoriteVideos").Delete(video)
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// FavoriteList returns a list of Favorite videos.
func FavoriteList(ctx context.Context, uid int64) ([]Video, error) {
	user := new(model.User)
	if err := DB.WithContext(ctx).First(user, uid).Error; err != nil {
		return nil, err
	}

	videos := []*model.Video{}
	// if err := DB.WithContext(ctx).First(&video, vid).Error; err != nil {
	// 	return nil, err
	// }

	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&videos); err != nil {
		return nil, err
	}
	return videos, nil
}*/
