package kit

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (t *TweetRequest) Validate() error {
	return prepareErrorResponse(validate.Struct(t))
}

func prepareErrorResponse(err error) error {
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				return maxResponse(e)
			}
		}
		return err
	}
	return nil
}

func maxResponse(e validator.FieldError) error {
	switch e.Tag() {
	case "max":
		return fmt.Errorf("el campo %s superar los %s caracteres", e.Field(), e.Param())
	case "required":
		return fmt.Errorf("el campo %s es requerido", e.Field())
	default:
		return fmt.Errorf("campo '%s' falló en la validación", e.Field())
	}
}
