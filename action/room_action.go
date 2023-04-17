package action

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
	"rocket-cat-example/cmd"
	"rocket-cat-example/dto"
	"rocket-cat-example/entity"
)

var matchQueue *common.MatchQueue
var roomManager = common.NewRoomManger()

func init() {
	room := RoomAction{}
	app.Gateway.Router().AddAction(cmd.Room, cmd.JoinMatch, room.joinMatch)
	app.Gateway.Router().AddAction(cmd.Room, cmd.Move, room.move)
	app.Gateway.Router().AddAction(cmd.Room, cmd.QuitRoom, room.quitRoom)

	matchQueue = common.NewMatchQueue(2, room.matchOk)
}

// RoomAction 房间管理器
type RoomAction struct {
}

func (a RoomAction) joinMatch(ctx *router.Context) {
	// 加入房间
	d := dto.MessageDto{}
	_ = ctx.Message.Bind(&d)
	r := app.Gateway.GetPlugin(core.LoginPluginId)
	// 登入代码
	loginPlugin := r.(*core.LoginPlugin)
	loginPlugin.Login(d.LongData, ctx.SocketId)
	log.Printf("加入匹配 -> %v - %v", d.LongData, ctx.SocketId)
	// 加入到匹配队列
	matchQueue.AddMatch(&entity.Player{Uid: d.LongData})
}

func (a RoomAction) move(ctx *router.Context) {
	var d dto.TestDto
	_ = ctx.Message.Bind(&d)
	fmt.Println("收到消息 -> ", &d)
	room := roomManager.GetByUserId(d.GetUserId())
	if room != nil {
		r := app.Gateway.GetPlugin(core.LoginPluginId)
		plugin := r.(*core.LoginPlugin)
		testDto := dto.ListTestDto{}
		testDto.List = append(testDto.List, &d)

		// 处理
		message := messages.ProtoMessage{}
		message.SetBody(&testDto)
		message.SetMerge(common.CmdKit.GetMerge(7, 1))
		log.Println("发送消息: -> ", &testDto)
		plugin.SendByUserIdMessage(app.Decoder.EncodeBytes(&message), room.GetUserIdList()...)
	}
	ctx.Message = nil
}

// 匹配成功
func (a RoomAction) matchOk(ids []int64) {
	roomId := roomManager.GetUniqueRoomId()
	room := common.DefaultRoom{RoomId: roomId}
	roomManager.AddRoom(&room)
	r := app.Gateway.GetPlugin(core.LoginPluginId)
	loginPlugin := r.(*core.LoginPlugin)
	testDto := dto.ListTestDto{}
	log.Println("创建房间")
	for _, uid := range ids {
		roomManager.JoinRoom(&entity.Player{Uid: uid}, room.RoomId)
		testDto.List = append(testDto.List, &dto.TestDto{UserId: uid})
	}
	message := messages.ProtoMessage{}
	message.SetBody(&testDto)
	message.SetMerge(common.CmdKit.GetMerge(7, 2))
	loginPlugin.SendByUserIdMessage(app.Decoder.EncodeBytes(&message), ids...)
}

func (a RoomAction) quitRoom(ctx *router.Context) {
	d := dto.MessageDto{}
	_ = ctx.Message.Bind(&d)
	// 退出房间并退出登入了
	room := roomManager.GetByUserId(d.LongData)
	if room != nil {
		roomManager.RemoveRoom(room.GetRoomId())
	}
	r := app.Gateway.GetPlugin(core.LoginPluginId)
	loginPlugin := r.(*core.LoginPlugin)
	loginPlugin.LogoutByUserId(d.LongData)
}
