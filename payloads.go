package rabbit

type TextMessage struct {
	SenderID    string  `json:"senderID"`
	RecipientID string  `json:"recipientID"`
	MessageID   string  `json:"mid"`
	Text        string  `json:"text"`
	Timestamp   float64 `json:"timestamp"`
}

func NewTextMessage(senderID, recipientID, messageID, text string, timestamp float64) *TextMessage {
	return &TextMessage{
		SenderID:    senderID,
		RecipientID: recipientID,
		MessageID:   messageID,
		Text:        text,
		Timestamp:   timestamp,
	}
}

// func MakeNewTextMessage(entry * Response)
