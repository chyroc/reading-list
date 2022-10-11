package transport

import "github.com/go-playground/validator"

var validate = validator.New()

type Inputs struct {
	URL string `validate:"required,url" query:"url"`
}

func (i *Inputs) Validate() error {
	return validate.Struct(i)
}
