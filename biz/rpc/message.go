package rpc

import (
	"context"
	"douyin/kitex_gen/message"
	"douyin/kitex_gen/message/messageservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var messageClient messageservice.Client

func initMessage() {
	r, err := etcd.NewEtcdResolver([]string{"192.168.5.54:2379"})
	if err != nil {
		panic(err)
	}

	c, err := messageservice.NewClient(
		"messageservice",
		client.WithTransportProtocol(transport.GRPC),
		client.WithResolver(r),

		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	messageClient = c
}

func MessageChat(ctx context.Context, req *message.DouyinMessageChatRequest) (*message.DouyinMessageChatResponse, error) {
	resp, err := messageClient.MessageChat(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}
func MessageAction(ctx context.Context, req *message.DouyinRelationActionRequest) (*message.DouyinRelationActionResponse, error) {
	resp, err := messageClient.MessageAction(ctx, req)
	if err != nil {
		panic(err)
		return resp, err
	}

	return resp, nil
}
