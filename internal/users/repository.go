package users

import (
	custom_models "admin-phone-shop-api/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
}

type UserRepoImpl struct {
	userCtx *custom_models.AdminContext
	db_pool *sqlx.DB
}

func NewUserRepoImpl(uCtx *custom_models.AdminContext, db_pool *sqlx.DB) *UserRepoImpl {
	return &UserRepoImpl{
		userCtx: uCtx,
		db_pool: db_pool,
	}
}
