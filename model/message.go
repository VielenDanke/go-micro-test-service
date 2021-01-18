package model

import (
	"github.com/google/uuid"
	pb "github.com/vielendanke/test-service/proto"
)

// Message ...
type Message struct {
	ID                 string `json:"id,omitempty"`
	MessageDescription string `json:"description,omitempty"`
}

// NewMessage ...
func NewMessage() *Message {
	return &Message{}
}

// BeforeSaving ...
func (m *Message) BeforeSaving() {
	m.ID = uuid.New().String()
}

// MapToMessage ...
func (m *Message) MapToMessage(req *pb.SaveMessageRequest) {
	m.MessageDescription = req.GetDescription()
}
