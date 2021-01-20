package service

import (
	"context"
	"database/sql"

	grpcsrv "github.com/unistack-org/micro-server-grpc"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/server"
	servicehandler "github.com/vielendanke/test-service/handler"
	pb "github.com/vielendanke/test-service/proto"
	messagerepoimpl "github.com/vielendanke/test-service/repository/impl"
)

// StartGRPCService ...
func StartGRPCService(ctx context.Context, errCh chan<- error, dbCh <-chan *sql.DB) {
	grpcServer := grpcsrv.NewServer(
		server.Name("test-service-grpc"),
		server.Version("latest"),
		server.Address(":7074"),
		server.Context(ctx),
		grpcsrv.Reflection(true),
	)

	options := append([]micro.Option{},
		micro.Server(grpcServer),
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
	db := <-dbCh

	messageRepo := messagerepoimpl.NewMessageRepositoryImpl(db)
	if err := pb.RegisterMessageServiceHandler(
		srv.Server(),
		servicehandler.NewGrpcHandler(messageRepo),
	); err != nil {
		errCh <- err
	}
	if err := srv.Run(); err != nil {
		errCh <- err
	}
}
