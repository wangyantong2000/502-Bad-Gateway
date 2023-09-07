package main

import (
	"context"
	"douyin/dal/db"
	relation "douyin/kitex_gen/relation"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

func Authenticate(token string, id int64) bool {
	claims, err := JwtParser.ParseToken(token)
	if err != nil {
		return false
	}
	return claims.Id == id
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (resp *relation.DouyinRelationFriendListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.DouyinRelationFriendListResponse)
	if !Authenticate(req.Token, req.UserId) {
		resp.StatusCode = -1
		resp.StatusMsg = "认证失败"
		return resp, err
	}
	friendList, err := db.GetFriendListByUserId(db.DB, req.UserId)
	if err != nil {
		return &relation.DouyinRelationFriendListResponse{StatusCode: -1}, err
	}
	var friendUsers []*relation.FriendUser

	for _, friend := range friendList {
		latestMessage, msgType, err := db.GetLatestMessageByUserId(db.DB, friend.ID, req.UserId)
		if err != nil {
			return &relation.DouyinRelationFriendListResponse{StatusCode: -1}, err
		}
		friendUser := relation.FriendUser{
			Id:              friend.ID,
			Name:            friend.UserName,
			FollowCount:     friend.FollowCount,
			FollowerCount:   friend.FollowerCount,
			IsFollow:        true,
			Avatar:          friend.Avatar,
			BackgroundImage: friend.BackgroundImage,
			Signature:       friend.Signature,
			TotalFavorited:  friend.TotalFavorited,
			WorkCount:       friend.WorkCount,
			FavoriteCount:   friend.FavoriteCount,
			Message:         latestMessage,
			MsgType:         msgType,
		}
		friendUsers = append(friendUsers, &friendUser)
	}
	resp = &relation.DouyinRelationFriendListResponse{
		StatusCode: 0,
		StatusMsg:  "获取成功",
		UserList:   friendUsers,
	}
	return resp, nil
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.DouyinRelationFollowerListResponse)
	if !Authenticate(req.Token, req.UserId) {
		resp.StatusCode = -1
		resp.StatusMsg = "认证失败"
		return resp, err
	}
	followers, err := db.GetFollowerListById(db.DB, req.GetUserId())
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "获取粉丝列表失败"
		return resp, err
	}
	followingIds, err := db.GetFollowingIdsByUserId(db.DB, req.UserId)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "获取关注列表失败"
		return resp, err
	}
	var followList []*relation.User
	for _, follower := range followers {
		user := &relation.User{
			Id:              follower.ID,
			Name:            follower.UserName,
			FollowCount:     follower.FollowCount,
			FollowerCount:   follower.FollowerCount,
			Avatar:          follower.Avatar,
			BackgroundImage: follower.BackgroundImage,
			Signature:       follower.Signature,
			TotalFavorited:  follower.TotalFavorited,
			WorkCount:       follower.WorkCount,
			FavoriteCount:   follower.FavoriteCount,
		}
		for _, followingID := range followingIds {
			if followingID == follower.ID {
				user.IsFollow = true
				break
			}
		}
		followList = append(followList, user)
	}
	resp.UserList = followList
	resp.StatusCode = 0
	resp.StatusMsg = "获取粉丝列表成功"
	return resp, nil
}

// GetFollowingList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowingList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.DouyinRelationFollowListResponse)
	following, err := db.GetFollowingUserById(db.DB, req.GetUserId())
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "获取关注列表失败"
		return resp, err
	}
	var followingList []*relation.User
	for _, follower := range following {
		user := &relation.User{
			Id:              follower.ID,
			Name:            follower.UserName,
			FollowCount:     follower.FollowCount,
			FollowerCount:   follower.FollowerCount,
			IsFollow:        true,
			Avatar:          follower.Avatar,
			BackgroundImage: follower.BackgroundImage,
			Signature:       follower.Signature,
			TotalFavorited:  follower.TotalFavorited,
			WorkCount:       follower.WorkCount,
			FavoriteCount:   follower.FavoriteCount,
		}
		followingList = append(followingList, user)
	}
	resp.UserList = followingList
	resp.StatusCode = 0
	resp.StatusMsg = "获取关注列表成功"
	return resp, nil
}

// ChangeRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ChangeRelation(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	// TODO: Your code here...
	resp = new(relation.DouyinRelationActionResponse)
	claims, err := JwtParser.ParseToken(req.Token)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "认证失败"
		return resp, err
	}
	action := req.ActionType
	tx := db.DB.Begin()
	if action == 1 {
		err := db.InsertRelationByUserId(tx, claims.Id, req.ToUserId)
		if err != nil {
			// 回滚事务
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.IncreaseFollowingCountByUserId(tx, claims.Id)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.IncreaseFollowerCountByUserId(tx, req.ToUserId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
	} else {
		err := db.DeleteRelationByUserId(tx, claims.Id, req.ToUserId)
		if err != nil {
			// 回滚事务
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.DecreaseFollowingCountByUserId(tx, claims.Id)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "操作失败"
			tx.Rollback()
			return resp, err
		}
		err = db.DecreaseFollowerCountByUserId(tx, req.ToUserId)
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
