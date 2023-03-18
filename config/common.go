package config

import (
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
)


func Init() {
	log.Println("初始化配置")
	// 添加中间件
	app.Gateway.Router().AddProxy(&MyProxy{})
}

type MyProxy struct {
	proxy router.Proxy
}

func (m *MyProxy) InvokeFunc(ctx *router.Context) {
	//log.Println("执行前")
	m.proxy.InvokeFunc(ctx)
	//log.Println("执行后")
}

func (m *MyProxy) SetProxy(proxy router.Proxy) {
	m.proxy = proxy
}


