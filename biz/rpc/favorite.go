package rpc

import (
	"context"
	"douyin/kitex_gen/favorite"
	"douyin/kitex_gen/favorite/favoritesrv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var favoriteClient favoritesrv.Client

func initFavorite() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.5.54:2379"})
	if err != nil {
		panic(err)
	}
	c, err := favoritesrv.NewClient(
		"favoritesrv",
		client.WithTransportProtocol(transport.GRPC),
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	favoriteClient = c
}
func FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (res *favorite.DouyinFavoriteActionResponse, err error) {
	resp, err := favoriteClient.FavoriteAction(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}

func FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (res *favorite.DouyinFavoriteListResponse, err error) {
	resp, err := favoriteClient.FavoriteList(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}
