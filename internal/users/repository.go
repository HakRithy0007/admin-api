package users


import (
	model "admin-phone-shop-api/pkg/model"
	"admin-phone-shop-api/pkg/utils/error"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	ShowAllUser(userReqest ShowUserRequest) (*ShowUserResponse, *error.ErrorResponse)
}

type UserRepoImpl struct {
	userCtx *model.UserContext
	db_pool *sqlx.DB
}

func NewUserRepoImpl(uCtx *model.UserContext, db_pool *sqlx.DB) *UserRepoImpl {
	return &UserRepoImpl{
		userCtx: uCtx,
		db_pool: db_pool,
	}
}

func (u *UserRepoImpl) ShowAllUser(userReqest ShowUserRequest) (*ShowUserResponse, *error.ErrorResponse) {

	return &ShowUserResponse{}, nil
}