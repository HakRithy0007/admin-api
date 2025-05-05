package users

import (
	"admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	error_response "admin-phone-shop-api/pkg/utils/error"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse)
}

type UserRepoImpl struct {
	userCtx *custom_models.AdminContext
	db_pool *sqlx.DB
}

// Function to create new repository implementation (was missing)
func NewUserRepoImpl(uCtx *custom_models.AdminContext, db_pool *sqlx.DB) UserRepo {
	return &UserRepoImpl{
		userCtx: uCtx,
		db_pool: db_pool,
	}
}

func (u *UserRepoImpl) CreateUser(createUserReq CreateUserRequest) (*UserResponese, *error_response.ErrorResponse) {
	msg := error_response.ErrorResponse{}

	// Begin transaction
	tx, err := u.db_pool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("transaction_start_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("transaction_start_failed", fmt.Errorf("transaction start failed"))
	}

	// Ensure rollback on panic or failure
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			custom_log.NewCustomLog("create_user_panic", fmt.Sprintf("%v", p), "error")
			panic(p)
		} else if err != nil {
			tx.Rollback() // Only rollback if there was an error
		}
	}()

	var newUser = NewUser{}

	err = newUser.New(createUserReq, u.userCtx, u.db_pool)
	if err != nil {
		custom_log.NewCustomLog("create_member_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("create_member_failed", err)
	}

	query := `
		INSERT INTO tbl_user (
			first_name, last_name, user_name, password, email, login_session, status_id, "order", phone, created_by, created_at, role_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id
	`

	var userID int64
	err = tx.QueryRowx(query,
		newUser.FirstName,
		newUser.LastName,
		newUser.Username,    // Added missing username field
		newUser.Password,
		newUser.Email,
		newUser.LoginSession,
		newUser.StatusID,
		newUser.Order,
		newUser.PhoneNumber, // Fixed field name to match database column 'phone'
		newUser.CreatedBy,
		newUser.CreatedAt,
		newUser.RoleID,      // Added missing role_id parameter
	).Scan(&userID)

	if err != nil {
		custom_log.NewCustomLog("insert_user_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("insert_user_failed", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		custom_log.NewCustomLog("transaction_commit_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("transaction_commit_failed", fmt.Errorf("could not commit transaction"))
	}

	// Return response
	return &UserResponese{
		User: User{
			ID:        int(userID),
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Phone:     newUser.PhoneNumber,
			RoleID:    newUser.RoleID,
			StatusID:  newUser.StatusID,
		},
	}, nil
}