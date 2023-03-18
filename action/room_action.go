package action

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/router"
	"rocket-cat-example/app"
)

func init() {
	room := RoomAction{}
	app.Gateway.Router().AddAction(common.CmdKit.GetMerge(2, 1), room.joinMatch)
	app.Gateway.Router().AddAction(common.CmdKit.GetMerge(2, 2), room.quitRoom)
}

// RoomAction 房间管理器
type RoomAction struct {
}

func (a RoomAction) joinMatch(ctx *router.Context) {

}

func (a RoomAction) quitRoom(ctx *router.Context) {

}
