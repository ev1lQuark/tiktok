package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/ev1lQuark/tiktok/service/user/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/handler"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, interface{}) {
		return http.StatusInternalServerError, err.Error()
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
