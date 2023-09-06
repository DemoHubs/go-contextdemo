package ctxdemo

import (
	v1 "contextdemo/api/hello/v1"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/context"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) demo(ctx context.Context, req *v1.Req) (res *v1.Res, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
