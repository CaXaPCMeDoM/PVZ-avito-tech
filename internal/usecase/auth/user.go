package auth

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/infrastructure/repo"
	"PVZ-avito-tech/internal/infrastructure/security"
	"context"
)

type UserUsecase struct {
	repo   repo.UserRepo
	hasher security.PasswordHasher
}

func NewUserUsecase(
	repo repo.UserRepo,
	hasher security.PasswordHasher,
) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		hasher: hasher,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, u *entity.User) (RegisterResponse, error) {
	response := RegisterResponse{}

	hashedPass, err := uc.hasher.Hash(u.Password)
	if err != nil {
		return response, err
	}

	u.Password = hashedPass

	if err = uc.repo.Create(ctx, u); err != nil {
		return response, err
	}

	response = RegisterResponse{
		Id:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}

	return response, err
}

func (uc *UserUsecase) Login(ctx context.Context, email string, rawPassword string) (LoginResponse, error) {
	u, err := uc.repo.GetByEmail(ctx, email)

	if err != nil {
		return LoginResponse{}, err
	}

	if err = uc.hasher.Verify(u.Password, rawPassword); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Role: u.Role}, nil
}
