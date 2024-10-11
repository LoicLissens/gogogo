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
	whereClause, positionalParam, params := writeWhereClase(opts)
	// Retrieve count
	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereClause)
	var NbItems int
	err := gv.conn.QueryRow(context.Background(), countStatement, params...).Scan(&NbItems)
	if err != nil {
		return dtos.GuildeListViewDTO{}, fmt.Errorf("error while retrieving the count of guilde %w", err)
	}

	// Retrieve Items
	order, orderBy := portView.CheckOrderBy(opts.OrderingMethod, portView.OrderBy(opts.OrderBy))
	params = append(params, limit, (page-1)*limit)
	statement := fmt.Sprintf("SELECT * FROM %s %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		tableName, whereClause, orderBy, order, positionalParam+1, positionalParam+2)
	rows, err := gv.conn.Query(context.Background(), statement, params...)
	if err != nil {
		return dtos.GuildeListViewDTO{}, fmt.Errorf("error while retrieving the guildes %w", err)
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

func writeWhereClase(opts views.ListGuildesViewOpts) (string, int, []interface{}) {
	positionalParam := 0
	var sb strings.Builder
	var params []interface{}

	if opts.Name != "" {
		positionalParam++
		params = append(params, "%"+opts.Name+"%")
		sb.WriteString(fmt.Sprintf("%s LOWER(name) LIKE LOWER($%d)", getFilterKeyWord(positionalParam), positionalParam))
	}
	if opts.Exists != nil {
		positionalParam++
		params = append(params, *opts.Exists)
		sb.WriteString(fmt.Sprintf("%s exists = $%d", getFilterKeyWord(positionalParam), positionalParam))
	}
	if opts.Validated != nil {
		positionalParam++
		params = append(params, *opts.Validated)
		sb.WriteString(fmt.Sprintf("%s validated = $%d", getFilterKeyWord(positionalParam), positionalParam))
	}
	if opts.Active != nil {
		positionalParam++
		params = append(params, *opts.Active)
		sb.WriteString(fmt.Sprintf("%s active = $%d", getFilterKeyWord(positionalParam), positionalParam))
	}
	if !opts.CreationDateSince.IsZero() {
		positionalParam++
		params = append(params, opts.CreationDateSince)
		sb.WriteString(fmt.Sprintf("%s creation_date >= $%d", getFilterKeyWord(positionalParam), positionalParam))
	}
	if !opts.CreationDateUntil.IsZero() {
		positionalParam++
		params = append(params, opts.CreationDateUntil)
		sb.WriteString(fmt.Sprintf("%s creation_date <= $%d", getFilterKeyWord(positionalParam), positionalParam))
	}

	whereClause := sb.String()
	return whereClause, positionalParam, params
}
func getFilterKeyWord(positionalParam int) string {
	if positionalParam == 1 {
		return "WHERE"
	} else {
		return " AND"
	}
}
