package rpc

import (
	"context"
	"douyin/kitex_gen/user"
	"douyin/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var userClient userservice.Client

func initUser() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.5.54:2379"})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
		"userservice",
		client.WithTransportProtocol(transport.GRPC),
		client.WithResolver(r),

		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (*user.DouyinUserRegisterResponse, error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}

func Login(ctx context.Context, req *user.DouyinUserLoginRequest) (*user.DouyinUserLoginResponse, error) {
	resp, err := userClient.Login(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}

func UserInfo(ctx context.Context, req *user.DouyinUserRequest) (*user.DouyinUserResponse, error) {
	resp, err := userClient.UserInfo(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}
