package player

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/decoder"
	"rocket-cat-example/dto"
)

var JsonDecoder = &decoder.JsonDecoder{}
var P1 = 123
var P2 = 456

const Cmd = 1
const (
	Login    = 1
	AddMatch = 2
	SyncCmd  = 3
)

// Message 消息分发
func Message(data []byte) {
	message := JsonDecoder.DecoderBytes(data)
	// 分解命令
	subCmd := common.CmdKit.GetSubCmd(message.GetMerge())
	if subCmd == SyncCmd {
		// 同步消息
		var list dto.ListTestDto
		_ = message.Bind(&list)
		fmt.Println("同步消息-> ", &list)
	}
}

func SendSyncMsg(user int) []byte {
	testDto := dto.TestDto{UserId: int64(user)}
	testDto.X = 1
	tool := JsonDecoder.Tool(Cmd, SyncCmd, &testDto)
	return tool
}
