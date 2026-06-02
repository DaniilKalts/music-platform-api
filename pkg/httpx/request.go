package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const maxRequestBodyBytes = 1 << 20

var validate = validator.New()

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			WriteError(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return false
		}

		WriteError(w, http.StatusBadRequest, "invalid request body")
		return false
	}

	if err := validate.Struct(v); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			msgs := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				field := e.Field()
				if jsonName := jsonFieldName(e.StructField(), v); jsonName != "" {
					field = jsonName
				}
				msgs = append(msgs, fmt.Sprintf("field %s is invalid", field))
			}
			WriteError(w, http.StatusBadRequest, strings.Join(msgs, ", "))
			return false
		}

		WriteError(w, http.StatusBadRequest, "validation failed")
		return false
	}

	return true
}

func jsonFieldName(field string, v any) string {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return ""
	}

	structField, ok := typeOf.FieldByName(field)
	if !ok {
		return ""
	}

	name := strings.Split(structField.Tag.Get("json"), ",")[0]
	if name == "-" {
		return ""
	}

	return name
}
