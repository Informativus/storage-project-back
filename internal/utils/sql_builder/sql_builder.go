package sql_builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

func InsertArgs[T any](model T) (cols []string, vals []any, placeholders []string, err error) {
	val := reflect.ValueOf(model)
	typ := reflect.TypeOf(model)

	if val.Kind() != reflect.Struct {
		err = errors.New("model must be a struct")
		log.Error().Msg(err.Error())
		return nil, nil, nil, err
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		if field.Tag.Get("json") == "" {
			err = errors.New("model must be a struct")
			log.Error().Msg(err.Error())
			return nil, nil, nil, err
		}

		if field.Tag.Get("json") == "-" {
			continue
		}

		cols = append(cols, field.Tag.Get("json"))
		vals = append(vals, val.Field(i).Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(vals)))
	}

	return cols, vals, placeholders, nil
}

func BuildInsertQuery(table string, cols []string, placeholders []string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
}
