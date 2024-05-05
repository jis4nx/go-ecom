package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (a *App) ReadJson(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// limiting maximum request body size to 1MB
	maxBytes := 10_00_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	jsonDec := json.NewDecoder(r.Body)
	jsonDec.DisallowUnknownFields()

	var (
		ErrSyntaxJson     *json.SyntaxError
		ErrUnmarshalType  *json.UnmarshalTypeError
		ErrInvalidMarshal *json.InvalidUnmarshalError
	)

	err := jsonDec.Decode(dst)
	if err != nil {
		switch {

		case errors.As(err, &ErrSyntaxJson):
			return fmt.Errorf("Body contains badly-formatted JSON (at character %d)", ErrSyntaxJson.Offset)

		case errors.As(err, &ErrUnmarshalType):
			if ErrUnmarshalType.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON type (at character %d)", ErrUnmarshalType.Offset)
			}

			return fmt.Errorf("Body contains incorrect JSON field of %q", ErrUnmarshalType.Field)
		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return errors.New("Request body must not be larger than 1MB")

		case errors.As(err, &ErrInvalidMarshal):
			panic(err)

		default:
			return err
		}
	}
	return nil
}
