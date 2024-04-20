package repositories

import (
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/domain/models/guilde"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeRepository struct {
	conn *pgxpool.Pool
}

func NewGuildeRepository(connectionPool *pgxpool.Pool) GuildeRepository {
	return GuildeRepository{conn: connectionPool}
}
func (repository *GuildeRepository) GetByUUID(uuid uuid.UUID) (guilde.Guilde, error) {
	var tableName string = tables.GuildeTable{}.GetTableName()
	entity, err := repository.ScanRow(GetEntityByUuid(repository.conn, uuid, tableName))
	if err != nil {
		return entity, HandleSQLErrors(err, tableName, uuid)
	}
	return entity, nil
}

func (repository *GuildeRepository) Save(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity.Name, entity.Img_url, entity.Page_url, entity.Uuid, entity.Created_at, entity.Updated_at)
	savedEntity, err := repository.ScanRow(SaveEntity(table, repository.conn))

	if err != nil {
		return savedEntity, HandleSQLErrors(err, table.GetTableName(), entity.Uuid)
	}

	return savedEntity, err
}

func (repository *GuildeRepository) Delete(uuid uuid.UUID) error {
	tableName := tables.GuildeTable{}.GetTableName()
	rowsAffected, err := DeleteEntity(tableName, uuid, repository.conn)
	return HandleSQLDelete(rowsAffected, err, tableName, uuid)
}

func (repository *GuildeRepository) ScanRow(row pgx.Row) (guilde.Guilde, error) {
	entity := guilde.Guilde{}
	err := row.Scan(&entity.Uuid, &entity.Created_at, &entity.Updated_at, &entity.Name, &entity.Img_url, &entity.Page_url)

	if err != nil {
		return entity, err
	}
	return entity, nil
}
