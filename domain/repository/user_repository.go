package repository

import (
	"context"

	"graphql-project/domain/model"
)

type UserRepository DataSource

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := FindEntity(ctx, (*DataSource)(r), &user, `SELECT {fields} FROM users WHERE id = $1 AND "deletedAt" IS NULL`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := FindEntity(ctx, (*DataSource)(r), &user, `SELECT {fields} FROM users WHERE email = $1 AND "deletedAt" IS NULL`, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, offset int32, limit int32) ([]model.User, error) {
	users := model.NewUsers(max(int(limit), 128))
	err := FindEntities(ctx, (*DataSource)(r), &users, `SELECT {fields} FROM users WHERE "deletedAt" IS NULL ORDER BY id`, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserByIds(ctx context.Context, ids []int64) ([]*model.User, []error) {
	users := model.NewPtrUsers(max(len(ids), 128))
	err := FindEntities(ctx, (*DataSource)(r), &users, `SELECT {fields} FROM users JOIN UNNEST($1::BIGINT[]) WITH ORDINALITY t(id, n) USING(id) WHERE "deletedAt" IS NULL ORDER BY t.n`, 0, 0, ids)
	if err != nil {
		return nil, []error{err}
	}
	if len(users) < len(ids) {
		buffer := make([]*model.User, len(ids))
		n := 0
		for i, id := range ids {
			user := users[n]
			if user.ID == id {
				buffer[i] = user
				n++
			}
		}
		users = buffer
	}
	return users, nil
}

func NewUserRepository(dataSource *DataSource) *UserRepository {
	return (*UserRepository)(dataSource)
}
