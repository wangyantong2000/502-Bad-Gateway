package rpc

import (
	"context"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var relationClient relationservice.Client

func InitClient() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.5.54:2379"})
	if err != nil {
		panic(err)
	}
	c, err := relationservice.NewClient(
		"relationservice",
		client.WithTransportProtocol(transport.GRPC),
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	relationClient = c
}

func ChangeRelation(ctx context.Context, req *relation.DouyinRelationActionRequest) (*relation.DouyinRelationActionResponse, error) {
	resp, err := relationClient.ChangeRelation(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func GetFriendsList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (*relation.DouyinRelationFriendListResponse, error) {
	resp, err := relationClient.GetFriendList(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func GetFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (*relation.DouyinRelationFollowerListResponse, error) {
	resp, err := relationClient.GetFollowerList(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}

func GetFollowingList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (*relation.DouyinRelationFollowListResponse, error) {
	resp, err := relationClient.GetFollowingList(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}
