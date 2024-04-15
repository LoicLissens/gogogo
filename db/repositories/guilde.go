package repositories

import (
	"errors"
	"fmt"
	"jiva-guildes/db/tables"
	"jiva-guildes/domain/models/guilde"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeRepository struct {
	conn *pgxpool.Pool
}

func (repository GuildeRepository) GetByUUID(uuid uuid.UUID) (guilde.Guilde, error) {
	var tableName string = tables.GuildeTable{}.GetTableName()
	entity, err := repository.ScanRow(GetEntityByUuid(repository.conn, uuid, tableName))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity, fmt.Errorf("No entity with UUID %s found in table %s", uuid, tableName)
		} else {
			return entity, err
		}
	}
	return entity, nil
}

func (repository GuildeRepository) Save(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity.Name, entity.Img_url, entity.Page_url, entity.Uuid, entity.Created_at, entity.Updated_at)

	savedEntity, err := repository.ScanRow(SaveEntity(table, repository.conn))
	if err != nil {
		log.Fatal(err)
	}
	return savedEntity, err
}

func (repository GuildeRepository) ScanRow(row pgx.Row) (guilde.Guilde, error) {
	entity := guilde.Guilde{}
	err := row.Scan(&entity.Uuid, &entity.Created_at, &entity.Updated_at, &entity.Name, &entity.Img_url, &entity.Page_url)

	if err != nil {
		return entity, err
	}
	return entity, nil
}
