package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	jsoncodec "github.com/unistack-org/micro-codec-json"
	fileconfig "github.com/unistack-org/micro-config-file"
	httpsrv "github.com/unistack-org/micro-server-http"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	"github.com/vielendanke/test-service/handler"
	pb "github.com/vielendanke/test-service/proto"
)

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
	for _, v := range endpoints {
		if strings.Contains(v.Name, "Get") {
			router.HandleFunc(v.Path[0], handler.GetTest).Methods(v.Method[0])
		}
		if strings.Contains(v.Name, "Post") {
			router.HandleFunc(v.Path[0], handler.PostTest).Methods(v.Method[0])
		}
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
