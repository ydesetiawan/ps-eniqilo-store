package repository

import (
	"ps-eniqilo-store/internal/user/model"
	"ps-eniqilo-store/pkg/errs"
	"strings"

	"github.com/jmoiron/sqlx"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetUserByPhoneNumber(phoneNumber string) (model.User, error) {
	var user model.User
	query := "select * from users where phone_number = $1 "
	err := r.db.Get(&user, query, phoneNumber)
	if err != nil {
		return model.User{}, errs.NewErrInternalServerErrors("execute query error [GetUserByPhoneNumber]: ", err.Error())
	}
	return user, err
}

func (r *userRepositoryImpl) GetUserByPhoneNumberAndId(phoneNumber string, id int64) (model.User, error) {
	var user model.User
	query := "select * from users where phone_number = $1 and id = $2 "
	err := r.db.Get(&user, query, phoneNumber, id)
	if err != nil {
		return model.User{}, errs.NewErrInternalServerErrors("execute query error [GetUserByPhoneNumber]: ", err.Error())
	}
	return user, err
}

func (r *userRepositoryImpl) RegisterUser(user *model.User) (int64, error) {
	var lastInsertId int64 = 0
	query := "insert into users (phone_number, name, password) values($1, $2, $3) RETURNING id"

	err := r.db.QueryRowx(query, user.PhoneNumber, user.Name, user.Password).Scan(&lastInsertId)
	if err != nil {
		if strings.Contains(err.Error(), "users_phone_number_key") {
			return 0, errs.NewErrDataConflict("phoneNumber already exist", user.PhoneNumber)
		}
		return 0, errs.NewErrInternalServerErrors("execute query error [GetUserByPhoneNumber]: ", err.Error())
	}

	return lastInsertId, nil
}
