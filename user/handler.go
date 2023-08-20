package main

import (
	"context"
	"crypto/md5"
	"douyin/dal/db"
	"douyin/dal/model"
	user "douyin/kitex_gen/user"
	"douyin/middleware/jwt"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	name := req.Username
	password := req.Password
	if len(name) == 0 {
		log.Printf("名字为空，注册失败")
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "名字为空，注册失败",
		}
		return res, nil
	}
	if len(password) == 0 {
		log.Printf("密码为空，注册失败")
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "密码为空，注册失败",
		}
		return res, nil
	}
	if len(name) > 32 {
		log.Printf("名字太长超过32字符，注册失败")
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "名字太长超过32字符，注册失败",
		}
		return res, nil
	}
	if len(password) > 32 {
		log.Printf("密码太长超过32字符，注册失败")
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "密码太长超过32字符，注册失败",
		}
		return res, nil
	}
	u, err := db.GetUserByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("用户名不重复，可以注册")
		} else {
			log.Printf(err.Error())
			res := &user.DouyinUserRegisterResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误，注册失败",
			}
			return res, nil
		}
	}
	if u != nil {
		log.Printf("名字重复，注册失败")
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "名字重复，注册失败",
		}
		return res, nil
	}
	h := md5.New()
	if _, err = io.WriteString(h, password); err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，注册失败",
		}
		return res, nil
	}
	password = fmt.Sprintf("%x", h.Sum(nil))
	u = new(model.User)
	u.UserName = name
	u.Password = password

	if err := db.CreateUser(u); err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，注册失败",
		}
		return res, nil
	}
	claims := jwt.CustomClaims{Id: u.ID}
	claims.ExpiresAt = time.Now().Add(time.Minute * 5).Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：token 创建失败",
		}
		return res, nil
	}
	log.Printf(string(u.ID), token)
	res := &user.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserId:     u.ID,
		Token:      token,
	}
	return res, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	name := req.Username
	password := req.Password
	if len(name) == 0 {
		log.Printf("名字为空，登录失败")
		res := &user.DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "名字为空，登录失败",
		}
		return res, nil
	}
	if len(password) == 0 {
		log.Printf("密码为空，登录失败")
		res := &user.DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "密码为空，登录失败",
		}
		return res, nil
	}
	u, err := db.GetUserByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("用户名不存在，登录失败")
			res := &user.DouyinUserLoginResponse{
				StatusCode: -1,
				StatusMsg:  "用户名不存在，登录失败",
			}
			return res, nil
		} else {
			log.Printf(err.Error())
			res := &user.DouyinUserLoginResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误，登录失败",
			}
			return res, nil
		}
	}

	h := md5.New()
	if _, err := io.WriteString(h, password); err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，登录失败",
		}
		return res, nil
	}
	password = fmt.Sprintf("%x", h.Sum(nil))
	if u.Password != password {
		log.Printf("用户名与密码不匹配，登录失败")
		res := &user.DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名与密码不匹配，登录失败",
		}
		return res, nil
	}
	claims := jwt.CustomClaims{Id: u.ID}
	claims.ExpiresAt = time.Now().Add(time.Hour * 12).Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：token 创建失败",
		}
		return res, nil
	}
	res := &user.DouyinUserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     u.ID,
		Token:      token,
	}
	return res, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	userId := req.UserId
	token := req.Token
	claims, err := Jwt.ParseToken(token)
	if err != nil {
		log.Printf(err.Error())
		res := &user.DouyinUserResponse{
			StatusCode: -1,
			StatusMsg:  "token 解析失败",
		}
		return res, nil
	}
	myId := claims.Id
	u, err := db.GetUserById(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("用户不存在")
			res := &user.DouyinUserResponse{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			}
			return res, nil
		} else {
			log.Printf(err.Error())
			res := &user.DouyinUserResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}
			return res, nil
		}
	}
	is_follow := true
	if _, err := db.IsFollowById(myId, u.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			is_follow = false
		} else {
			log.Printf(err.Error())
			res := &user.DouyinUserResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}
			return res, nil
		}
	}
	res := &user.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户信息成功",
		User: &user.User{
			Id:              u.ID,
			Name:            u.UserName,
			FollowCount:     u.FollowCount,
			FollowerCount:   u.FollowerCount,
			IsFollow:        is_follow,
			Avatar:          u.Avatar,
			BackgroundImage: u.BackgroundImage,
			Signature:       u.Signature,
			TotalFavorited:  u.TotalFavorited,
			WorkCount:       u.WorkCount,
			FavoriteCount:   u.FavoriteCount,
		},
	}
	return res, nil
}
