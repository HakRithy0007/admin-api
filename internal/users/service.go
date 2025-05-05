package users

import (
	custom_models "admin-phone-shop-api/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserCreator interface {
}

type UserService struct {
	userCtx  *custom_models.AdminContext
	db_pool  *sqlx.DB
	userRepo UserRepo
}

func NewUserService(uCtx *custom_models.AdminContext, db_pool *sqlx.DB) *UserService {
	repo := NewUserRepoImpl(uCtx, db_pool)
	return &UserService{
		userCtx:  uCtx,
		db_pool:  db_pool,
		userRepo: repo,
	}
}
