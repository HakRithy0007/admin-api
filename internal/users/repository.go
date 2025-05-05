package users

import (
	"admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	error_response "admin-phone-shop-api/pkg/utils/error"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse)
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

func (u UserRepoImpl) CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse) {

	// Transaction 
	tx, err := u.db_pool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("transaction_start_failed", err.Error(), "error")
		return nil, &error_response.ErrorResponse{
			MessageID: ,
		}
	}
}