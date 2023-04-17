package app

import (
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/router"
)

var Gateway = core.DefaultGateway()
var Socket = &connect.WebSocket{}
var Decoder = decoder.ProtoDecoder{}

func init() {
	// 添加中间件
	Gateway.Router().AddProxy(&MyProxy{})
}

type MyProxy struct {
	router.ProxyFunc
}

func (m *MyProxy) InvokeFunc(ctx *router.Context) {
	//log.Println("执行前")
	m.ProxyFunc.InvokeFunc(ctx)
	//log.Println("执行后")
}
