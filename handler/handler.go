package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/codec"
	pb "github.com/vielendanke/test-service/proto"
)

type Handler struct {
	Codec codec.Codec
}

func (h *Handler) GetTest(w http.ResponseWriter, r *http.Request) {
	res := &pb.GetTestResponse{Result: "success"}
	h.Codec.Write(w, nil, res)
}

func (h *Handler) PostTest(w http.ResponseWriter, r *http.Request) {
	req := &pb.PostTestRequest{}
	h.Codec.ReadBody(r.Body, req)
	id := uuid.New().String()
	h.Codec.Write(w, nil, &pb.PostTestResponse{Id: id})
}
