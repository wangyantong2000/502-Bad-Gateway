// Code generated by hertz generator.

package userapi

import (
	"context"
	userapi "douyin/biz/model/userapi"
	"douyin/biz/rpc"
	"douyin/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
)

// Register .
// @router /douyin/user/register/ [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req userapi.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		log.Printf(err.Error())
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	res, err := rpc.Register(ctx, &user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Printf(err.Error())
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := &user.DouyinUserRegisterResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		UserId:     res.UserId,
		Token:      res.Token,
	}
	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /douyin/user/login/ [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req userapi.DouyinUserLoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	res, err := rpc.Login(context.Background(), &user.DouyinUserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Printf(err.Error())
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := &user.DouyinUserRegisterResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		UserId:     res.UserId,
		Token:      res.Token,
	}
	c.JSON(consts.StatusOK, resp)
}

// UserInfo .
// @router /douyin/user/ [GET]
func UserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req userapi.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	res, err := rpc.UserInfo(context.Background(), &user.DouyinUserRequest{
		UserId: req.UserId,
		Token:  req.Token,
	})
	if err != nil {
		log.Printf(err.Error())
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := &user.DouyinUserResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		User:       res.User,
	}
	c.JSON(consts.StatusOK, resp)
}