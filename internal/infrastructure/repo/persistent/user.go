package persistent

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pool *postgres.Postgres) *UserRepo {
	return &UserRepo{pool}
}

func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	query := `
        INSERT INTO users (email, password, role)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	err := r.Pool.QueryRow(ctx, query, u.Email, u.Password, u.Role).
		Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%w: %s", entity.ErrUserAlreadyExists, err)
			}
		}

		return fmt.Errorf("%w: %s", entity.ErrInternal, err)
	}

	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
        SELECT id, email, password, role, created_at
        FROM users
        WHERE email = $1
    `
	var u entity.User
	err := r.Pool.QueryRow(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}

		return nil, entity.ErrInternal
	}
	return &u, nil
}
