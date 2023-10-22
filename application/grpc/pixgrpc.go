package grpc

import (
	"context"

	"github.com/jfilipedias/codepix-go/application/grpc/pb"
	"github.com/jfilipedias/codepix-go/application/usecase"
)

type PixGrpcService struct {
	PixUseCase usecase.PixKeyUseCase
	pb.UnimplementedPixServiceServer
}

func (pixGrpcService PixGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	key, err := pixGrpcService.PixUseCase.Register(in.Key, in.Kind, in.AccountId)
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "Not created",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Id:     key.ID,
		Status: "Created",
	}, nil
}

func (pixGrpcService PixGrpcService) FindKeyByKind(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixKey, err := pixGrpcService.PixUseCase.FindKeyByKind(in.Key, in.Kind)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id:   pixKey.ID,
		Key:  pixKey.Key,
		Kind: pixKey.Kind,
		Account: &pb.Account{
			Id:        pixKey.AccountID,
			Number:    pixKey.Account.Number,
			BankId:    pixKey.Account.BankID,
			BankName:  pixKey.Account.Bank.Name,
			OwnerName: pixKey.Account.OwnerName,
			CreatedAt: pixKey.Account.CreatedAt.String(),
		},
		CreatedAt: pixKey.CreatedAt.String(),
	}, nil
}

func NewPixGrpcService(usecase usecase.PixKeyUseCase) *PixGrpcService {
	return &PixGrpcService{
		PixUseCase: usecase,
	}
}
