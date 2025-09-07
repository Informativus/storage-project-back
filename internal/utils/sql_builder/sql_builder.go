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

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		err = errors.New("model must be a struct")
		log.Error().Msg(err.Error())
		return nil, nil, nil, err
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		if field.Tag.Get("json") == "" {
			err = errors.New("the model must be described correctly")
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
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(cols, ", "),
	)
}

func GetStructCols[T any](model T) (cols []string, err error) {
	typ := reflect.TypeOf(model)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Tag.Get("json") == "" {
			err = errors.New("the model must be described correctly")
			log.Error().Msg(err.Error())
			return nil, err
		}

		if field.Tag.Get("json") == "-" {
			continue
		}

		cols = append(cols, field.Tag.Get("json"))
	}

	return cols, nil
}

func BuildSelectQuery(table string, cols []string, whereExpression *string) string {
	var query string
	if whereExpression == nil {
		query = fmt.Sprintf("SELECT %s FROM %s",
			strings.Join(cols, ", "),
			table,
		)
		return query
	}

	query = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(cols, ", "),
		table,
		*whereExpression,
	)
	return query
}

func BuildUpdateQueryReturn(table string, setClauses []string, whereExpression string, returning []string) string {
	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s RETURNING %s",
		table,
		strings.Join(setClauses, ", "),
		whereExpression,
		strings.Join(returning, ", "),
	)
}

func BuildUpdateQuery(table string, fildsToUpdate []string, whereExpression string) string {
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(fildsToUpdate, ", "), whereExpression)
}

func BuildDeleteQuery(table string, whereExpression string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereExpression)
}

func BuildSetClauses(fildsToUpdate []string) []string {
	var setClauses []string
	for _, fild := range fildsToUpdate {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", fild, len(setClauses)+1))
	}
	return setClauses
}
