package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"moviegiver.kish.net/internal/validator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		var maxBytesError *http.MaxBytesError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-forment JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for filed %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at charect %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not empty")

		case strings.HasPrefix(err.Error(), "json: unknown filed "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown filed ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contains a single JSON value")
	}
	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "invalid value, must be an integer value")
		return defaultValue
	}

	return i
}
