package main

import (
	"rocket-cat-example/action"
	"rocket-cat-example/app"
	"rocket-cat-example/config"
)

func main() {
	action.Init()                       // 让action先注册
	config.Init()                       // 注册配置
	app.Gateway.SetDecoder(app.Decoder) // 设置编码器
	app.Gateway.SetSocket(app.Socket)   // 设置socket
	app.Gateway.Start(":10100")         //启动服务,这行注释，把下面的注释取消启动则可以看到一个消息发送demo

}
