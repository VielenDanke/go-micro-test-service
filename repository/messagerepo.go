package repository

import (
	"context"

	"github.com/vielendanke/test-service/model"
)

// MessageRepository ...
type MessageRepository interface {
	GetMessageByID(context.Context, string) (*model.Message, error)
	FindAllMessages(context.Context) ([]*model.Message, error)
	SaveMessage(context.Context, *model.Message) error
}
