package onebot

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type Message []Segment

func (m Message) String() string {
	j, _ := json.Marshal(m)
	return "onebot.Message" + bytesToString(j)
}

func MessageFromJSON(j gjson.Result) (Message, error) {
	if j.Type == gjson.String {
		return Message{TextSegment(j.Str)}, nil
	}

	var msgJSONString string
	if j.IsObject() {
		msgJSONString = "[" + j.Raw + "]"
	} else if j.IsArray() {
		msgJSONString = j.Raw
	} else {
		return nil, fmt.Errorf("消息解析失败, 不是有效的消息格式")
	}

	msg := Message{}
	err := json.Unmarshal(stringToBytes(msgJSONString), &msg)
	if err != nil {
		return nil, fmt.Errorf("消息解析失败, 错误: %v", err)
	}
	return msg, nil
}
