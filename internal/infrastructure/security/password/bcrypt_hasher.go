package password

import (
	"PVZ-avito-tech/config"
	"PVZ-avito-tech/internal/entity"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cfg *config.Config) *BcryptHasher {
	return &BcryptHasher{cost: cfg.Security.PasswordCost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	if len(password) == 0 {
		return "", entity.ErrInvalidPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)

	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "", entity.ErrPasswordTooLong
	}

	if err != nil {
		return "", entity.ErrPasswordHashing
	}
	return string(bytes), nil
}

func (h *BcryptHasher) Verify(hashedPassword, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(inputPassword),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return entity.ErrPasswordVerify
		}
		return err
	}

	return nil
}
