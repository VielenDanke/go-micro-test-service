package handler

import (
	"context"

	"github.com/vielendanke/go-micro-test-service/model"
	pb "github.com/vielendanke/go-micro-test-service/proto"
	"github.com/vielendanke/go-micro-test-service/repository"
)

// GrpcHandler ...
type GrpcHandler struct {
	messageRepository repository.MessageRepository
}

// NewGrpcHandler ...
func NewGrpcHandler(messageRepository repository.MessageRepository) pb.MessageServiceHandler {
	return &GrpcHandler{messageRepository: messageRepository}
}

// MessageByID ...
func (gh *GrpcHandler) MessageByID(ctx context.Context, req *pb.SingleMessageRequest, resp *pb.SingleMessageResponse) error {
	msg, err := gh.messageRepository.GetMessageByID(ctx, req.GetMessageId())
	if err != nil {
		return err
	}
	resp.MessageId = msg.ID
	resp.MessageDescription = msg.MessageDescription
	return nil
}

// SaveMessage ...
func (gh *GrpcHandler) SaveMessage(ctx context.Context, req *pb.SaveMessageRequest, resp *pb.SaveMessageResponse) error {
	msg := model.NewMessage()
	msg.BeforeSaving()
	msg.MapToMessage(req)
	if err := gh.messageRepository.SaveMessage(ctx, msg); err != nil {
		return err
	}
	resp.MessageId = msg.ID
	return nil
}

// GetMessageStream ...
func (gh *GrpcHandler) GetMessageStream(ctx context.Context, req *pb.GetListMessageRequest, stream pb.MessageService_GetMessageStreamStream) error {
	messages, err := gh.messageRepository.FindAllMessages(ctx)
	if err != nil {
		return err
	}
	for _, v := range messages {
		if err := stream.Send(&pb.GetListMessageResponse{
			MessageId:          v.ID,
			MessageDescription: v.MessageDescription,
		}); err != nil {
			return err
		}
	}
	return nil
}

// GetValidAPICall ...
func (gh *GrpcHandler) GetValidAPICall(ctx context.Context, req *pb.APIRequest, resp *pb.APIResponse) error {
	return nil
}

// GetInvalidAPICall ...
func (gh *GrpcHandler) GetInvalidAPICall(ctx context.Context, req *pb.APIRequest, resp *pb.APIResponse) error {
	return nil
}
