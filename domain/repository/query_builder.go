package repository

import (
	"strconv"

	gotils "github.com/savsgio/gotils/strconv"

	"graphql-project/core"
	. "graphql-project/interface/core"
	"graphql-project/interface/model"
)

type updateQueryBuilder struct {
	query []byte
	args  []any
}

type insertQueryBuilder updateQueryBuilder

type selectQueryBuilder updateQueryBuilder

type selectQuery interface {
	Build(entity model.Entity)
	Query() string
	Args() []any
}

func Update(tableName string) *updateQueryBuilder {
	query := make([]byte, 0, 256)
	query = core.AppendStrings(query, `UPDATE "`, tableName, `" SET `)
	return &updateQueryBuilder{query, make([]any, 0, 32)}
}

func (q *updateQueryBuilder) Value(name string, value any) {
	q.args = append(q.args, value)
	i := len(q.args)
	if i > 1 {
		q.query = append(q.query, ',')
	}
	q.query = append(q.query, '"')
	q.query = core.AppendStrings(q.query, name, `"=$`)
	q.query = strconv.AppendInt(q.query, int64(i), 10)
}

func (q *updateQueryBuilder) Where(name string, value any) *updateQueryBuilder {
	q.args = append(q.args, value)
	q.query = core.AppendStrings(q.query, ` WHERE "`, name, `"=$`)
	q.query = strconv.AppendInt(q.query, int64(len(q.args)), 10)
	return q
}

func (q *updateQueryBuilder) And(name string, value any) *updateQueryBuilder {
	q.args = append(q.args, value)
	q.query = core.AppendStrings(q.query, ` AND "`, name, `"=$`)
	q.query = strconv.AppendInt(q.query, int64(len(q.args)), 10)
	return q
}

func (q *updateQueryBuilder) Query(fields Iterator) string {
	q.query = appendQuoted(" RETURNING ", q.query, fields)
	return gotils.B2S(q.query)
}

func (q *updateQueryBuilder) Args() []any {
	return q.args
}

func InsertInto(tableName string) *insertQueryBuilder {
	query := make([]byte, 0, 256)
	query = core.AppendStrings(query, `INSERT INTO "`, tableName, `"(`)
	return &insertQueryBuilder{query, make([]any, 0, 32)}
}

func (q *insertQueryBuilder) Value(name string, value any) {
	q.args = append(q.args, value)
	if len(q.args) > 1 {
		q.query = append(q.query, ',')
	}
	q.query = append(q.query, '"')
	q.query = append(q.query, name...)
	q.query = append(q.query, '"')
}

func (q *insertQueryBuilder) Query(fields Iterator) string {
	q.query = append(q.query, ") VALUES("...)
	for i := 1; i <= len(q.args); i++ {
		if i > 1 {
			q.query = append(q.query, ',')
		}
		q.query = append(q.query, '$')
		q.query = strconv.AppendInt(q.query, int64(i), 10)
	}
	q.query = appendQuoted(") RETURNING ", q.query, fields)
	return gotils.B2S(q.query)
}

func (q *insertQueryBuilder) Args() []any {
	return q.args
}

func appendQuoted(p string, b []byte, iter Iterator) []byte {
	var i uint32 = 0
	for iter.Next() {
		if i > 0 {
			b = append(b, ',')
		} else {
			b = append(b, p...)
		}
		b = append(b, '"')
		b = append(b, iter.Get()...)
		b = append(b, '"')
		i++
	}
	return b
}

func Select(fields Iterator) *selectQueryBuilder {
	query := make([]byte, 0, 256)
	query = appendQuoted("SELECT ", query, fields)
	return &selectQueryBuilder{query, make([]any, 0, 16)}
}

func (q *selectQueryBuilder) From(tableName string) *selectQueryBuilder {
	q.query = core.AppendStrings(q.query, ` FROM "`, tableName)
	q.query = append(q.query, '"')
	return q
}

func (q *selectQueryBuilder) Where(name string, value any) *selectQueryBuilder {
	q.query = append(q.query, ` WHERE `...)
	if value == nil {
		q.query = append(q.query, name...)
	} else {
		q.args = append(q.args, value)
		q.query = append(q.query, '"')
		q.query = append(q.query, name...)
		q.query = append(q.query, `"=$`...)
		q.query = strconv.AppendInt(q.query, int64(len(q.args)), 10)
	}
	return q
}

func (q *selectQueryBuilder) And(name string, value any) *selectQueryBuilder {
	q.query = append(q.query, ` AND `...)
	if value == nil {
		q.query = append(q.query, name...)
	} else {
		q.args = append(q.args, value)
		q.query = append(q.query, '"')
		q.query = core.AppendStrings(q.query, name, `"=$`)
		q.query = strconv.AppendInt(q.query, int64(len(q.args)), 10)
	}
	return q
}

func (q *selectQueryBuilder) OrderBy(name string, desc bool) *selectQueryBuilder {
	q.query = core.AppendStrings(q.query, ` ORDER BY "`, name)
	q.query = append(q.query, '"')
	if desc {
		q.query = append(q.query, " DESC"...)
	}
	return q
}

func (q *selectQueryBuilder) Offset(offset int32) *selectQueryBuilder {
	if offset > 0 {
		q.query = append(q.query, ` OFFSET `...)
		q.query = strconv.AppendInt(q.query, int64(offset), 10)
	}
	return q
}

func (q *selectQueryBuilder) Limit(limit int32) *selectQueryBuilder {
	if limit > 0 {
		q.query = append(q.query, ` LIMIT `...)
		q.query = strconv.AppendInt(q.query, int64(limit), 10)
	}
	return q
}

func (q *selectQueryBuilder) Query() string {
	return gotils.B2S(q.query)
}

func (q *selectQueryBuilder) Args() []any {
	return q.args
}
