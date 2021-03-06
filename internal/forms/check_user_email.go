package forms

import (
	"github.com/clevergo/demo/internal/validations"
	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CheckUserEmail struct {
	db    *sqlex.DB
	Email string `json:"email"`
}

func NewCheckUserEmail(db *sqlex.DB) *CheckUserEmail {
	return &CheckUserEmail{db: db}
}

func (f *CheckUserEmail) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Email,
			validation.Required,
			is.Email,
			validation.By(validations.IsUserEmailTaken(f.db)),
		),
	)
}
