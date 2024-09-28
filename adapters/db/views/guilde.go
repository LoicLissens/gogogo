package views

import (
	"context"
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	portView "jiva-guildes/domain/ports/views"
	"jiva-guildes/domain/ports/views/dtos"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tableName string = tables.GuildeTable{}.GetTableName()

type GuildeView struct {
	conn *pgxpool.Pool
}

func NewGuildeView(connectionPool *pgxpool.Pool) portView.GuildeView {
	return GuildeView{conn: connectionPool}
}
func (gv GuildeView) Fetch(uuid uuid.UUID) (dtos.GuildeViewDTO, error) {
	statement := fmt.Sprintf("SELECT * FROM %s WHERE uuid = $1", tableName)
	row := gv.conn.QueryRow(context.Background(), statement, uuid)
	dto := dtos.GuildeViewDTO{}
	err := row.Scan(&dto.Uuid, &dto.Created_at, &dto.Updated_at, &dto.Name, &dto.Img_url, &dto.Page_url)
	if err != nil {
		return dtos.GuildeViewDTO{}, db.HandleSQLErrors(err, tableName, uuid)
	}
	return dto, nil
}

// TODO: add ordering (asc et desc) and filtering
func (gv GuildeView) List(page int, limit int) (dtos.GuildeListViewDTO, error) { //TODO: Will need to add additional check for pagination to add in utils method
	whereClause := "" // For later use /!\ adapt the $ sign of the statement to match the correct query parameter
	params := []interface{}{}
	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	var NbItems int
	err := gv.conn.QueryRow(context.Background(), countStatement, params...).Scan(&NbItems)
	if err != nil {
		return dtos.GuildeListViewDTO{}, err
	}
	statement := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", tableName)
	rows, err := gv.conn.Query(context.Background(), statement, limit, (page-1)*limit)
	if err != nil {
		return dtos.GuildeListViewDTO{}, err
	}
	defer rows.Close()
	dtoList := make([]dtos.GuildeViewDTO, 0)
	for rows.Next() {
		dto := dtos.GuildeViewDTO{}
		err := rows.Scan(&dto.Uuid, &dto.Created_at, &dto.Updated_at, &dto.Name, &dto.Img_url, &dto.Page_url)
		if err != nil {
			return dtos.GuildeListViewDTO{}, err
		}
		dtoList = append(dtoList, dto)
	}
	return dtos.GuildeListViewDTO{Items: dtoList, NbItems: NbItems}, nil

}
