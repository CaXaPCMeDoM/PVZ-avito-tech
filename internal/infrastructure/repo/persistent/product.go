package persistent

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (r *ProductRepo) AddProduct(ctx context.Context, pvzID uuid.UUID, productType entity.ProductType) (*entity.Product, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress' 
		LIMIT 1 FOR UPDATE`,
		pvzID,
	).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrNoActiveReception
		}
		return nil, fmt.Errorf("failed to get active reception: %w", err)
	}

	var product entity.Product
	err = tx.QueryRow(ctx, `
		INSERT INTO products (reception_id, type)
		VALUES ($1, $2)
		RETURNING id, reception_id, type, created_at
	`, receptionID, productType).Scan(
		&product.ID,
		&product.ReceptionID,
		&product.Type,
		&product.DateTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &product, nil
}

func (r *ProductRepo) DeleteProductLIFO(ctx context.Context, pvzID uuid.UUID) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id FROM receptions 
		WHERE pvz_id = $1 AND status = 'in_progress' 
		LIMIT 1 FOR UPDATE`,
		pvzID,
	).Scan(&receptionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrNoActiveReception
		}
		return fmt.Errorf("failed to get active reception: %w", err)
	}

	var productID uuid.UUID
	err = tx.QueryRow(ctx, `
		DELETE FROM products 
		WHERE id = (
			SELECT id FROM products 
			WHERE reception_id = $1 
			ORDER BY created_at DESC 
			LIMIT 1
		)
		RETURNING id`,
		receptionID,
	).Scan(&productID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrNoProducts
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
