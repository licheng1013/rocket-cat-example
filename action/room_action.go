package action

import (
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
	"rocket-cat-example/dto"
	"sync"
	"time"
)

var matchMap sync.Map //匹配map
const maxPlayer = 2

func init() {
	room := RoomAction{}
	app.Gateway.Router().AddAction(7, 0, room.joinMatch)
	app.Gateway.Router().AddAction(7, 1, room.quitRoom)

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
	log.Println("加入匹配 -> ", d.LongData)

	app.Gateway.UsePlugin(core.LoginPluginId, func(r core.Plugin) {
		// 登入代码
		r.(*core.LoginPlugin).Login(d.LongData, ctx.SocketId)
		// 加入到匹配队列
		matchMap.Store(d.LongData, nil)
	})
}

func (a RoomAction) quitRoom(ctx *router.Context) {

}
