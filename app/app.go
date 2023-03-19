package app

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
)

var Gateway = core.DefaultGateway()
var Socket = &connect.WebSocket{}
var Decoder = decoder.ProtoDecoder{}

func init() {
	Socket.Pool = common.NewPool(100, 30)
}
