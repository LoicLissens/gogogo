package views

type OrderingMethod string
type OrderBy string

const (
	ASC        OrderingMethod = "ASC"
	DESC       OrderingMethod = "DESC"
	CREATED_AT OrderBy        = "created_at"
)

type BaseListViewOpts struct {
	Page           int
	Limit          int
	OrderingMethod OrderingMethod
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
func CheckOrderBy(orderingMethod OrderingMethod, orderBy OrderBy) (OrderingMethod, OrderBy) {
	if orderingMethod == "" {
		orderingMethod = DESC
	}
	if orderBy == "" {
		orderBy = CREATED_AT
	}
	return orderingMethod, orderBy
}
