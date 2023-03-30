package action

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
	"rocket-cat-example/app"
	"rocket-cat-example/dto"
	"rocket-cat-example/entity"
	"rocket-cat-example/player"
	"time"
)

var queue *common.MatchQueue

func init() {
	queue = common.NewMatchQueue(2, func(userIds []int64) {
		fmt.Println("匹配成功 -> ", userIds)
		// 构建房间管理器
		room := common.RoomManger.CreateRoom()
		common.RoomManger.AddRoom(room)
		for _, userId := range userIds {
			common.RoomManger.PlayerJoinRoom(&entity.Player{Uid: userId}, room.RoomId)
		}
		// 帧编号
		var frameId int32
		room.List = append(room.List, &common.SafeList{})
		fmt.Println("帧同步开始")
		room.StartCustom(func() {
			list := room.List[frameId]
			dataList := list.GetList()
			dtoList := convert(dataList)
			// 发送第一帧数据
			testDto := dto.ListTestDto{Frame: frameId, List: dtoList}
			// 获取登入插件并发送数据
			r := app.Gateway.GetPlugin(1)
			loginPlugin := r.(*core.LoginPlugin)
			loginPlugin.SendByUserIdMessage(app.Decoder.Tool(player.Cmd, player.SyncCmd, &testDto), room.UserIds()...)
			frameId++
			room.List = append(room.List, &common.SafeList{})
		}, time.Second/20) // 20帧
	})

	user := UserAction{}
	app.Gateway.Router().AddAction(player.Cmd, player.Login, user.Login)
	app.Gateway.Router().AddAction(player.Cmd, player.AddMatch, user.AddMatch)

	// 清理房间
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 1)
	//		common.RoomManger.RoomClear(3)
	//	}
	//}()
}

// dataList转换为TestDto
func convert(dataList []interface{}) []*dto.TestDto {
	var list []*dto.TestDto
	for _, data := range dataList {
		list = append(list, data.(*dto.TestDto))
	}
	return list
}

type UserAction struct {
}

func (a UserAction) Login(ctx *router.Context) {
	var login entity.Login
	_ = ctx.Message.Bind(&login)
	fmt.Println("服务端-> ", login)
	login.UserId = 2222
	ctx.Message.SetBody(login)
}

func (a UserAction) AddMatch(ctx *router.Context) {
	var login entity.Login
	_ = ctx.Message.Bind(&login)
	// 使用登入插件
	r := app.Gateway.GetPlugin(1)
	loginPlugin := r.(*core.LoginPlugin)
	loginPlugin.Login(login.UserId, ctx.SocketId)
	queue.AddMatch(login.UserId)
}
