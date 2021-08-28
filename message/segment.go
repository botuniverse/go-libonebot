package message

type Segment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func TextSegment(text string) Segment {
	return Segment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

func MentionSegment(userID string) Segment {
	return Segment{
		Type: "mention",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	}
}
