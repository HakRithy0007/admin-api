package users

import (
	"admin-phone-shop-api/pkg/custom_log"
	model "admin-phone-shop-api/pkg/model"
	custom_sql "admin-phone-shop-api/pkg/sql"
	"admin-phone-shop-api/pkg/utils/error"
	"fmt"

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

// Show all user
func (u *UserRepoImpl) ShowAllUser(userReqest ShowUserRequest) (*ShowUserResponse, *error.ErrorResponse) {

	err_msg := &error.ErrorResponse{}

	var perPage = userReqest.PagingOption.PerPage
	var page = userReqest.PagingOption.Page
	var offset = (page - 1) * perPage
	var limitClause = fmt.Sprintf("LIMIT %d OFFSET %d", perPage, offset)
	var sqlOrderBy = custom_sql.BuildSQLSort(userReqest.Sorts)

	sqlFilters, argsFilters := custom_sql.BuildSQLFilter(userReqest.Filters)
	if len(argsFilters) > 0 {
		sqlFilters = " AND " + sqlFilters
	}
	

	query := fmt.Sprintf(`
		SELECT 
				u.id, 
				u.first_name,
				u.last_name,
				u.user_name,
				u.email,
				u.phone,
				u.status_id,
				u.created_at,
				u.deleted_at,
				u.created_by,
				u.role_id,
				r.user_role_name,
				creator.user_name AS operator
			FROM
				tbl_users u
			INNER JOIN 
				tbl_users_roles r ON u.role_id = r.id 
			INNER JOIN
				tbl_users creator ON u.created_by = creator.id
			WHERE 
				u.deleted_at IS NULL 
				AND r.deleted_at IS NULL   %s %s %s
	`, sqlFilters, sqlOrderBy, limitClause)

	var Users []User

	err := u.db_pool.Select(&Users, query, argsFilters...)
	if err != nil {
		custom_log.NewCustomLog("could_not_query", err.Error(), "error")
		return nil, err_msg.NewErrorResponse("could_not_query", fmt.Errorf("can select user the database error"))
	}

	totalQuery := fmt.Sprintf(`
		SELECT
					COUNT(*) as total
				FROM 
					tbl_users u
				INNER JOIN 
					tbl_users_roles r ON u.role_id = r.id 
				INNER JOIN
					tbl_users creator ON u.created_by = creator.id
				WHERE u.deleted_at IS NULL %s
	`, sqlFilters)

	var total TotalRecord

	err = u.db_pool.Get(&total, totalQuery, argsFilters...)
	if err != nil {
		custom_log.NewCustomLog("could_not_query", err.Error(), "error")
		return nil, err_msg.NewErrorResponse("could_not_query", fmt.Errorf("can select user the database error"))
	}

	return &ShowUserResponse{
		Users: Users,
		Total: total.Total,
	}, nil
}
