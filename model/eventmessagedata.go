package model

import (
	pb "github.com/vielendanke/test-service/proto"
)

// MessageData ...
type MessageData struct {
	Event     *pb.EventMessage
	MediaType string
	DestTopic string
}

// Topic ...
func (md *MessageData) Topic() string {
	return md.DestTopic
}

// Payload ...
func (md *MessageData) Payload() interface{} {
	return md.Event
}

// ContentType ...
func (md *MessageData) ContentType() string {
	return md.MediaType
}
