package repository

import (
	"context"

	"graphql-project/core"
	"graphql-project/interface/model"
)

type selectByIdQuery struct {
	selectQueryBuilder
	ctx context.Context
	id  int64
}

func (q *selectByIdQuery) Build(entity model.Entity) {
	q.selectQueryBuilder = *Select(getFields(q.ctx, entity)).From(entity.Table()).Where(entity.Identity(), q.id).And(`"deletedAt" IS NULL`, nil)
}

func SelectById(ctx context.Context, id int64) *selectByIdQuery {
	return &selectByIdQuery{ctx: ctx, id: id}
}

type selectManyQuery struct {
	selectQueryBuilder
	ctx    context.Context
	offset int32
	limit  int32
}

func (q *selectManyQuery) Build(entity model.Entity) {
	q.selectQueryBuilder = *Select(getFields(q.ctx, entity)).From(entity.Table()).Where(`"deletedAt" IS NULL`, nil).
		OrderBy(entity.Identity(), false).Offset(q.offset).Limit(q.limit)
}

func SelectMany(ctx context.Context, offset int32, limit int32) *selectManyQuery {
	return &selectManyQuery{ctx: ctx, offset: offset, limit: limit}
}

type selectByIdsQuery struct {
	selectQueryBuilder
	ctx context.Context
	ids []int64
}

func (q *selectByIdsQuery) Build(entity model.Entity) {
	query := make([]byte, 0, 144)
	query = appendQuoted("SELECT ", query, getFields(q.ctx, entity))
	query = core.AppendStrings(query, ` FROM "`, entity.Table(), `" JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t("`, entity.Identity(), `", n) USING("`, entity.Identity(), `") WHERE "deletedAt" IS NULL ORDER BY t.n`)
	q.selectQueryBuilder = selectQueryBuilder{query, []any{q.ids}}
}

func SelectByIds(ctx context.Context, ids []int64) *selectByIdsQuery {
	return &selectByIdsQuery{ctx: ctx, ids: ids}
}

type selectByQuery struct {
	selectQueryBuilder
	ctx   context.Context
	name  string
	value any
}

func (q *selectByQuery) Build(entity model.Entity) {
	q.selectQueryBuilder = *Select(getFields(q.ctx, entity)).From(entity.Table()).Where(q.name, q.value).And(`"deletedAt" IS NULL`, nil)
}

func SelectBy(ctx context.Context, name string, value any) *selectByQuery {
	return &selectByQuery{ctx: ctx, name: name, value: value}
}

type selectByRefIdQuery struct {
	selectQueryBuilder
	ctx    context.Context
	ref    string
	refId  int64
	offset int32
	limit  int32
}

func (q *selectByRefIdQuery) Build(entity model.Entity) {
	q.selectQueryBuilder = *Select(getFields(q.ctx, entity)).From(entity.Table()).Where(q.ref, q.refId).And(`"deletedAt" IS NULL`, nil)
}

func SelectByRefId(ctx context.Context, ref string, refId int64, offset int32, limit int32) *selectByRefIdQuery {
	return &selectByRefIdQuery{ctx: ctx, ref: ref, refId: refId, offset: offset, limit: limit}
}

type selectByRefIdsQuery struct {
	selectQueryBuilder
	ctx    context.Context
	ref    string
	refIds []int64
	offset int32
	limit  int32
}

func (q *selectByRefIdsQuery) Build(entity model.Entity) {
	query := make([]byte, 0, 320)
	args := make([]any, 0, 4)
	var from int64
	if q.offset > 0 {
		from = int64(q.offset)
	} else {
		from = 1
	}
	args = append(args, q.refIds)
	args = append(args, from)
	query = core.AppendStrings(query, `WITH o AS (SELECT *, ROW_NUMBER() OVER (PARTITION BY "`, q.ref, `" ORDER BY "`, entity.Identity(), `") AS r FROM "`, entity.Table(), `" WHERE "deletedAt" IS NULL AND "`, q.ref)
	query = appendQuoted(`" = ANY($1::BIGINT[])) SELECT `, query, getFields(q.ctx, entity))
	query = core.AppendStrings(query, `,"`, q.ref, `" FROM o JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t("`, q.ref, `", n) USING("`, q.ref, `") WHERE r >= $2`)
	if q.limit > 0 {
		to := from + int64(q.limit)
		args = append(args, to)
		query = append(query, ` AND r < $3`...)
	}
	query = append(query, ` ORDER BY t.n, r`...)
	q.selectQueryBuilder = selectQueryBuilder{query, args}
}

func SelectByRefIds(ctx context.Context, ref string, refIds []int64, offset int32, limit int32) *selectByRefIdsQuery {
	return &selectByRefIdsQuery{ctx: ctx, ref: ref, refIds: refIds, offset: offset, limit: limit}
}
