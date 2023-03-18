package action

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
	"time"
)

func init() {
	user := UserAction{}
	app.Gateway.Router().
		AddAction(common.CmdKit.GetMerge(1, 1), user.Hello)

}

type UserAction struct {
}

var startTime = time.Now()
var count int64

func (a UserAction) Hello(ctx *router.Context) {
	count++
	endTime := time.Now()
	if endTime.UnixMilli()-startTime.UnixMilli() > 1000 {
		log.Println("1s数量:", count)
		startTime = endTime
		count = 0
	}
	//app.Socket.SendMessage(ctx.Message.SetBody([]byte("Hi")).GetBytesResult()) //广播功能
	ctx.Message.SetBody([]byte("Hello"))
	//log.Println(string(ctx.Message.GetBody()))
}
