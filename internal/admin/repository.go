package admin

import (
	"admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	error_response "admin-phone-shop-api/pkg/utils/error"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AdminRepo interface {
	ShowOne(id int) (*AdminResponse, *error_response.ErrorResponse)
}

type AdminRepoImpl struct {
	adminCtx *custom_models.AdminContext
	db_pool  *sqlx.DB
}

func NewAdminRepoImpl(aCtx *custom_models.AdminContext, db_pool *sqlx.DB) AdminRepo {
	return &AdminRepoImpl{
		adminCtx: aCtx,
		db_pool:  db_pool,
	}
}

// Show One
func (u *AdminRepoImpl) ShowOne(id int) (*AdminResponse, *error_response.ErrorResponse) {
	msg := error_response.ErrorResponse{}

	if id <= 0 {
		return nil, msg.NewErrorResponse("invalid_user_id", fmt.Errorf("invalid user id '%d' received", id))
	}

	var admin AdminOne
	query :=
		`
		SELECT
			id,
			first_name,
			last_name,
			admin_name,
			email,
			phone,
			status_id,
			created_at,
			created_by,
			deleted_at,
			deleted_by,
			role_id
		FROM tbl_admin
		WHERE id = $1
	`
	err := u.db_pool.Get(&admin, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			custom_log.NewCustomLog("admin_not_found", fmt.Sprintf("No admin found with id '%d'", id), "error")
			return nil, msg.NewErrorResponse("admin_not_found", fmt.Errorf("no admin found with id '%d'", id))
		}
		custom_log.NewCustomLog("server_error", err.Error(), "error")
		return nil, msg.NewErrorResponse("server_error", fmt.Errorf("data error: %v", err))
	}
	return &AdminResponse{
		AdminInfo: admin,
	}, nil
}
