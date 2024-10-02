package tables

import (
	"jiva-guildes/domain/models/guilde"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeTable struct {
	BaseModelTable
	Name          string     `db:"name" sql_properties:"VARCHAR(255) NOT NULL"`
	Img_url       string     `db:"img_url" sql_properties:"VARCHAR(255)  NOT NULL"`
	Page_url      string     `db:"page_url" sql_properties:"VARCHAR(255)  NOT NULL"`
	Exists        bool       `db:"exists" sql_properties:"BOOLEAN NOT NULL"`
	Validated     bool       `db:"validated" sql_properties:"BOOLEAN NOT NULL"`
	Active        *bool      `db:"active" sql_properties:"BOOLEAN"`
	Creation_date *time.Time `db:"creation_date" sql_properties:"TIMESTAMP"`
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
func NewGuildeTable(g guilde.Guilde) GuildeTable {
	return GuildeTable{
		BaseModelTable: BaseModelTable{
			Uuid:       g.Uuid,
			Created_at: g.Created_at,
			Updated_at: g.Updated_at,
		},
		Name:          g.Name,
		Img_url:       g.Img_url,
		Page_url:      g.Page_url,
		Exists:        g.Exists,
		Active:        g.Active,
		Validated:     g.Validated,
		Creation_date: g.Creation_date,
	}
}
