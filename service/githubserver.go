package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jsoncodec "github.com/unistack-org/micro-codec-json"
	httpsrv "github.com/unistack-org/micro-server-http"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/server"
	apiendpoints "github.com/vielendanke/test-service/endpoints"
	servicehandler "github.com/vielendanke/test-service/handler"
	pb "github.com/vielendanke/test-service/proto"
)

// StartGithubService ...
func StartGithubService(ctx context.Context, errCh chan<- error) {
	options := append([]micro.Option{},
		micro.Server(httpsrv.NewServer()),
		micro.Context(ctx),
		micro.Name("github-service"),
		micro.Version("latest"),
	)
	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		errCh <- err
	}
	if err := svc.Init(
		micro.Server(httpsrv.NewServer(
			server.Name("github-service"),
			server.Version("latest"),
			server.Address(":9095"),
			server.Context(ctx),
			server.Codec("application/json", jsoncodec.NewCodec()),
		)),
	); err != nil {
		errCh <- err
	}
	router := mux.NewRouter()

	handler := servicehandler.NewGithubHandler(jsoncodec.NewCodec())

	endpoints := pb.NewGithubServiceEndpoints()

	if err := apiendpoints.ConfigureHandlerToEndpoints(router, handler, endpoints); err != nil {
		errCh <- err
	}
	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("Not found, %v\n", r)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf("Method not allowed, %v\n", r)
	})
	if err := svc.Server().Handle(svc.Server().NewHandler(router)); err != nil {
		errCh <- err
	}
	if err := svc.Run(); err != nil {
		errCh <- err
	}
}
