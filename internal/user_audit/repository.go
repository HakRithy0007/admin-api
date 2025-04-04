package user_audit

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	"admin-phone-shop-api/pkg/utils/error"
	model "admin-phone-shop-api/pkg/model"
	sql "admin-phone-shop-api/pkg/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserAuditRepo interface {
	Show(audit_req AuditShowRequest) (*UserAuditResponse, *error.ErrorResponse)
}

type UserAuditRepoImpl struct {
	userCtx *model.UserContext
	db      *sqlx.DB
}

func NewUserAuditRepoImpl(u *model.UserContext, db *sqlx.DB) *UserAuditRepoImpl {
	return &UserAuditRepoImpl{
		userCtx: u,
		db:      db,
	}
}

func (ua *UserAuditRepoImpl) Show(audit_req AuditShowRequest) (*UserAuditResponse, *error.ErrorResponse) {
	msg := error.ErrorResponse{}
	// Prepare paging options
	perPage := audit_req.PageOptions.PerPage
	page := audit_req.PageOptions.Page
	offset := (page - 1) * perPage
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", perPage, offset)

	// Prepare order by clause
	orderByClause := sql.BuildSQLSort(audit_req.Sorts)

	// Prepare filter clause
	filterClause, filterArgs := sql.BuildSQLFilter(audit_req.Filters)
	if len(filterArgs) > 0 {
		filterClause = " AND " + filterClause
	}

	// SQL query for user audits
	query := fmt.Sprintf(`
		SELECT
			user_id, 
			user_audit_context, 
			user_audit_desc, 
			audit_type_id,
			user_agent, 
			operator, 
			ip, status_id, 
			"order", 
			created_at
		FROM
			tbl_users_audit
		WHERE
			deleted_at IS NULL %s %s %s
	`, filterClause, orderByClause, limitClause)

	// SQL query for total count
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(*) AS total
		FROM
			tbl_users_audit
		WHERE
			deleted_at IS NULL %s
	`, filterClause)

	// Execute the user audits query
	var userAudits []UserAudit
	err := ua.db.Select(&userAudits, query, filterArgs...)
	if err != nil {
		custom_log.NewCustomLog("execute_users_audits_show_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("execute_users_audits_show_failed", fmt.Errorf("execute_users_audits_show_failed."))
	}

	// Execute the total count query
	var total []TotalRecord
	err = ua.db.Select(&total, countQuery, filterArgs...)
	if err != nil {
		return nil, msg.NewErrorResponse("users_audits_show_failed", fmt.Errorf("users_audits_show_failed."))
	}

	// Return the response
	return &UserAuditResponse{
		UserAudits: userAudits,
		Total:      total[0].Total,
	}, nil
}
