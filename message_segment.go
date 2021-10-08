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
	SegTypeText       = "text"        // 纯文本消息段
	SegTypeMention    = "mention"     // 提及 (即 @) 消息段
	SegTypeMentionAll = "mention_all" // 提及所有人消息段
	SegTypeImage      = "image"       // 图片消息段
	SegTypeVoice      = "voice"       // 语音消息段
	SegTypeVideo      = "video"       // 视频消息段
	SegTypeFile       = "file"        // 文件消息段
	SegTypeLocation   = "location"    // 位置消息段
	SegTypeReply      = "reply"       // 回复消息段
)

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

// MentionAllSegment 构造一个提及所有人消息段.
func MentionAllSegment() Segment {
	return CustomSegment(SegTypeMentionAll, map[string]interface{}{})
}

// ImageSegment 构造一个图片消息段.
func ImageSegment(fileID string) Segment {
	return CustomSegment(SegTypeImage, map[string]interface{}{
		"file_id": fileID,
	})
}

// VoiceSegment 构造一个语音消息段.
func VoiceSegment(fileID string) Segment {
	return CustomSegment(SegTypeVoice, map[string]interface{}{
		"file_id": fileID,
	})
}

// VideoSegment 构造一个视频消息段.
func VideoSegment(fileID string) Segment {
	return CustomSegment(SegTypeVideo, map[string]interface{}{
		"file_id": fileID,
	})
}

// FileSegment 构造一个文件消息段.
func FileSegment(fileID string) Segment {
	return CustomSegment(SegTypeFile, map[string]interface{}{
		"file_id": fileID,
	})
}

// LocationSegment 构造一个位置消息段.
func LocationSegment(latitude float64, longitude float64, title string, content string) Segment {
	return CustomSegment(SegTypeLocation, map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
		"title":     title,
		"content":   content,
	})
}

// ReplySegment 构造一个回复消息段.
func ReplySegment(messageID string, userID string) Segment {
	return CustomSegment(SegTypeReply, map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
	})
}
