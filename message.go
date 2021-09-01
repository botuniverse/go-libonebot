package libonebot

import (
	"fmt"
)

type Message []Segment

func (m *Message) Reduce() {
	for i := 0; i < len(*m)-1; i++ {
		j := i + 1
		for ; j < len(*m) && (*m)[i].TryMerge((*m)[j]); j++ {
		}
		if i+1 != j {
			*m = append((*m)[:i+1], (*m)[j:]...)
		}
	}
}

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
