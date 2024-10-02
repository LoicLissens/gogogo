package repositories

import (
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/domain/models/guilde"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeRepository struct {
	conn *pgxpool.Pool
}

var tableName string = tables.GuildeTable{}.GetTableName()

func NewGuildeRepository(connectionPool *pgxpool.Pool) GuildeRepository {
	return GuildeRepository{conn: connectionPool}
}
func (repository *GuildeRepository) GetByUUID(uuid uuid.UUID) (guilde.Guilde, error) {
	entity, err := repository.ScanRow(GetEntityByUuid(repository.conn, uuid, tableName))
	if err != nil {
		return entity, fmt.Errorf("error while fetching entity %w", db.HandleSQLErrors(err, tableName, uuid))
	}
	return entity, nil
}

func (repository *GuildeRepository) Save(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity)
	savedEntity, err := repository.ScanRow(SaveEntity(table, repository.conn))

	if err != nil {
		return savedEntity, fmt.Errorf("error while saving entity %w", db.HandleSQLErrors(err, tableName, entity.Uuid))
	}

	return savedEntity, err
}

func (repository *GuildeRepository) Delete(uuid uuid.UUID) error {
	rowsAffected, err := DeleteEntity(tableName, uuid, repository.conn)
	return HandleSQLDelete(rowsAffected, err, tableName, uuid)
}

func (repository *GuildeRepository) ScanRow(row pgx.Row) (guilde.Guilde, error) {
	entity := guilde.Guilde{}
	err := row.Scan(&entity.Uuid, &entity.Created_at, &entity.Updated_at, &entity.Name, &entity.Img_url, &entity.Page_url, &entity.Exists, &entity.Validated, &entity.Active, &entity.Creation_date)

	if err != nil {
		return entity, err
	}
	return entity, nil
}
func (repository *GuildeRepository) Update(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity)

	updatedentity, err := repository.ScanRow(UpdateEntity(table, repository.conn, entity.Uuid))

	if err != nil {
		return updatedentity, fmt.Errorf("error while updated entity %w", db.HandleSQLErrors(err, tableName, entity.Uuid))
	}

	return updatedentity, err
}
