package usecase

import (
	"fmt"

	"github.com/jfilipedias/codepix-go/domain/model"
)

type PixKeyUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
}

func (useCase *PixKeyUseCase) Register(key string, kind string, accountId string) (*model.PixKey, error) {
	account, err := useCase.PixKeyRepository.FindAccountById(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(key, kind, account)
	if err != nil {
		return nil, err
	}

	useCase.PixKeyRepository.RegisterKey(pixKey)
	if pixKey.ID == "" {
		return nil, fmt.Errorf("Unable to create a new key at this moment.")
	}

	return pixKey, nil
}

func (useCase *PixKeyUseCase) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	pixKey, err := useCase.PixKeyRepository.FindKeyByKind(key, kind)
	if err != nil {
		return nil, err
	}
	return pixKey, err
}
