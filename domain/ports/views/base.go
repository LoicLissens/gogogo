package views

import "jiva-guildes/settings"

type OrderingMethod string
type OrderBy string

const (
	ASC        OrderingMethod = "ASC"
	DESC       OrderingMethod = "DESC"
	CREATED_AT OrderBy        = "created_at"
)

func CheckPagination(page, limit int) (int, int) {
	if page < 1 {
		page = settings.AppSettings.DEFAULT_PAGE
	}
	if limit < 1 {
		limit = settings.AppSettings.DEFAULT_PAGE_LIMIT
	}
	return page, limit
}
func CheckOrderBy(orderingMethod OrderingMethod, orderBy OrderBy) (OrderingMethod, OrderBy) {
	if orderingMethod == "" {
		orderingMethod = DESC
	}
	if orderBy == "" {
		orderBy = CREATED_AT
	}
	return orderingMethod, orderBy
}
