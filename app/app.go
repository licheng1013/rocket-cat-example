package app

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
)

var Gateway = core.DefaultGateway()
var Socket = &connect.WebSocket{}
var Decoder = decoder.JsonDecoder{}

func init() {
	Socket.OnClose(func(socketId uint32) {
		fmt.Println(socketId, "连接关闭")
	})
}
