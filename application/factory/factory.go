package factory

import (
	"github.com/jfilipedias/codepix-go/application/usecase"
	"github.com/jfilipedias/codepix-go/infra/repository"
	"gorm.io/gorm"
)

func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}
    transactionUseCase := usecase.TransactionUseCase{
    	TransactionRepository:  &transactionRepository,
    	PixKeyRepository:      pixRepository,
    }

    return transactionUseCase
}
