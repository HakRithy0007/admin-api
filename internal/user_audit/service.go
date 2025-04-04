package user_audit

import (
	"admin-phone-shop-api/pkg/utils/error"
	model "admin-phone-shop-api/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserAuditCreator interface {
	Show(audit_req AuditShowRequest) (*UserAuditResponse, *error.ErrorResponse)
}

type UserAuditService struct {
	userCtx       *model.UserContext
	db_pool       *sqlx.DB
	userAuditRepo UserAuditRepo
}

func NewUserAuditService(uCtx *model.UserContext, db_pool *sqlx.DB) *UserAuditService {
	repo := NewUserAuditRepoImpl(uCtx, db_pool)
	return &UserAuditService{
		userCtx:       uCtx,
		db_pool:       db_pool,
		userAuditRepo: repo,
	}
}

func (ua *UserAuditService) Show(audit_req AuditShowRequest) (*UserAuditResponse, *error.ErrorResponse) {
	success, err := ua.userAuditRepo.Show(audit_req)
	return success, err
}
