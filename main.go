package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	jsoncodec "github.com/unistack-org/micro-codec-json"
	fileconfig "github.com/unistack-org/micro-config-file"
	grpcsrv "github.com/unistack-org/micro-server-grpc"
	httpsrv "github.com/unistack-org/micro-server-http"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/api"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	servicehandler "github.com/vielendanke/test-service/handler"
	pb "github.com/vielendanke/test-service/proto"
)

func configureHandlerToEndpoints(router *mux.Router, handler interface{}, endpoints []*api.Endpoint) error {
	for _, v := range endpoints {
		fName := strings.Split(v.Name, ".")[1]
		f, ok := reflect.ValueOf(handler).MethodByName(fName).Interface().(func(w http.ResponseWriter, r *http.Request))
		if !ok {
			return fmt.Errorf("Naming of method is incorrect: %v:%v", fName, v.Path)
		}
		for k, j := range v.Path {
			router.HandleFunc(j, f).Methods(v.Method[k])
		}
	}
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)

	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.TraceLevel))

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		errCh <- fmt.Errorf("%s", <-sigChan)
	}()

	go func() {
		grpcServer := grpcsrv.NewServer(
			server.Name("test-service-grpc"),
			server.Version("latest"),
			server.Address(":9094"),
		)
		options := append([]micro.Option{},
			micro.Context(ctx),
			micro.Name("test-service-grpc"),
			micro.Version("latest"),
		)
		srv := micro.NewService(options...)

		if err := srv.Init(); err != nil {
			errCh <- err
		}
		if err := srv.Init(
			micro.Server(grpcServer),
		); err != nil {
			errCh <- err
		}
		if err := srv.Server().Handle(srv.Server().NewHandler(&servicehandler.GrpcHandler{})); err != nil {
			errCh <- err
		}
		if err := srv.Run(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		cfg := NewConfig("test-service", "1.0")

		if err := config.Load(ctx,
			config.NewConfig(
				config.Struct(cfg),
			),
			config.NewConfig(
				config.AllowFail(true),
				config.Struct(cfg),
				config.Codec(jsoncodec.NewCodec()),
				fileconfig.Path("local.json"),
			),
		); err != nil {
			logger.Fatalf("Error while loading config %v\n", err)
		}
		options := append([]micro.Option{},
			micro.Server(httpsrv.NewServer()),
			micro.Context(ctx),
			micro.Name(cfg.Server.Name),
			micro.Version(cfg.Server.Version),
		)
		svc := micro.NewService(options...)

		if err := svc.Init(); err != nil {
			errCh <- err
		}

		if err := svc.Init(
			micro.Server(httpsrv.NewServer(
				server.Name(cfg.Server.Name),
				server.Version(cfg.Server.Version),
				server.Address(cfg.Server.Addr),
				server.Codec("application/json", jsoncodec.NewCodec()),
			)),
		); err != nil {
			errCh <- err
		}
		router := mux.NewRouter()

		handler := &servicehandler.Handler{Codec: jsoncodec.NewCodec()}

		endpoints := pb.NewTestServiceEndpoints()

		if err := configureHandlerToEndpoints(router, handler, endpoints); err != nil {
			errCh <- err
		}

		router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logger.Infof("Not found, %v\n", r)
		})
		router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logger.Infof("Method not allowed, %v\n", r)
		})
		if err := svc.Server().Handle(svc.Server().NewHandler(router)); err != nil {
			errCh <- err
		}
		if err := svc.Run(); err != nil {
			errCh <- err
		}
	}()
	fmt.Printf("Terminated: %v\n", <-errCh)
}
