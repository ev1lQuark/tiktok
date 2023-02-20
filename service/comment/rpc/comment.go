package main

import (
	"flag"
	"fmt"

	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/config"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/server"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/svc"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/types/comment"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/comment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		comment.RegisterLikeServer(grpcServer, server.NewLikeServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
