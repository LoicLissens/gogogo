package tables

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeTable struct {
	BaseModelTable
	Name     string `db:"name" sql_properties:"VARCHAR(255) NOT NULL"`
	Img_url  string `db:"img_url" sql_properties:"VARCHAR(255)"`
	Page_url string `db:"page_url" sql_properties:"VARCHAR(255)"`
}

func (table GuildeTable) GetTableName() string {
	return "guildes"
}
func (table GuildeTable) CreateTable(conn *pgxpool.Pool) {
	CreateTable(conn, table)
}
func (table GuildeTable) DropTable(conn *pgxpool.Pool) {
	DropTable(conn, table)
}
func NewGuildeTable(name, img_url, page_url string, Uuid uuid.UUID, Created_at, Updated_at time.Time) GuildeTable {
	return GuildeTable{
		BaseModelTable: BaseModelTable{
			Uuid:       Uuid,
			Created_at: Created_at,
			Updated_at: Updated_at,
		},
		Name:     name,
		Img_url:  img_url,
		Page_url: page_url,
	}
}
