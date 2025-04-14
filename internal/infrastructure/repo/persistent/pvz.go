package persistent

import (
	"PVZ-avito-tech/internal/controller/http/dto"
	"PVZ-avito-tech/internal/entity"
	"PVZ-avito-tech/internal/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type PVZRepo struct {
	*postgres.Postgres
}

func NewPVZRepo(pg *postgres.Postgres) *PVZRepo {
	return &PVZRepo{pg}
}

func (r *PVZRepo) Create(ctx context.Context, pvz *entity.PVZ) error {
	columns := []string{"city"}
	values := []interface{}{pvz.City}

	if pvz.ID != nil && *pvz.ID != uuid.Nil {
		columns = append(columns, "id")
		values = append(values, *pvz.ID)
	}

	if pvz.RegistrationDate != nil && !pvz.RegistrationDate.IsZero() {
		columns = append(columns, "created_at")
		values = append(values, *pvz.RegistrationDate)
	}

	builder := r.Builder.
		Insert("pvz").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	row := r.Pool.QueryRow(ctx, query, args...)
	var id uuid.UUID
	var created time.Time
	if err := row.Scan(&id, &created); err != nil {
		return entity.ErrCreatePVZ
	}

	pvz.ID = &id
	pvz.RegistrationDate = &created

	return nil
}

func (r *PVZRepo) GetPVZWithReceptions(
	ctx context.Context,
	filter dto.ReceptionFilter,
) (*[]dto.PVZInfo, error) {
	baseQuery := r.Builder.
		Select(
			"pvz.id AS pvz_id",
			"pvz.city AS pvz_city",
			"pvz.created_at AS pvz_created_at",
			"r.id AS reception_id",
			"r.created_at AS reception_created_at",
			"r.status AS reception_status",
			"p.id AS product_id",
			"p.type AS product_type",
			"p.created_at AS product_created_at",
		).
		From("pvz").
		LeftJoin("receptions r ON pvz.id = r.pvz_id").
		LeftJoin("products p ON r.id = p.reception_id").
		Where(sq.And{
			sq.Or{
				sq.GtOrEq{"r.created_at": filter.StartDate},
				sq.Eq{"r.created_at": nil},
			},
			sq.Or{
				sq.LtOrEq{"r.created_at": filter.EndDate},
				sq.Eq{"r.created_at": nil},
			},
		}).
		OrderBy("pvz.created_at DESC", "r.created_at DESC", "p.created_at DESC").
		Limit(uint64(filter.Limit)).
		Offset(uint64((filter.Page - 1) * filter.Limit)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	pvzMap := make(map[uuid.UUID]*dto.PVZInfo)
	receptionMap := make(map[uuid.UUID]*dto.ReceptionGroup)

	for rows.Next() {
		var (
			pvzID           uuid.UUID
			pvzCity         entity.City
			pvzCreatedAt    time.Time
			receptionID     uuid.NullUUID
			receptionDate   pq.NullTime
			receptionStatus sql.NullString
			productID       uuid.NullUUID
			productType     sql.NullString
			productDate     pq.NullTime
		)

		err := rows.Scan(
			&pvzID,
			&pvzCity,
			&pvzCreatedAt,
			&receptionID,
			&receptionDate,
			&receptionStatus,
			&productID,
			&productType,
			&productDate,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		// Обработка PVZ
		if _, exists := pvzMap[pvzID]; !exists {
			pvzMap[pvzID] = &dto.PVZInfo{
				PVZ: dto.PVZWithReceptions{
					ID:               pvzID,
					City:             pvzCity,
					RegistrationDate: pvzCreatedAt,
				},
				Receptions: []dto.ReceptionGroup{},
			}
		}

		// Обработка Reception
		if receptionID.Valid {
			receptionKey := receptionID.UUID
			if _, exists := receptionMap[receptionKey]; !exists {
				receptionMap[receptionKey] = &dto.ReceptionGroup{
					Reception: dto.ReceptionWithProducts{
						ID:       receptionID.UUID,
						DateTime: receptionDate.Time,
						PVZID:    pvzID,
						Status:   entity.ReceptionsStatus(receptionStatus.String),
					},
					Products: []dto.ProductDTO{},
				}
				pvzMap[pvzID].Receptions = append(
					pvzMap[pvzID].Receptions,
					*receptionMap[receptionKey],
				)
			}

			// Обработка Product
			if productID.Valid {
				product := dto.ProductDTO{
					ID:          productID.UUID,
					DateTime:    productDate.Time,
					Type:        entity.City(productType.String), // Возможна ошибка типа!
					ReceptionID: receptionID.UUID,
				}
				receptionMap[receptionKey].Products = append(
					receptionMap[receptionKey].Products,
					product,
				)
			}
		}
	}

	result := make([]dto.PVZInfo, 0, len(pvzMap))
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return &result, rows.Err()
}
