package persistent

import (
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ReceptionRepo struct {
	*postgres.Postgres
}

func NewReceptionRepo(pool *postgres.Postgres) *ReceptionRepo {
	return &ReceptionRepo{pool}
}

func (r *ReceptionRepo) CreateReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	query := `
		INSERT INTO receptions (pvz_id, status)
		VALUES ($1, $2)
		RETURNING id, pvz_id, status, created_at
	`

	var reception entity.Reception
	err := r.Pool.QueryRow(ctx, query, pvzID, entity.InProgressStatus).Scan(
		&reception.ID,
		&reception.PVZID,
		&reception.Status,
		&reception.DateTime,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				return nil, entity.ErrReceptionConflict
			case pgErr.Code == "23503":
				return nil, entity.ErrPVZNotFound
			}
		}
		return nil, fmt.Errorf("failed to create reception: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) CloseActiveReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	query := `
		UPDATE receptions 
		SET status = $1 
		WHERE pvz_id = $2 
		AND status = $3
		RETURNING id, pvz_id, status, created_at
	`

	var reception entity.Reception
	err := r.Pool.QueryRow(
		ctx,
		query,
		entity.CloseStatus,
		pvzID,
		entity.InProgressStatus,
	).Scan(
		&reception.ID,
		&reception.PVZID,
		&reception.Status,
		&reception.DateTime,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrNoActiveReception
		}
		return nil, fmt.Errorf("failed to close reception: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) CloseActiveReceptionWithTx(
	ctx context.Context,
	tx pgx.Tx,
	pvzID uuid.UUID,
) (*entity.Reception, error) {
	query := `
		UPDATE receptions 
		SET status = $1 
		WHERE pvz_id = $2 
		AND status = $3
		RETURNING id, pvz_id, status, created_at
	`

	var reception entity.Reception
	err := tx.QueryRow(
		ctx,
		query,
		entity.CloseStatus,
		pvzID,
		entity.InProgressStatus,
	).Scan(
		&reception.ID,
		&reception.PVZID,
		&reception.Status,
		&reception.DateTime,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrNoActiveReception
		}
		return nil, fmt.Errorf("failed to close reception: %w", err)
	}

	return &reception, nil
}
