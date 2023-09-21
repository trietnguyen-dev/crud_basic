package test

import (
	"book-project/config"
	pb "book-project/protobuf/gen/go"
	"book-project/server"
	"book-project/service"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

var (
	conf    *config.Config
	logger  *zap.Logger
	userSvc *service.UserScv
	bookSvc *service.BookScv
)

func server2(ctx context.Context) (pb.BookClient, func()) {
	var opts []grpc.ServerOption

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	logger.Info("server listening at", zap.String("addr", lis.Addr().String()))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBookServer(grpcServer, server.NewBookSrvServer(bookSvc))

	// Create a connection to the gRPC server
	conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		logger.Fatal("failed to connect to server", zap.Error(err))
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		grpcServer.Stop()
		conn.Close() // Close the gRPC connection
	}

	client := pb.NewBookClient(conn)

	return client, closer
}

func TestTelephoneServer_GetContact(t *testing.T) {
	ctx := context.Background()

	client, closer := server2(ctx)
	defer closer()

	type expectation struct {
		out *pb.CreateBookRes
		err error
	}

	tests := map[string]struct {
		in       *pb.CreateBookReq
		expected expectation
	}{
		"Must_Success": {
			in: &pb.CreateBookReq{
				Price:       123,
				Language:    "VN",
				Author:      "123",
				Quality:     123,
				Description: "123",
				Name:        "123",
				Category:    "123s",
			},
			expected: expectation{
				out: &pb.CreateBookRes{
					Success: true,
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.CreateBook(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Success != out.GetSuccess() {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}

		})
	}
}
