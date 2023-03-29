package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/connect"
	"log"
	"net/url"
	"rocket-cat-example/entity"
	"rocket-cat-example/player"
)

func main() {
	u := url.URL{Scheme: "ws", Host: connect.Addr, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err.Error())
	}
	e := &entity.Login{UserId: int64(player.P1)}
	tool := player.JsonDecoder.Tool(player.Cmd, player.AddMatch, e)
	err = conn.WriteMessage(websocket.BinaryMessage, tool)
	for {
		_, m, err := conn.ReadMessage()
		player.Message(m)
		fmt.Println("客户端-> ", e)
		if err != nil {
			log.Println("读取消息错误:", err)
		}
	}
}
