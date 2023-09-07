package main

import (
	"context"
	"douyin/dal/db"
	favorite "douyin/kitex_gen/favorite"
	"douyin/middleware/minio"
	"gorm.io/gorm"
	"log"
)

// FavoriteSrvImpl implements the last service interface defined in the IDL.
type FavoriteSrvImpl struct{}

func Authenticate(token string, id int64) bool {
	claims, err := JwtParser.ParseToken(token)
	if err != nil {
		return false
	}
	return claims.Id == id
}

// FavoriteAction implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	resp = new(favorite.DouyinFavoriteActionResponse)
	claims, err := JwtParser.ParseToken(req.Token)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "认证失败"
		return resp, err
	}
	action := req.ActionType
	tx := db.DB.Begin()
	if action == 1 {
		err := db.InsertFavoriteByVideoId(tx, claims.Id, req.VideoId)
		if err != nil {
			// 回滚事务
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.IncreaseFavoritingCountByVideoId(tx, req.VideoId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.IncreaseFavoriteCountByUserId(tx, claims.Id)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
	} else {
		err := db.DeleteFavoriteByVideoId(tx, claims.Id, req.VideoId)
		if err != nil {
			// 回滚事务
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.DecreaseFavoritingCountByVideoId(tx, req.VideoId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.DecreaseFavoriteCountByUserId(tx, claims.Id)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
	}
	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		// 处理提交事务错误
		resp.StatusCode = -1
		resp.StatusMsg = "操作失败"
		tx.Rollback()
		return resp, err
	}
	resp.StatusCode = 0
	resp.StatusMsg = "操作成功"
	return resp, err
}

// FavoriteList implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (resp *favorite.DouyinFavoriteListResponse, err error) {
	// TODO: Your code here...
	resp = new(favorite.DouyinFavoriteListResponse)
	if !Authenticate(req.Token, req.UserId) {
		resp.StatusCode = -1
		resp.StatusMsg = "认证失败"
		return resp, err
	}
	FavoriteList, err := db.GetFavoriteListByVideoId(db.DB, req.UserId)
	if err != nil {
		return &favorite.DouyinFavoriteListResponse{StatusCode: -1}, err
	}
	var FavoriteVideos []*favorite.Video

	for _, f := range FavoriteList {
		if err != nil {
			return &favorite.DouyinFavoriteListResponse{StatusCode: -1}, err
		}
		author, err := db.GetUserById(f.AuthorID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("用户不存在")
				res := &favorite.DouyinFavoriteListResponse{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				}
				return res, nil
			} else {
				log.Printf(err.Error())
				res := &favorite.DouyinFavoriteListResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误，获取喜欢列表失败",
				}
				return res, nil
			}
		}
		is_follow := true
		if _, err := db.IsFollowById(req.UserId, author.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				is_follow = false
			} else {
				log.Printf(err.Error())
				res := &favorite.DouyinFavoriteListResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误",
				}
				return res, nil
			}
		}
		playUrl, err := minio.GetFileUrl(minio.VideoBucketName, f.PlayUrl)
		if err != nil {
			log.Printf("minio服务器错误，获取视频url失败")
			res := &favorite.DouyinFavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "minio服务器错误，获取视频url失败",
			}
			return res, nil
		}
		coverUrl, err := minio.GetFileUrl(minio.VideoBucketName, f.CoverUrl)
		if err != nil {
			log.Printf("minio服务器错误，获取封面url失败")
			res := &favorite.DouyinFavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "minio服务器错误，获取封面url失败",
			}
			return res, nil
		}
		FavoriteVideo := favorite.Video{
			Id: f.ID,
			Author: &favorite.User{
				Id:              author.ID,
				Name:            author.UserName,
				FollowCount:     author.FollowCount,
				FollowerCount:   author.FollowerCount,
				IsFollow:        is_follow,
				Avatar:          author.Avatar,
				BackgroundImage: author.BackgroundImage,
				Signature:       author.Signature,
				TotalFavorited:  author.TotalFavorited,
				WorkCount:       author.WorkCount,
				FavoriteCount:   author.FavoriteCount,
			},
			PlayUrl:       playUrl.String(),
			CoverUrl:      coverUrl.String(),
			FavoriteCount: f.FavoriteCount,
			CommentCount:  f.CommentCount,
			IsFavorite:    true,
			Title:         f.Title,
		}
		FavoriteVideos = append(FavoriteVideos, &FavoriteVideo)
	}
	resp = &favorite.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "获取成功",
		VideoList:  FavoriteVideos,
	}
	return resp, nil
}
