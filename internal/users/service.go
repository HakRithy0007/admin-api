package users

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	audit "admin-phone-shop-api/pkg/utils/audit"
	error_response "admin-phone-shop-api/pkg/utils/error"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserCreator interface {
	CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse)
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

func (u *UserService) CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse) {
	success, err := u.userRepo.CreateUser(createUserReq)
	if err != nil {
		return nil, err
	}

	// Add audit
	auditDesc := fmt.Sprintf(`Admin ID: %d has logged out`, createUserReq.UserName)
	_, auditErr := audit.AddMemeberAuditLog(success.User.ID, "Logout", auditDesc, 1, "adminAgent", "", "ip", u.userCtx.AdminID, u.db_pool)
	if auditErr != nil {
		custom_log.NewCustomLog("add_audit_log_failed", auditErr.Error(), "error")
		errResponse := &error_response.ErrorResponse{}
		return nil, errResponse.NewErrorResponse("add_audit_log_failed", fmt.Errorf("cannot insert logout audit log"))
	}

	return success, nil

}
