package guilde

import (
	"time"

	"jiva-guildes/domain/commands"
	customerrors "jiva-guildes/domain/custom_errors"
	"jiva-guildes/domain/models"

	"github.com/google/uuid"
)

type Guilde struct {
	models.BaseModel
	Name          string
	Img_url       string
	Page_url      string
	Exists        bool
	Validated     bool
	Active        *bool
	Creation_date *time.Time
}

type GuildeOptions struct {
	Name          string
	Img_url       string
	Page_url      string
	Exists        bool
	Validated     bool
	Active        *bool
	Creation_date *time.Time
}

func CreateFromCommand(cmd commands.CreateGuildeCommand, validated bool) (*Guilde, error) {
	opts := GuildeOptions{
		Name:          cmd.Name,
		Img_url:       cmd.Img_url,
		Page_url:      cmd.Page_url,
		Exists:        cmd.Exists,
		Validated:     validated,
		Active:        cmd.Active,
		Creation_date: cmd.Creation_date,
	}
	return New(opts)
}
func New(opts GuildeOptions) (*Guilde, error) {

	if !opts.Exists && (opts.Active != nil && *opts.Active) {
		return &Guilde{}, customerrors.NewValueError("A guilde can't be active if it doesn't exist")
	}

	g := Guilde{
		BaseModel: models.BaseModel{
			Uuid:       uuid.New(),
			Created_at: time.Now().UTC(),
			Updated_at: time.Now().UTC(),
		},
		Name:          opts.Name,
		Img_url:       opts.Img_url,
		Page_url:      opts.Page_url,
		Exists:        opts.Exists,
		Validated:     opts.Validated,
		Active:        opts.Active,
		Creation_date: opts.Creation_date,
	}
	return &g, nil
}
