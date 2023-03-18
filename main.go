package main

import (
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
	"time"
)

func main() {
	action.Init()                                  // 让action先注册
	config.Init()                                  // 注册配置
	app.Gateway.SetDecoder(decoder.JsonDecoder{})  // 设置编码器
	go app.Gateway.Start(connect.Addr, app.Socket) // 监听端口
	channel := make(chan int)

	startTime := time.Now()
	clientCount := 1
	for i := 0; i < clientCount; i++ {
		go WsTest(channel)
	}
	for i := 0; i < clientCount; i++ {
		select {
		case ok := <-channel:
			log.Println(ok)
		}
	}
	log.Println("总时间毫秒:", time.Now().UnixMilli()-startTime.UnixMilli())
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
	var count int64
	for {
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err = c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("写:", err)
			return
		}

		_, m, err := c.ReadMessage()
		jsonDecoder := decoder.JsonDecoder{}
		_ = jsonDecoder.DecoderBytes(m)
		if err != nil {
			log.Println("读取消息错误:", err)
			return
		}
		count++
		if count >= 1 { // 接受100w个数据后退出
			v <- 0
		}
	}
}
