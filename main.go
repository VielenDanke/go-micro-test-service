package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

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
	"github.com/vielendanke/test-service/handler"
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

	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.TraceLevel))
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

	grpcsrv.NewServer()
	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		logger.Fatalf("Failed to init service, %v\n", err)
	}

	if err := svc.Init(
		micro.Server(httpsrv.NewServer(
			server.Name(cfg.Server.Name),
			server.Version(cfg.Server.Version),
			server.Address(cfg.Server.Addr),
			server.Codec("application/json", jsoncodec.NewCodec()),
		)),
	); err != nil {
		logger.Fatalf("Failed to init service with server, %v\n", err)
	}
	router := mux.NewRouter()

	handler := &handler.Handler{Codec: jsoncodec.NewCodec()}

	endpoints := pb.NewTestServiceEndpoints()

	if err := configureHandlerToEndpoints(router, handler, endpoints); err != nil {
		logger.Fatalf("Error while mapping functions to handler %v", err)
	}

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof("Not found, %v\n", r)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Infof("Method not allowed, %v\n", r)
	})
	if err := svc.Server().Handle(svc.Server().NewHandler(router)); err != nil {
		logger.Fatalf("Failed to set http handler: %v\n", err)
	}
	if err := svc.Run(); err != nil {
		logger.Fatalf("Failed to run service: %v\n", err)
	}
}
