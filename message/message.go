package message

import (
	"encoding/json"
	"fmt"

	"github.com/botuniverse/go-libonebot/utils"
)

type Message []Segment

// TODO: Reduce and other methods

func MessageFromJSON(msgJSONString string) (Message, error) {
	msg := Message{}
	err := json.Unmarshal(utils.StringToBytes(msgJSONString), &msg)
	if err != nil {
		return nil, fmt.Errorf("消息解析失败, 错误: %v", err)
	}
	for idx, seg := range msg {
		if seg.Type == "" || seg.Data == nil {
			return nil, fmt.Errorf("消息解析失败, 第 %v 个消息段无效", idx)
		}
	}
	return msg, nil
}
