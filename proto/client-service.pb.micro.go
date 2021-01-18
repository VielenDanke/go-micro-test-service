// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/client-service.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

import (
	context "context"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
	server "github.com/unistack-org/micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for GithubService service

func NewGithubServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{
		&api.Endpoint{
			Name:    "GithubService.ValidMessage",
			Path:    []string{"/api/v1/valid-endpoint"},
			Method:  []string{"GET"},
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "GithubService.InvalidMessage",
			Path:    []string{"/api/v1/invalid-endpoint"},
			Method:  []string{"GET"},
			Handler: "rpc",
		},
	}
}

// Client API for GithubService service

type GithubService interface {
	ValidMessage(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
	InvalidMessage(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
}

type githubService struct {
	c    client.Client
	name string
}

func NewGithubService(name string, c client.Client) GithubService {
	return &githubService{
		c:    c,
		name: name,
	}
}

func (c *githubService) ValidMessage(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	rsp := &Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "GithubService.ValidMessage", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *githubService) InvalidMessage(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	rsp := &Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "GithubService.InvalidMessage", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// Server API for GithubService service

type GithubServiceHandler interface {
	ValidMessage(context.Context, *Request, *Response) error
	InvalidMessage(context.Context, *Request, *Response) error
}

func RegisterGithubServiceHandler(s server.Server, hdlr GithubServiceHandler, opts ...server.HandlerOption) error {
	type githubService interface {
		ValidMessage(ctx context.Context, req *Request, rsp *Response) error
		InvalidMessage(ctx context.Context, req *Request, rsp *Response) error
	}
	type GithubService struct {
		githubService
	}
	h := &githubServiceHandler{hdlr}
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "GithubService.ValidMessage",
		Path:    []string{"/api/v1/valid-endpoint"},
		Method:  []string{"GET"},
		Handler: "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "GithubService.InvalidMessage",
		Path:    []string{"/api/v1/invalid-endpoint"},
		Method:  []string{"GET"},
		Handler: "rpc",
	}))
	return s.Handle(s.NewHandler(&GithubService{h}, opts...))
}

type githubServiceHandler struct {
	GithubServiceHandler
}

func (h *githubServiceHandler) ValidMessage(ctx context.Context, req *Request, rsp *Response) error {
	return h.GithubServiceHandler.ValidMessage(ctx, req, rsp)
}

func (h *githubServiceHandler) InvalidMessage(ctx context.Context, req *Request, rsp *Response) error {
	return h.GithubServiceHandler.InvalidMessage(ctx, req, rsp)
}