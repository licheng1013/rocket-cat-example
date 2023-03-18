package action

import (
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"rocket-cat-example/app"
)

func Init() {
	app.Gateway.Router().(*router.DefaultRouter).DebugLog = true
	log.Println("初始化Action")
}
