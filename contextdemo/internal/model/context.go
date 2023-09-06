package model

const (
	ContextKey = "contextKey"
)

type Context struct {
	LoginName string   // 登录用户名
	Authority string   // 用户权限
	List      []string // 产品清单
}
