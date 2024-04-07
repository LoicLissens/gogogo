package models

type GuildeTable struct {
	BaseModelTable
	Name     string `db:"name" sql_properties:"VARCHAR(255) NOT NULL"`
	Img_url  string `db:"img_url" sql_properties:"VARCHAR(255)"`
	Page_url string `db:"page_url" sql_properties:"VARCHAR(255)"`
}

func (g GuildeTable) getTableName() string {
	return "guildes"
}
