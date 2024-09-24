package views

import (
	portView "jiva-guildes/domain/ports/views"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ViewsManager struct {
	conn *pgxpool.Pool
}

func NewViewsManager(connectionPool *pgxpool.Pool) portView.ViewsManager {
	return ViewsManager{conn: connectionPool}
}
func (vm ViewsManager) Guilde() portView.GuildeView {
	return NewGuildeView(vm.conn)
}
