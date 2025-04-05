package users

import (
	model "admin-phone-shop-api/pkg/model"
	"admin-phone-shop-api/pkg/utils/error"

	"github.com/jmoiron/sqlx"
)

type UserCreator interface {
	ShowAllUser(userRequest ShowUserRequest) (*ShowUserResponse, *error.ErrorResponse)
}

type UserService struct {
	userCtx  *model.UserContext
	db_pool  *sqlx.DB
	userRepo UserRepo
}

func NewUserService(uCtx *model.UserContext, db_pool *sqlx.DB) *UserService {
	repo := NewUserRepoImpl(uCtx, db_pool)
	return &UserService{
		userCtx:  uCtx,
		db_pool:  db_pool,
		userRepo: repo,
	}
}

// Show All
func (u *UserService) ShowAllUser(userRequest ShowUserRequest) (*ShowUserResponse, *error.ErrorResponse) {
	user, err := u.userRepo.ShowAllUser(userRequest)
	if err != nil {
		return nil, err
	}
	return user, nil
}
