package repository

import (
	"gopkg.in/reform.v1"

	"graphql-pro/domain/model"
	"graphql-pro/interface/repository"
)

type userRepository struct {
	DataSource
}

func (repository userRepository) GetUserByID(id int64) (*model.User, error) {
	if data, err := repository.FindByPrimaryKeyFrom(model.UserTable, id); err != nil {
		if err == reform.ErrNoRows {
			return nil, nil
		}
		return nil, err
	} else {
		user := data.(*model.User)
		if user.DeletedAt.Valid {
			return nil, nil
		}
		return user, nil
	}
}

func (repository userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repository.SelectOneTo(&user, `WHERE email = $1`, email); err != nil {
		if err == reform.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if user.DeletedAt.Valid {
		return nil, nil
	}
	return &user, nil
}

func NewUserRepository(ds DataSource) repository.UserRepository {
	return userRepository{ds}
}
