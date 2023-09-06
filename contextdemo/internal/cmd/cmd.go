package cmd

import (
	"context"
	"contextdemo/internal/controller/ctxdemo"
	"contextdemo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"contextdemo/internal/controller/hello"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware.Ctx)
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.New(),
				)
				group.Bind(ctxdemo.New())
			})
			s.Run()
			return nil
		},
	}
)
