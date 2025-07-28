package usecase

import (
	"mymodule/pkg/cypto"
)

type CryptoService interface {
	HashedPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}

type DefaultCryptoService struct{}

func (d *DefaultCryptoService) HashedPassword(password string) (string, error) {
	return cypto.HashedPassword(password)
}
func (d *DefaultCryptoService) ComparePassword(password, hash string) bool {
	return cypto.ComparePassword(password,hash)
}