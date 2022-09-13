// 接口定义 - 消息接口 - 消息数据类型
// https://12.onebot.dev/interface/message/type/

package libonebot

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

// Segment 表示一个消息段.
type Segment struct {
	Type string    // 消息段类型
	Data EasierMap // 消息段数据
}

func (s Segment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": s.Type,
		"data": s.Data.Value(),
	})
}

func (s Segment) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(map[string]interface{}{
		"type": s.Type,
		"data": s.Data.Value(),
	})
}

func segmentFromMap(m map[string]interface{}) (Segment, error) {
	em := EasierMapFromMap(m)
	t, _ := em.GetString("type")
	if t == "" {
		return Segment{}, fmt.Errorf("消息段 `type` 字段不存在或为空")
	}
	data, err := em.GetMap("data")
	if err != nil {
		data = EasierMapFromMap(make(map[string]interface{}))
	}
	return Segment{
		Type: t,
		Data: data,
	}, nil
}

// Message 表示一条消息.
type Message []Segment

// Reduce 合并消息中连续的可合并消息段 (如连续的纯文本消息段).
func (m *Message) Reduce() {
	for i := 0; i < len(*m)-1; i++ {
		j := i + 1
		for ; j < len(*m) && (*m)[i].tryMerge((*m)[j]); j++ {
		}
		if i+1 != j {
			*m = append((*m)[:i+1], (*m)[j:]...)
		}
	}
}

// ExtractText 提取消息中的纯文本消息段, 并合并为字符串.
func (m *Message) ExtractText() (text string) {
	for _, s := range *m {
		if s.Type == SegTypeText {
			t, _ := s.Data.GetString("text")
			text += t
		}
	}
	return
}

func messageFromInterface(interf interface{}) (Message, error) {
	switch v := interf.(type) {
	case string:
		return Message{TextSegment(v)}, nil
	case map[string]interface{}:
		seg, err := segmentFromMap(v)
		if err != nil {
			return nil, err
		}
		return Message{seg}, nil
	case []interface{}:
		segs := make([]Segment, len(v))
		for i, s := range v {
			switch s := s.(type) {
			case string:
				segs[i] = TextSegment(s)
			case map[string]interface{}:
				seg, err := segmentFromMap(s)
				if err != nil {
					return nil, err
				}
				segs[i] = seg
			default:
				return nil, fmt.Errorf("消息解析失败, 不是有效的消息格式")
			}
		}
		return segs, nil
	default:
		return nil, fmt.Errorf("消息解析失败, 不是有效的消息格式")
	}
}
