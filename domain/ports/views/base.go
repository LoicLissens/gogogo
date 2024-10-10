package views

import "time"

type Order string
type OrderBy string

const (
	ASC        Order   = "ASC"
	DESC       Order   = "DESC"
	CREATED_AT OrderBy = "created_at"
)

type BaseListViewOpts struct {
	Page            int
	Limit           int
	Order           Order
	CreatedAtFilter time.Time
	UpdatedAtFilter time.Time
	UuidFilters     []string
}

func CheckPagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return page, limit
}
func CheckOrderBy(order Order, orderBy OrderBy) (Order, OrderBy) {
	if order == "" {
		order = DESC
	}
	if orderBy == "" {
		orderBy = CREATED_AT
	}
	return order, orderBy
}
