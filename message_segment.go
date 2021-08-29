package onebot

import (
	"encoding/json"
	"fmt"
)

type Segment struct {
	Type string   `json:"type"`
	Data *easyMap `json:"data"`
}

func (s Segment) String() string {
	j, _ := json.Marshal(s)
	return "onebot.Segment" + bytesToString(j)
}

// UnmarshalJSON implements json.Unmarshaler for `Segment` with validation.
func (s *Segment) UnmarshalJSON(b []byte) error {
	tmp := struct {
		Type string   `json:"type"`
		Data *easyMap `json:"data"`
	}{}
	err := json.Unmarshal(b, &tmp) // this will do the normal unmarshalling
	if err != nil {
		return fmt.Errorf("消息段格式错误")
	}
	// validate the result
	if tmp.Type == "" {
		return fmt.Errorf("消息段 `type` 字段不存在或为空")
	}
	if tmp.Data == nil {
		return fmt.Errorf("消息段 `data` 字段不存在或为空")
	}
	s.Type = tmp.Type
	s.Data = tmp.Data
	return nil
}
}

func CustomSegment(type_ string, data map[string]interface{}) Segment {
	return Segment{
		Type: type_,
		Data: newEasyMapFromMap(data),
	}
}

func TextSegment(text string) Segment {
	return CustomSegment("text", map[string]interface{}{
		"text": text,
	})
}

func MentionSegment(userID string) Segment {
	return CustomSegment("mention", map[string]interface{}{
		"user_id": userID,
	})
}
