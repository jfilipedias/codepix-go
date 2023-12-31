package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

type Transactions struct {
	Transaction []Transaction
}

type TransationRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	FindById(id string) (*Transaction, error)
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"required"`
	AccountFromID     string   `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount            float64  `json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyToID        string   `gorm:"column:pix_key_to_id;type:uuid;" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255)" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func (transaction *Transaction) IsValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("The amount must be greater than 0.")
	}

	if transaction.Status != TransactionPending && transaction.Status != TransactionCompleted && transaction.Status != TransactionError && transaction.Status != TransactionConfirmed {
		return errors.New("Invalid status.")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("The source can not be the same as origin.")
	}

	if err != nil {
		return err
	}

	return nil
}

func (transation *Transaction) Complete() error {
	transation.Status = TransactionCompleted
	transation.UpdatedAt = time.Now()

	err := transation.IsValid()
	return err
}

func (transation *Transaction) Confirmed() error {
	transation.Status = TransactionConfirmed
	transation.UpdatedAt = time.Now()

	err := transation.IsValid()
	return err
}

func (transation *Transaction) Cancel(description string) error {
	transation.Status = TransactionError
	transation.UpdatedAt = time.Now()
	transation.CancelDescription = description

	err := transation.IsValid()
	return err
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Description: description,
		Status:      TransactionPending,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.IsValid()

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
