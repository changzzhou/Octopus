package main

import (
	"flag"
	"fmt"
	"net/http"

	"backend/internal/config"
	"backend/internal/errorx"
	"backend/internal/handler"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/workflow-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		if codeErr, ok := err.(*errorx.CodeError); ok {
			return codeErr.Code, &errorx.ErrorResponse{
				Code:    codeErr.Code,
				Message: codeErr.Message,
			}
		}
		return http.StatusInternalServerError, &errorx.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
