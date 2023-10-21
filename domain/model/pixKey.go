package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKeyRepositoryInterface interface {
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	Register(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	FindAccountById(id string) (*Account, error)
	FindBankById(id string) (*Bank, error)
}

type PixKey struct {
	Base      `valid:"required"`
	Key       string   `json:"key" valid:"notnull"`
	Kind      string   `json:"kind" valid:"notnull"`
	Account   *Account `valid:"-"`
	AccountID string   `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Status    string   `json:"status" valid:"notnull"`
}

func (pixKey *PixKey) IsValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("Invalid type of key.")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("Invalid type of status.")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewPixKey(key string, kind string, account *Account) (*PixKey, error) {
	pixKey := PixKey{
		Account: account,
		Kind:    kind,
		Key:     key,
		Status:  "active",
	}
	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.IsValid()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}
