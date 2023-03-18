package action

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/router"
	"rocket-cat-example/app"
	"rocket-cat-example/entity"
)

func init() {
	user := UserAction{}
	app.Gateway.Router().AddAction(1, 1, user.login)
}

type UserAction struct {
}

func (a UserAction) login(ctx *router.Context) {
	var login entity.Login
	_ = ctx.Message.Bind(&login)
	fmt.Println("服务端-> ", login)
	login.UserId = 2222
	ctx.Message.SetBody(login)
}
