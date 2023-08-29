package rpc

import (
	"context"
	"douyin/kitex_gen/video"
	"douyin/kitex_gen/video/videoservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var videoClient videoservice.Client

func initVideo() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.5.54.:2379"})
	if err != nil {
		panic(err)
	}
	c, err := videoservice.NewClient(
		"videoservice",
		client.WithTransportProtocol(transport.GRPC),
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}
func Feed(ctx context.Context, req *video.DouyinFeedRequest) (*video.DouyinFeedResponse, error) {
	resp, err := videoClient.Feed(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}
func PublishAction(ctx context.Context, req *video.DouyinPublishActionRequest) (*video.DouyinPublishActionResponse, error) {
	resp, err := videoClient.PublishAction(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}
func PublishList(ctx context.Context, req *video.DouyinPublishListRequest) (*video.DouyinPublishListResponse, error) {
	resp, err := videoClient.PublishList(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}
	return resp, nil
}
