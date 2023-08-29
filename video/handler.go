package main

import (
	"bytes"
	"context"
	"douyin/dal/db"
	"douyin/dal/model"
	video "douyin/kitex_gen/video"
	"douyin/middleware/minio"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.DouyinFeedRequest) (resp *video.DouyinFeedResponse, err error) {
	token := req.Token
	latestTime := time.Now().Unix()
	var userId int64
	if token != "" {
		claims, err := Jwt.ParseToken(token)
		if err != nil {
			log.Printf(err.Error())
			log.Printf(token)
			res := &video.DouyinFeedResponse{
				StatusCode: -1,
				StatusMsg:  "token 解析失败",
			}
			return res, nil
		}
		userId = claims.Id
	} else {
		userId = -1
	}
	videos, err := db.GetVideosByPubilshTime(30, latestTime)
	if err != nil {
		log.Printf(err.Error())
		res := &video.DouyinFeedResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，视频流失败",
		}
		return res, nil
	}
	videoList := make([]*video.Video, 0)
	for _, v := range videos {
		author, err := db.GetUserById(v.AuthorID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("用户不存在")
				res := &video.DouyinFeedResponse{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				}
				return res, nil
			} else {
				log.Printf(err.Error())
				res := &video.DouyinFeedResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误，获取发布列表失败",
				}
				return res, nil
			}
		}
		is_follow := true
		if token == "" {
			is_follow = false
		} else if _, err := db.IsFollowById(userId, author.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				is_follow = false
			} else {
				log.Printf(err.Error())
				res := &video.DouyinFeedResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误",
				}
				return res, nil
			}
		}
		is_favorite := true
		if token == "" {
			is_favorite = false
		} else if _, err := db.IsFavoriteById(userId, v.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				is_favorite = false
			} else {
				log.Printf(err.Error())
				res := &video.DouyinFeedResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误",
				}
				return res, nil
			}
		}
		playUrl, err := minio.GetFileUrl(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			log.Printf("minio服务器错误，获取视频url失败")
			res := &video.DouyinFeedResponse{
				StatusCode: -1,
				StatusMsg:  "minio服务器错误，获取视频url失败",
			}
			return res, nil
		}
		coverUrl, err := minio.GetFileUrl(minio.VideoBucketName, v.CoverUrl)
		if err != nil {
			log.Printf("minio服务器错误，获取封面url失败")
			res := &video.DouyinFeedResponse{
				StatusCode: -1,
				StatusMsg:  "minio服务器错误，获取封面url失败",
			}
			return res, nil
		}
		videoList = append(videoList, &video.Video{
			Id: v.ID,
			Author: &video.User{
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
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    is_favorite,
			Title:         v.Title,
		})

	}

	if len(videos) != 0 {
		latestTime = videos[len(videos)-1].PublishTime
	}
	res := &video.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "视频流获取成功",
		VideoList:  videoList,
		NextTime:   latestTime,
	}
	return res, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.DouyinPublishActionRequest) (resp *video.DouyinPublishActionResponse, err error) {
	token := req.Token
	data := req.Data
	title := req.Title
	claims, err := Jwt.ParseToken(token)
	if err != nil {
		log.Printf(err.Error())
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "token 解析失败",
		}
		return res, nil
	}
	userId := claims.Id
	if len(title) == 0 {
		log.Printf("视频标题长度为0，上传失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频标题长度为0，上传失败",
		}
		return res, nil
	}
	if len(data) == 0 {
		log.Printf("视频大小为0，上传失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频大小为0，上传失败",
		}
		return res, nil
	}
	//u, err := uuid.NewV4()
	//if err != nil {
	//	log.Printf("uuid错误，上传失败")
	//	res := &video.DouyinPublishActionResponse{
	//		StatusCode: -1,
	//		StatusMsg:  "uuid错误，上传失败",
	//	}
	//	return res, nil
	//}
	createTimestamp := time.Now().Unix()
	playReader := bytes.NewReader(data)
	playContentType := "video/mp4"
	//playName := u.String() + "." + "mp4"
	playName := fmt.Sprintf("%d_%s_%d.mp4", userId, title, createTimestamp)
	err = minio.PutFile(minio.VideoBucketName, playName, playReader, int64(len(data)), playContentType)
	if err != nil {
		log.Printf("minio服务器错误，视频上传失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "minio服务器错误，视频上传失败",
		}
		return res, nil
	}
	url, err := minio.GetFileUrl(minio.VideoBucketName, playName)
	if err != nil {
		log.Printf("minio服务器错误，获取视频url失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "minio服务器错误，获取视频url失败",
		}
		return res, nil
	}

	//playUrl := strings.Split(url.String(), "?")[0]
	//u1, err := uuid.NewV4()
	//if err != nil {
	//	log.Printf("uuid错误，上传失败")
	//	res := &video.DouyinPublishActionResponse{
	//		StatusCode: -1,
	//		StatusMsg:  "uuid错误，上传失败",
	//	}
	//	return res, nil
	//}
	//coverName := u1.String() + "." + "jpg"
	coverName := fmt.Sprintf("%d_%s_%d.png", userId, title, createTimestamp)
	coverData, err := GetImageBuffer(url.String())
	if err != nil {
		log.Printf("ffmpeg错误，上传失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "ffmpeg错误，上传失败",
		}
		return res, nil
	}
	//coverReader := bytes.NewReader(coverData)
	coverContentType := "image/png"
	err = minio.PutFile(minio.VideoBucketName, coverName, coverData, int64(coverData.Len()), coverContentType)
	if err != nil {
		log.Printf("minio服务器错误，封面上传失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "minio服务器错误，封面上传失败",
		}
		return res, nil
	}
	url, err = minio.GetFileUrl(minio.VideoBucketName, coverName)
	if err != nil {
		log.Printf("minio服务器错误，获取封面url失败")
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "minio服务器错误，获取封面url失败",
		}
		return res, nil
	}
	//coverUrl := strings.Split(url.String(), "?")[0]
	v := new(model.Video)
	v.AuthorID = userId
	v.PlayUrl = playName
	v.CoverUrl = coverName
	v.FavoriteCount = 0
	v.CommentCount = 0
	v.Title = title
	v.PublishTime = time.Now().Unix()
	err = db.CreateVideo(v)
	if err != nil {
		log.Printf(err.Error())
		res := &video.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，视频发布失败",
		}
		return res, nil
	}
	res := &video.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "视频上传成功",
	}
	return res, nil
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error) {
	token := req.Token
	userId := req.UserId
	claims, err := Jwt.ParseToken(token)
	if err != nil {
		log.Printf(err.Error())
		res := &video.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "token 解析失败",
		}
		return res, nil
	}
	myId := claims.Id
	videos, err := db.GetVideosByAuthorID(userId)
	if err != nil {

		log.Printf(err.Error())
		res := &video.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，获取发布列表失败",
		}
		return res, err
	}
	videoList := make([]*video.Video, 0)
	for _, v := range videos {
		author, err := db.GetUserById(v.AuthorID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("用户不存在")
				res := &video.DouyinPublishListResponse{
					StatusCode: -1,
					StatusMsg:  "用户不存在",
				}
				return res, nil
			} else {
				log.Printf(err.Error())
				res := &video.DouyinPublishListResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误，获取发布列表失败",
				}
				return res, nil
			}
		}
		is_follow := true
		if _, err := db.IsFollowById(myId, author.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				is_follow = false
			} else {
				log.Printf(err.Error())
				res := &video.DouyinPublishListResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误",
				}
				return res, nil
			}
		}
		is_favorite := true
		if _, err := db.IsFavoriteById(myId, v.ID); err != nil {
			if err == gorm.ErrRecordNotFound {
				is_favorite = false
			} else {
				log.Printf(err.Error())
				res := &video.DouyinPublishListResponse{
					StatusCode: -1,
					StatusMsg:  "服务器内部错误",
				}
				return res, nil
			}
		}
		videoList = append(videoList, &video.Video{
			Id: v.ID,
			Author: &video.User{
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
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    is_favorite,
			Title:         v.Title,
		})

	}
	res := &video.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "发布列表获取成功",
		VideoList:  videoList,
	}
	return res, nil
}
