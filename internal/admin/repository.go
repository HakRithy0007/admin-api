package admin

import (
	"admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	error_response "admin-phone-shop-api/pkg/utils/error"
	custom_sql "admin-phone-shop-api/pkg/sql"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AdminRepo interface {
	ShowAll(adminRequest AdminShowRequest) (*AdminShowResponse, *error_response.ErrorResponse)
	ShowOne(id int) (*AdminResponse, *error_response.ErrorResponse)
	CreateNewAdmin(crreq CreateAdminRequest) (*CreateAdminResponse, *error_response.ErrorResponse)
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

// Show All
func (u *AdminRepoImpl) ShowAll(adminRequest AdminShowRequest) (*AdminShowResponse, *error_response.ErrorResponse) {

	err_msg := &error_response.ErrorResponse{}

	var per_page = adminRequest.PageOption.PerPage
	var page = adminRequest.PageOption.Page
	var offset = (page - 1) * per_page
	var limit_clause = fmt.Sprintf("LIMIT %d OFFSET %d", per_page, offset)
	var sql_orderby = custom_sql.BuildSQLSort(adminRequest.Sort)

	sql_filters, args_filters := custom_sql.BuildSQLFilter(adminRequest.Filters)
	if len(args_filters) > 0 {
		sql_filters = "AND " + sql_filters
	}

	query := fmt.Sprintf(
		`
		SELECT
			a.id,
			a.first_name,
			a.last_name,
			a.admin_name AS user_name,
			a.email,
			a.phone,
			a.status_id,
			a.created_at,
			a.deleted_at,
			a.created_by,
			a.role_id,
			r.admin_role_name,
			creator.admin_name AS operator
		FROM
			tbl_admin a
		INNER JOIN
			tbl_admin_roles r ON a.role_id = r.id
		INNER JOIN
			tbl_admin creator ON a.created_by = creator.id
		WHERE
			a.deleted_at IS NULL
			AND r.deleted_at is NULL %s %s %s
		`, sql_filters, sql_orderby, limit_clause)

		var NewAdmin []Admin

		err := u.db_pool.Select(&NewAdmin, query, args_filters...)
		if err != nil {
			custom_log.NewCustomLog("could_not_query", err.Error(), "error")
			return nil, err_msg.NewErrorResponse("could_not_query", fmt.Errorf("can not select admin data the database is error"))
		}

		totalQuery := fmt.Sprintf(`
			SELECT
				COUNT(*) as total
				FROM 
					tbl_admin a
				INNER JOIN
					tbl_admin_roles r ON a.role_id = r.id
				INNER JOIN
					tbl_admin creator ON a.created_by = creator.id
				WHERE a.deleted_at IS NULL %s
		`, sql_filters)

		var total TotalRecord

		err = u.db_pool.Get(&total, totalQuery, args_filters...)
		if err != nil {
			custom_log.NewCustomLog("admin_show_failed", err.Error(), "error")
			return nil, err_msg.NewErrorResponse("admin_show_failed", fmt.Errorf("can not select admin total the database error"))
		}

		return &AdminShowResponse{
			Admin: NewAdmin,
			Total: total.Total,
		}, nil
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

// Create new admin
func (u *AdminRepoImpl) CreateNewAdmin(crreq CreateAdminRequest) (*CreateAdminResponse, *error_response.ErrorResponse) {

	err_msg := &error_response.ErrorResponse{}

	// Transaction
	tx, err := u.db_pool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("create_admin_failed", err.Error(), "error")
		return nil, err_msg.NewErrorResponse("create_admin_failed", fmt.Errorf("Create admin failed"))
	}
	var newAdmin = NewAdmin{}
}