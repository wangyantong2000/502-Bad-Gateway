package main

import (
	"douyin/dal/db"
	"douyin/kitex_gen/favorite/favoritesrv"
	"douyin/middleware/jwt"
	"douyin/middleware/minio"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

var (
	JwtParser *jwt.JWT
)

func main() {
	db.Init()
	JwtParser = jwt.NewJWT([]byte("signingKey"))
	minio.Init()
	r, err := etcd.NewEtcdRegistry([]string{"192.168.5.54:2379"}) // 服务器地址:2379
	if err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8883")
	svr := favoritesrv.NewServer(new(FavoriteSrvImpl),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "favoritesrv"}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
