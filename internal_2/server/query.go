package server

import "github.com/jackc/pgx/v5"

func NewQuery(q string, params []any) Query {
	return Query{
		query:  q,
		params: params,
		result: make(chan pgx.Row),
	}
}
