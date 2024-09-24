package views

import (
	"context"
	"fmt"
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
func (gv GuildeView) Fetch(uuid uuid.UUID) dtos.GuildeViewDTO {
	statement := fmt.Sprintf("SELECT * FROM %s WHERE uuid = $1", tableName)
	row := gv.conn.QueryRow(context.Background(), statement, uuid)
	dto := dtos.GuildeViewDTO{}
	err := row.Scan(&dto)
	if err != nil {
		panic(err)
	}
	return dto
}
func (gv GuildeView) List(page int, limit int) {
	whereClause := "" // For later use /!\ adapt the $ sign of the statement to match the correct query parameter
	params := []interface{}{}
	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	var total int
	err := gv.conn.QueryRow(context.Background(), countStatement, params...).Scan(&total)
	if err != nil {
		panic(err)
	}
	statement := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", tableName)
	rows, err := gv.conn.Query(context.Background(), statement, limit, (page-1)*limit)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		dto := dtos.GuildeViewDTO{}
		err := rows.Scan(&dto)
		if err != nil {
			panic(err)
		}
		fmt.Println(dto)
	}
}
