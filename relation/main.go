package main

import (
	"douyin/dal/db"
	"douyin/kitex_gen/relation/relationservice"
	"douyin/middleware/jwt"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

var (
	JwtParser *jwt.JWT
)

// relation_service的端口为8884
func main() {
	db.Init()
	JwtParser = jwt.NewJWT([]byte("signingKey"))
	r, err := etcd.NewEtcdRegistry([]string{"192.168.5.54:2379"}) // 服务器地址:2379
	if err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8884")
	svr := relationservice.NewServer(new(RelationServiceImpl),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "relationservice"}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
