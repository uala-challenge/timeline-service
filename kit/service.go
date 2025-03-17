package kit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var validate = validator.New()

func (t *Request) Validate() error {
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

func BytesToModel[O any](c []byte) (O, error) {
	h := *new(O)
	e := map[string]interface{}{}
	err := json.Unmarshal(c, &e)
	if err != nil {
		return h, fmt.Errorf("error converting data to model - unmarshal: %w", err)
	}
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err = decoder.Decode(e)
	if err != nil {
		return h, fmt.Errorf("error converting data to model - mapstructure: %w", err)
	}
	return h, nil
}
