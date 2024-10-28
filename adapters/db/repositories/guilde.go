package repositories

import (
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/domain/models/guilde"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GuildeRepository struct {
	db db.PsqlDB
}

var tableName string = tables.GuildeTable{}.GetTableName()

func NewGuildeRepository(db db.PsqlDB) GuildeRepository {
	return GuildeRepository{db: db}
}
func (repository *GuildeRepository) GetByUUID(uuid uuid.UUID) (guilde.Guilde, error) {
	entity, err := repository.ScanRow(GetEntityByUuid(repository.db, uuid, tableName))
	if err != nil {
		return entity, fmt.Errorf("error while fetching entity %w", db.HandleSQLErrors(err, tableName, uuid))
	}
	return entity, nil
}
func (repository *GuildeRepository) GetAll() ([]guilde.Guilde, error) {
	rows := GetAllEntities(repository.db, tableName)
	defer rows.Close()
	entities, err := repository.ScanRows(rows)
	if err != nil {
		return entities, fmt.Errorf("error while fetching entities %w", err)
	}
	return entities, nil
}
func (repository *GuildeRepository) Save(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity)
	savedEntity, err := repository.ScanRow(SaveEntity(table, repository.db))

	if err != nil {
		return savedEntity, fmt.Errorf("error while saving entity %w", db.HandleSQLErrors(err, tableName, entity.Uuid))
	}

	return savedEntity, err
}

func (repository *GuildeRepository) Delete(uuid uuid.UUID) error {
	rowsAffected, err := DeleteEntity(tableName, uuid, repository.db)
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

func (repository *GuildeRepository) ScanRows(rows pgx.Rows) ([]guilde.Guilde, error) {
	var entities []guilde.Guilde

	for rows.Next() {
		entity := guilde.Guilde{}
		err := rows.Scan(&entity.Uuid, &entity.Created_at, &entity.Updated_at, &entity.Name, &entity.Img_url, &entity.Page_url, &entity.Exists, &entity.Validated, &entity.Active, &entity.Creation_date)
		if err != nil {
			return entities, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
func (repository *GuildeRepository) Update(entity guilde.Guilde) (guilde.Guilde, error) {
	table := tables.NewGuildeTable(entity)
	table.Updated_at = time.Now().UTC() //TODO: lame
	updatedentity, err := repository.ScanRow(UpdateEntity(table, repository.db, entity.Uuid))

	if err != nil {
		return updatedentity, fmt.Errorf("error while updated entity %w", db.HandleSQLErrors(err, tableName, entity.Uuid))
	}

	return updatedentity, err
}
