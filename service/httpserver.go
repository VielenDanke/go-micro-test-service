package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	segmentiomicro "github.com/unistack-org/micro-broker-segmentio/v3"
	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	consulconfig "github.com/unistack-org/micro-config-consul/v3"
	fileconfig "github.com/unistack-org/micro-config-file/v3"
	promwrapper "github.com/unistack-org/micro-metrics-prometheus/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
	serviceconfig "github.com/vielendanke/go-micro-test-service/config"
	apiendpoints "github.com/vielendanke/go-micro-test-service/endpoints"
	servicehandler "github.com/vielendanke/go-micro-test-service/handler"
	pb "github.com/vielendanke/go-micro-test-service/proto"
	messagerepoimpl "github.com/vielendanke/go-micro-test-service/repository/impl"
)

// StartHTTPService ...
func StartHTTPService(ctx context.Context, errCh chan<- error, dbCh <-chan *sql.DB) {
	cfg := serviceconfig.NewConfig("message-service", "1.0")

	if err := config.Load(ctx,
		config.NewConfig(
			config.Struct(cfg),
		),
		config.NewConfig(
			config.AllowFail(true),
			config.Struct(cfg),
			config.Codec(jsoncodec.NewCodec()),
			fileconfig.Path("../local.json"),
		),
		consulconfig.NewConfig(
			config.AllowFail(true),
			config.Codec(jsoncodec.NewCodec()),
			config.BeforeLoad(func(ctx context.Context, conf config.Config) error {
				return conf.Init(
					consulconfig.Address("localhost:8500"),
					consulconfig.Path("api/message-service"),
				)
			}),
		),
	); err != nil {
		errCh <- err
	}
	// Broker
	br := segmentiomicro.NewBroker(broker.Addrs("kafka:9092"), broker.Codec(jsoncodec.NewCodec()))

	options := append([]micro.Option{},
		micro.Server(httpsrv.NewServer()),
		micro.Client(httpcli.NewClient()),
		micro.Context(ctx),
		micro.Name(cfg.Server.Name),
		micro.Version(cfg.Server.Version),
		micro.Broker(br),
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
			server.Context(ctx),
			server.Codec("application/json", jsoncodec.NewCodec()),
			server.Broker(br),
			server.WrapHandler(promwrapper.NewHandlerWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
		)),
		micro.Client(httpcli.NewClient(
			client.ContentType("application/json"),
			client.Codec("application/json", jsoncodec.NewCodec()),
			client.Wrap(promwrapper.NewClientWrapper(
				promwrapper.ServiceName(svc.Server().Options().Name),
				promwrapper.ServiceVersion(svc.Server().Options().Version),
				promwrapper.ServiceID(svc.Server().Options().Id),
			)),
		)),
		micro.Configs(consulconfig.NewConfig(
			config.AllowFail(true),
			config.Codec(jsoncodec.NewCodec()),
			config.BeforeLoad(func(ctx context.Context, conf config.Config) error {
				return conf.Init(
					consulconfig.Address("localhost:8500"),
					consulconfig.Path("api/message-service"),
				)
			}))),
	); err != nil {
		errCh <- err
	}
	router := mux.NewRouter()

	db := <-dbCh

	messageRepo := messagerepoimpl.NewMessageRepositoryImpl(db)

	handler := servicehandler.NewHandler(messageRepo, jsoncodec.NewCodec(), svc.Client())

	endpoints := pb.NewMessageServiceEndpoints()

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

	// Kafka
	h := func(ctx context.Context, event *pb.EventMessage) error {
		logger.Infof(ctx, "Message processed: %v", event.GetEvent())
		return nil
	}
	if err := micro.RegisterSubscriber("message-topic", svc.Server(), h, server.InternalSubscriber(true)); err != nil {
		errCh <- fmt.Errorf("Failed to register kafka subscriber %v", err)
	}

	// consulregister.SetupConfig("http://localhost:8500")
	// consulService := consulregister.NewService(fmt.Sprintf("message-service-%s", uuid.New().String()), "message-service", 9090)
	// consulService.Tags = append(consulService.Tags, "go")
	// err := consulregister.Register(consulService)
	// if err != nil {
	// 	errCh <- err
	// }

	if err := svc.Run(); err != nil {
		errCh <- err
	}

}
