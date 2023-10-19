package usecase

import (
	"fmt"

	"github.com/jfilipedias/codepix-go/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransationRepositoryInterface
	PixKeyRepository      model.PixKeyRepositoryInterface
}

func (useCase TransactionUseCase) Register(accountId string, amount float64, pixKeyTo string, pixKeyToKind string, description string) (*model.Transaction, error) {
	account, err := useCase.PixKeyRepository.FindAccountById(accountId)

	if err != nil {
		return nil, err
	}

	pixKey, err := useCase.PixKeyRepository.FindKeyByKind(pixKeyTo, pixKeyToKind)

	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)

	if err != nil {
		return nil, err
	}

	useCase.TransactionRepository.Save(transaction)

	if transaction.ID == "" {
		return nil, fmt.Errorf("Unable to process this transaction")
	}

	return transaction, nil
}

func (useCase TransactionUseCase) Confirm(transactionId string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.FindById(transactionId)

	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed

	err = useCase.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (useCase TransactionUseCase) Complete(transactionId string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.FindById(transactionId)

	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionCompleted

	err = useCase.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (useCase TransactionUseCase) Error(transactionId string, reason string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.FindById(transactionId)

	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionError
	transaction.CancelDescription = reason

	err = useCase.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
