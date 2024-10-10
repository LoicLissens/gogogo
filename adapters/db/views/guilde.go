package views

import (
	"context"
	"fmt"
	"jiva-guildes/adapters/db"
	"jiva-guildes/adapters/db/tables"
	"jiva-guildes/domain/ports/views"
	portView "jiva-guildes/domain/ports/views"
	"jiva-guildes/domain/ports/views/dtos"
	"log"
	"strings"

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
	err := row.Scan(&dto.Uuid, &dto.Created_at, &dto.Updated_at, &dto.Name, &dto.Img_url, &dto.Page_url, &dto.Exists, &dto.Validated, &dto.Active, &dto.Creation_date)
	if err != nil {
		return dtos.GuildeViewDTO{}, db.HandleSQLErrors(err, tableName, uuid)
	}
	return dto, nil
}

func (gv GuildeView) List(opts views.ListGuildesViewOpts) (dtos.GuildeListViewDTO, error) {
	page, limit := portView.CheckPagination(opts.Page, opts.Limit)
	positionalParams := 0
	var sb strings.Builder
	for _, filter := range opts.FilterBy {
		positionalParams++
		var whereClause string
		if positionalParams == 1 {
			whereClause = fmt.Sprintf("WHERE %s = $%d", filter, positionalParams)
		} else {
			whereClause = fmt.Sprintf(" AND %s = $%d ", filter, positionalParams)
		}
		sb.WriteString(whereClause)
	}
	whereClause := sb.String()
	// Retrieve count
	params := []interface{}{}
	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	var NbItems int
	err := gv.conn.QueryRow(context.Background(), countStatement, params...).Scan(&NbItems)
	if err != nil {
		return dtos.GuildeListViewDTO{}, err
	}
	// Retrieve Items
	order, orderBy := portView.CheckOrderBy(opts.Order, portView.OrderBy(opts.OrderBy))

	statement := fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2", tableName, orderBy, order)
	rows, err := gv.conn.Query(context.Background(), statement, limit, (page-1)*limit)
	if err != nil {
		return dtos.GuildeListViewDTO{}, err
	}

	defer rows.Close()
	dtoList := make([]dtos.GuildeViewDTO, 0)
	for rows.Next() {
		dto := dtos.GuildeViewDTO{}
		err := rows.Scan(&dto.Uuid, &dto.Created_at, &dto.Updated_at, &dto.Name, &dto.Img_url, &dto.Page_url, &dto.Exists, &dto.Validated, &dto.Active, &dto.Creation_date)
		if err != nil {
			return dtos.GuildeListViewDTO{}, err
		}
		dtoList = append(dtoList, dto)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return dtos.GuildeListViewDTO{Items: dtoList, NbItems: NbItems}, nil
}
