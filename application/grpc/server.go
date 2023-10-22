package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/jfilipedias/codepix-go/application/grpc/pb"
	"github.com/jfilipedias/codepix-go/application/usecase"
	"github.com/jfilipedias/codepix-go/infra/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pixKeyRepository := repository.PixKeyRepositoryDb{Db: database}
	pixKeyUseCase := usecase.PixKeyUseCase{PixKeyRepository: pixKeyRepository}
	pixGrpcService := NewPixGrpcService(pixKeyUseCase)
	pb.RegisterPixServiceServer(grpcServer, pixGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Cannot start grpc server", err)
	}

	log.Printf("gRPC server has been started on port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start grpc server", err)
	}
}
