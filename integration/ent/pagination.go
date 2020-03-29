// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/jmoiron/sqlx/reflectx"
)

var mapper = reflectx.NewMapper("json")

// Predicate creates a predicate
type Predicate = func(s *sql.Selector)

// CursorPosition represets a cursor position
type CursorPosition struct {
	Column    string
	Direction string
	Value     interface{}
}

// Cursor represents the cursor
type Cursor struct {
	positions []*CursorPosition
}

// DecodeCursor decodes a cursor from its base-64 string representation.
func DecodeCursor(order, token string) (*Cursor, error) {
	cursor := &Cursor{}

	if err := cursor.positionsAt(order); err != nil {
		return nil, err
	}

	if err := cursor.valuesAt(token); err != nil {
		return nil, err
	}

	return cursor, nil
}

// String returns a base-64 string representation of a cursor.
func (c *Cursor) String() string {
	count := len(c.positions)

	if count == 0 {
		return ""
	}

	values := make([]interface{}, count)

	for index, position := range c.positions {
		values[index] = position.Value
	}

	data, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// Next returns the next cursor
func (c *Cursor) Next(input interface{}) *Cursor {
	var (
		next   = Cursor{}
		source = reflect.ValueOf(input)
	)

	//TODO: do not use reflection
	if source.Type().Kind() == reflect.Slice {
		if source.Len() == 0 {
			return &next
		}

		source = source.Index(source.Len() - 1)
	}

	for _, position := range c.positions {
		index := &CursorPosition{
			Column:    position.Column,
			Direction: position.Direction,
			Value:     mapper.FieldByName(source, position.Column).Interface(),
		}

		next.positions = append(next.positions, index)
	}

	return &next
}

func (c *Cursor) positionsAt(order string) error {
	const (
		separator = ","
		asc       = "+"
		desc      = "-"
	)

	for _, field := range strings.Split(order, separator) {
		field = strings.TrimSpace(field)

		if field == "" {
			continue
		}

		position := &CursorPosition{
			Column:    field,
			Direction: asc,
		}

		switch {
		case strings.HasPrefix(field, asc):
			position = &CursorPosition{
				Column:    field[1:],
				Direction: asc,
			}
		case strings.HasPrefix(field, desc):
			position = &CursorPosition{
				Column:    field[1:],
				Direction: desc,
			}
		}

		// TODO: validate column name
		c.positions = append(c.positions, position)
	}

	return nil
}

func (c *Cursor) valuesAt(token string) error {
	values := []interface{}{}

	if token == "" {
		return nil
	}

	if n := len(token) % 4; n != 0 {
		token += strings.Repeat("=", 4-n)
	}

	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	for index, position := range c.positions {
		if index >= len(values) {
			return fmt.Errorf("ent: invalid pagination cursor")
		}

		position.Value = values[index]
	}

	return nil
}

// Pagination represents the pagination
type Pagination struct {
	Predicate Predicate
	Order     []Predicate
	Limit     int
}

// NewPagination creates a new pagination
func NewPagination() *Pagination {
	return &Pagination{}
}

// WithLimit returns the pagination with given limit
func (pg Pagination) WithLimit(limit int) *Pagination {
	pg.Limit = limit
	return &pg
}

// WithCursor set the pagination cursor
func (pg Pagination) WithCursor(cursor *Cursor) *Pagination {
	pg.Predicate = pg.apply(cursor.positions)

	for _, position := range cursor.positions {
		switch position.Direction {
		case "+":
			pg.Order = append(pg.Order, asc(position.Column))
		case "-":
			pg.Order = append(pg.Order, desc(position.Column))
		}
	}

	return &pg
}

func (pg *Pagination) apply(positions []*CursorPosition) Predicate {
	var (
		predicate        Predicate = func(*sql.Selector) {}
		predicateCompare Predicate = func(*sql.Selector) {}
		predicateEqual   Predicate = func(*sql.Selector) {}
	)

	if len(positions) == 0 {
		return predicate
	}

	position := positions[0]

	if position.Value != nil {
		predicateEqual = eq(position.Column, position.Value)

		switch position.Direction {
		case "+":
			predicateCompare = gt(position.Column, position.Value)
		case "-":
			predicateCompare = lt(position.Column, position.Value)
		default:
			predicateCompare = gt(position.Column, position.Value)
		}
	}

	positions = positions[1:]
	predicate = predicateCompare

	if len(positions) > 0 {
		predicate = or(predicateCompare,
			and(predicateEqual, pg.apply(positions)))
	}

	return predicate
}

func (pq *ProductQuery) Paginate(page *Pagination) *ProductQuery {
	pq = pq.Where(page.Predicate)
	pq = pq.Limit(page.Limit)

	for _, predicate := range page.Order {
		pq = pq.Order(predicate)
	}

	return pq
}

func eq(field string, value interface{}) Predicate {
	return func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(field), value))
	}
}

func gt(field string, value interface{}) Predicate {
	return func(s *sql.Selector) {
		s.Where(sql.GT(s.C(field), value))
	}
}

func lt(field string, value interface{}) Predicate {
	return func(s *sql.Selector) {
		s.Where(sql.LT(s.C(field), value))
	}
}

func and(predicates ...Predicate) Predicate {
	return func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	}
}

func or(predicates ...Predicate) Predicate {
	return func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	}
}

func asc(fields ...string) Predicate {
	return func(s *sql.Selector) {
		for _, f := range fields {
			s.OrderBy(sql.Asc(f))
		}
	}
}

func desc(fields ...string) Predicate {
	return func(s *sql.Selector) {
		for _, f := range fields {
			s.OrderBy(sql.Desc(f))
		}
	}
}
