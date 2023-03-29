package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"log"
	"net/url"
	"os"
	"os/signal"
	"rocket-cat-example/action"
	"rocket-cat-example/app"
	"rocket-cat-example/config"
	"rocket-cat-example/entity"
)

func main() {
	action.Init()                       // 让action先注册
	config.Init()                       // 注册配置
	app.Gateway.SetDecoder(app.Decoder) // 设置编码器
	app.Gateway.SetSocket(app.Socket)   // 设置socket
	app.Gateway.Start(connect.Addr)     //启动服务,这行注释，把下面的注释取消启动则可以看到一个消息发送demo

	//go app.Gateway.Start(connect.Addr) // 监听端口
	//channel := make(chan int)
	//
	//startTime := time.Now()
	//go WsTest(channel)
	//select {
	//case _ = <-channel:
	//	//	time.Sleep(1 * time.Millisecond)
	//	log.Println("总时间毫秒:", time.Now().UnixMilli()-startTime.UnixMilli())
	//}
}

// WsTest 模拟客户端
func WsTest(v chan int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: connect.Addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	for {
		e := &entity.Login{UserId: 123456}
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.SetBody(e)
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err = c.WriteMessage(websocket.BinaryMessage, jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("写:", err)
			return
		}
		_, m, err := c.ReadMessage()
		jsonDecoder := decoder.JsonDecoder{}
		t := jsonDecoder.DecoderBytes(m)
		_ = t.Bind(&e)
		fmt.Println("客户端-> ", e)
		if err != nil {
			log.Println("读取消息错误:", err)
		}
		v <- 0
		break
	}
}
