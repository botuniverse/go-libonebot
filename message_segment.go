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

// SegTypeXxx 表示 OneBot 标准定义的核心消息段类型.
const (
	SegTypeText    = "text"    // 纯文本消息段
	SegTypeMention = "mention" // 提及 (即 @) 消息段
)

func segmentFromMap(m map[string]interface{}) (Segment, error) {
	em := EasierMapFromMap(m)
	t, _ := em.GetString("type")
	if t == "" {
		return Segment{}, fmt.Errorf("消息段 `type` 字段不存在或为空")
	}
	data, err := em.GetMap("data")
	if err != nil {
		data = EasierMapFromMap(map[string]interface{}{})
	}
	return Segment{
		Type: t,
		Data: data,
	}, nil
}

func (s *Segment) tryMerge(next Segment) bool {
	switch s.Type {
	case SegTypeText:
		if next.Type == SegTypeText {
			text1, err1 := s.Data.GetString("text")
			text2, err2 := next.Data.GetString("text")
			if err1 != nil && err2 == nil {
				s.Data.Set("text", text2)
			} else if err1 == nil && err2 != nil {
				s.Data.Set("text", text1)
			} else if err1 == nil && err2 == nil {
				s.Data.Set("text", text1+text2)
			} else {
				s.Data.Set("text", "")
			}
			return true
		}
	}
	return false
}

// CustomSegment 构造一个指定类型的消息段.
func CustomSegment(type_ string, data map[string]interface{}) Segment {
	return Segment{
		Type: type_,
		Data: EasierMapFromMap(data),
	}
}

// TextSegment 构造一个纯文本消息段.
func TextSegment(text string) Segment {
	return CustomSegment(SegTypeText, map[string]interface{}{
		"text": text,
	})
}

// MentionSegment 构造一个提及消息段.
func MentionSegment(userID string) Segment {
	return CustomSegment(SegTypeMention, map[string]interface{}{
		"user_id": userID,
	})
}
