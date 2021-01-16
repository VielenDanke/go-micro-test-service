package handler

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/vielendanke/test-service/proto"
)

type GrpcHandler struct {
}

func (*GrpcHandler) GetTest(ctx context.Context, req *pb.GetTestRequest, resp *pb.GetTestResponse) error {
	resp.Result = "good job"
	return nil
}
func (*GrpcHandler) PostTest(ctx context.Context, req *pb.PostTestRequest, resp *pb.PostTestResponse) error {
	resp.Id = uuid.New().String()
	return nil
}
