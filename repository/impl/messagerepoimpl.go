package impl

import (
	"context"
	"database/sql"
	"time"

	"github.com/vielendanke/test-service/model"
	"github.com/vielendanke/test-service/repository"
)

// MessageRepositoryImpl ...
type MessageRepositoryImpl struct {
	db *sql.DB
}

// NewMessageRepositoryImpl ...
func NewMessageRepositoryImpl(db *sql.DB) repository.MessageRepository {
	return &MessageRepositoryImpl{
		db: db,
	}
}

// GetMessageByID ...
func (mr *MessageRepositoryImpl) GetMessageByID(ctx context.Context, id string) (*model.Message, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	message := model.NewMessage()

	if err := mr.db.QueryRowContext(
		timeoutCtx,
		"SELECT ID, MESSAGE_DESCRIPTION FROM MESSAGES WHERE ID=$1",
		id,
	).Scan(&message.ID, &message.MessageDescription); err != nil {
		return nil, err
	}
	return message, nil
}

// FindAllMessages ...
func (mr *MessageRepositoryImpl) FindAllMessages(ctx context.Context) ([]*model.Message, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	messages := []*model.Message{}
	rows, err := mr.db.QueryContext(timeoutCtx, "SELECT * FROM MESSAGES")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		message := model.NewMessage()
		rows.Scan(&message.ID, &message.MessageDescription)
		messages = append(messages, message)
	}
	return messages, nil
}

// SaveMessage ...
func (mr *MessageRepositoryImpl) SaveMessage(ctx context.Context, message *model.Message) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	row := mr.db.QueryRowContext(
		timeoutCtx, "INSERT INTO MESSAGES(ID, MESSAGE_DESCRIPTION) VALUES($1, $2)", message.ID, message.MessageDescription,
	)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}
