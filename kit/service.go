package kit

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

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

func BytesToSlice[O any](c []byte) ([]O, error) {
	var items []O
	err := json.Unmarshal(c, &items)
	if err == nil {
		return items, nil
	}
	var rawItems []map[string]interface{}
	err = json.Unmarshal(c, &rawItems)
	if err != nil {
		return nil, fmt.Errorf("error converting data to model - unmarshal: %w", err)
	}
	for _, raw := range rawItems {
		var item O
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &item,
			TagName:  "json",
		}
		decoder, err := mapstructure.NewDecoder(cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating mapstructure decoder: %w", err)
		}
		if err := decoder.Decode(raw); err != nil {
			return nil, fmt.Errorf("error converting data to model - mapstructure: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func SliceToBytes[O any](c []O) ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("error converting struct to bytes: %w", err)
	}
	return b, nil
}

func StructToMap[O any](obj O) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("error serializando la estructura: %w", err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo JSON a mapa: %w", err)
	}
	return result, nil
}

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
