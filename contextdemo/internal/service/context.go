package service

import (
	"context"
	"contextdemo/internal/model"
	"errors"

	"io/ioutil"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Context 上下文管理服务
var Context = contextService{}

type contextService struct{}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *contextService) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(model.ContextKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func (s *contextService) Get(ctx context.Context) *model.Context {
	value := ctx.Value(model.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// SetUser 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *contextService) SetUser(ctx context.Context, name string) error {
	save := s.Get(ctx)
	if save == nil {
		return errors.New("context信息获取失败")
	}
	users := readUser()
	if _, ok := users[name]; ok {
		save.LoginName = name
		save.Authority = users[name]
	} else {
		return errors.New("用户不存在")
	}
	save.List = readProduct()
	return nil
}

// 模拟读取用户信息，文件路径请根据运行环境自行修改。
func readUser() map[string]string {
	f, err := ioutil.ReadFile("/home/windf/go/src/contextdemo/app/model/users.txt")
	if err != nil {
		return nil
	}
	userList := strings.Split(string(f), "\n")
	list := make(map[string]string)
	for _, i := range userList {
		line := strings.Fields(i)
		list[line[0]] = line[1]
	}
	return list
}

// 模拟读取产品信息，文件路径请根据运行环境自行修改。
func readProduct() []string {
	f, err := ioutil.ReadFile("/home/windf/go/src/contextdemo/app/model/list.txt")
	if err != nil {
		return nil
	}
	list := strings.Split(string(f), "\n")
	return list
}

// 模拟登录
var loginName = "张三"

func login() string {
	return loginName
}

// ListTable 获取context中的产品列表
func (s *contextService) ListTable(ctx context.Context) []string {
	if v := Context.Get(ctx); v != nil {
		return v.List
	}
	return nil
}

// AuthZero 更改context中的用户权限
func (s *contextService) AuthZero(ctx context.Context) {
	if v := Context.Get(ctx); v != nil {
		v.Authority = "0"
	}
}
