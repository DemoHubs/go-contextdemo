package service

import (
	"contextdemo/internal/model"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Middleware 中间件管理服务
var Middleware = middlewareService{}

type middlewareService struct{}

// Ctx 自定义上下文对象
func (s *middlewareService) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	name := login()
	users := readUser()
	if _, ok := users[name]; ok {
		customCtx := model.Context{
			LoginName: name,
			Authority: users[name],
			List:      readProduct(),
		}
		Context.Init(r, &customCtx)

		// 给模板传递上下文对象中的键值对
		r.Assigns(g.Map{
			"user": customCtx.LoginName,
			"auth": customCtx.Authority,
			"list": customCtx.List,
		})
	} else {
		fmt.Println("用户不存在，请重新登录") // 此处是为了演示做简化，正常情况应该在登录验证函数做好处理
	}

	// 执行后续中间件
	r.Middleware.Next()
}
