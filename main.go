package main

import (
	"book-project/config"
	"book-project/daos"
	pb "book-project/protobuf/gen/go"
	"book-project/server"
	"book-project/service"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	// EnvProduction :
	EnvProduction = "production"
)

var (
	conf    *config.Config
	logger  *zap.Logger
	userSvc *service.UserScv
	bookSvc *service.BookScv
)

func init() {
	conf = config.GetConfig()

	// init log
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	dao, err := daos.NewDAO(conf)
	if err != nil {
		logger.Fatal("failed to init daos", zap.Error(err))
	}

	userSvc = service.NewUserSvc(dao, conf)
	bookSvc = service.NewBookSvc(dao, conf)
}

func main() {

	var (
		opts []grpc.ServerOption
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	logger.Info("server listening at", zap.String("addr", lis.Addr().String()))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserSrvServer(grpcServer, server.NewUserSrvServer(
		userSvc,
	))
	pb.RegisterBookServer(grpcServer, server.NewBookSrvServer(
		bookSvc))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
