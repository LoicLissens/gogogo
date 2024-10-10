package views

import (
	"jiva-guildes/domain/ports/views/dtos"
	"time"

	"github.com/google/uuid"
)

type GuildeView interface {
	Fetch(uuid uuid.UUID) (dtos.GuildeViewDTO, error)
	List(ListGuildesViewOpts) (dtos.GuildeListViewDTO, error)
}

type OrderByGuilde OrderBy

const (
	OrderByName OrderByGuilde = "name"
)

type ListGuildesViewOpts struct {
	BaseListViewOpts
	OrderBy               OrderByGuilde
	NameFilter            string
	ExistsFilter          *bool
	ValidetedFilter       *bool
	ActiveFilter          *bool
	MinCreationDateFilter time.Time
	MaxCreationDateFilter time.Time
}
