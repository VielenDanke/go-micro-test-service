package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/vielendanke/test-service/model"
	pb "github.com/vielendanke/test-service/proto"
	"github.com/vielendanke/test-service/repository"
)

// Handler ...
type Handler struct {
	codec             codec.Codec
	messageRepository repository.MessageRepository
}

// NewHandler ...
func NewHandler(messageRepository repository.MessageRepository, codec codec.Codec) *Handler {
	return &Handler{messageRepository: messageRepository, codec: codec}
}

// MessageByID ...
func (h *Handler) MessageByID(w http.ResponseWriter, r *http.Request) {
	messageID, ok := mux.Vars(r)["message_id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	msg, err := h.messageRepository.GetMessageByID(context.Background(), messageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.codec.Write(w, nil, msg)
}

// SaveMessage ...
func (h *Handler) SaveMessage(w http.ResponseWriter, r *http.Request) {
	msg := &pb.SaveMessageRequest{}
	h.codec.ReadBody(r.Body, msg)
	modelMsg := model.NewMessage()
	modelMsg.BeforeSaving()
	modelMsg.MapToMessage(msg)
	if err := h.messageRepository.SaveMessage(context.Background(), modelMsg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetMessageStream ...
func (h *Handler) GetMessageStream(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageRepository.FindAllMessages(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.codec.Write(w, nil, messages)
}
