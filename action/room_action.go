package action

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
	"rocket-cat-example/dto"
	"rocket-cat-example/entity"
	"sync"
	"time"
)

var matchMap sync.Map //匹配map
const maxPlayer = 2

func init() {
	room := RoomAction{}
	app.Gateway.Router().AddAction(7, 0, room.joinMatch)
	app.Gateway.Router().AddAction(7, 1, room.move)

	ticker := time.NewTicker(1 * time.Second) // 匹配管理器
	go func() {
		for _ = range ticker.C {
			var userIds []int64
			i := 0
			matchMap.Range(func(key, value any) bool {
				i++
				userIds = append(userIds, key.(int64))
				if i == maxPlayer { // 如果存在两个人，则匹配成功，就立即退出
					log.Println("两个匹配玩家id -> ", userIds)
					for _, item := range userIds { // 移除匹配列表
						matchMap.Delete(item)
					}
					room.matchOk(userIds)
					return false
				}
				return true
			})
		}
	}()
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
	matchMap.Store(d.LongData, nil)
}

const roomId = 200

func (a RoomAction) move(ctx *router.Context) {
	d := dto.TestDto{}
	_ = ctx.Message.Bind(&d)
	fmt.Println("收到消息 -> ", &d)
	ctx.Message = nil
	room := common.RoomManger.GetByRoomId(roomId)
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
		plugin.SendByUserIdMessage(app.Decoder.EncodeBytes(&message), room.UserIds()...)
	}
}

// 匹配成功
func (a RoomAction) matchOk(ids []int64) {
	room := common.RoomManger.CreateRoom()
	room.RoomId = roomId // 暂时写死
	common.RoomManger.AddRoom(room)
	r := app.Gateway.GetPlugin(core.LoginPluginId)
	loginPlugin := r.(*core.LoginPlugin)
	testDto := dto.ListTestDto{}
	for _, uid := range ids {
		common.RoomManger.PlayerJoinRoom(&entity.Player{Uid: uid}, room.RoomId)
		testDto.List = append(testDto.List, &dto.TestDto{UserId: uid})
	}
	message := messages.ProtoMessage{}
	message.SetBody(&testDto)
	message.SetMerge(common.CmdKit.GetMerge(7, 2))
	loginPlugin.SendByUserIdMessage(app.Decoder.EncodeBytes(&message), ids...)
}
