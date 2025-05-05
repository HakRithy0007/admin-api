package users

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	sql "admin-phone-shop-api/pkg/sql"
	custom_validator "admin-phone-shop-api/pkg/validator"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateUserRequest struct {
	FirstName   string `json:"first_name" validate:"reqired"`
	LastName    string `json:"last_name" validate:"reqired"`
	UserName    string `json:"username" validate:"reqired"`
	Password    string `json:"password" validate:"reqired, min=6"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleID      int    `json:"role_id"`
}

func (u CreateUserRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.BodyParser(u); err != nil {
		return err
	}

	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID           int     `db:"id" json:"id"`
	FirstName    string  `db:"first_name" json:"first_name"`
	LastName     string  `db:"last_name" json:"last_name"`
	Username     string  `db:"user_name" json:"username"`
	Email        string  `db:"email" json:"email"`
	LoginSession *string `db:"login_session" json:"-"`
	Phone        string  `db:"phone" json:"phone"`
	Password     string  `db:"password" json:"-"`
	StatusID     int     `db:"status_id" json:"status_id"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
	CreatedBy    int     `db:"created_by" json:"created_by"`
	DeletedAt    *string `db:"deleted_at" json:"-"`
	RoleID       int     `db:"role_id" json:"role_id"`
	UserRoleName string  `db:"user_role_name" json:"user_role_name"`
	Operator     string  `json:"operator" db:"operator"`
}

type UserResponese struct {
	User User `json:"user"`
}

type NewUser struct {
	ID           int       `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Username     string    `db:"user_name"`
	Email        string    `db:"email"`
	LoginSession *string   `db:"login_session"`
	Phone        string    `db:"phone"`
	Password     string    `db:"password"`
	StatusID     int       `db:"status_id"`
	OrderBy      int       `db:"order"`
	CreatedAt    time.Time `db:"created_at"`
	CreatedBy    int       `db:"created_by"`
	RoleID       int       `db:"role_id"`
}

func (u *NewUser) New(createUserReq CreateUserRequest, aCtx *custom_models.AdminContext, db_pool *sqlx.DB) error {

	if aCtx.RoleID > createUserReq.RoleID {
		return fmt.Errorf("failed you role can not create this user")
	}

	login_session, err := uuid.NewV7()
	if err != nil {
		custom_log.NewCustomLog("get_uuid_failed", err.Error(), "error")
		return err
	}
	sessionString := login_session.String()

	app_timezone := os.Getenv("TIME_ZONE")
	location, err := time.LoadLocation(app_timezone)
	if err != nil {
		return fmt.Errorf("failed to load location: %w", err)
	}
	local_now := time.Now().In(location)

	is_username, err := sql.IsExits("tbl_user", "user_name", createUserReq.UserName, db_pool)
	if err != nil {
		return err
	} else {
		if is_username {
			return fmt.Errorf("%s", fmt.Sprintf("username:`%s` already exists", createUserReq.UserName))
		}
	}
	createdByID, err := sql.GetAdminIdByField("tbl_users", "user_name", aCtx.Admin_Name, db_pool)
	if err != nil {
		return err
	}

	orderValue, err := sql.GetSeqNextVal("tbl_users_id_seq", db_pool)
	if err != nil {
		return fmt.Errorf("failed to generate order value: %w", err)
	}

	u.FirstName = createUserReq.FirstName
	u.LastName = createUserReq.LastName
	u.Password = createUserReq.Password
	u.Email = createUserReq.Email
	u.LoginSession = &sessionString
	u.StatusID = 1
	u.OrderBy = *orderValue
	u.Phone = createUserReq.PhoneNumber
	u.CreatedBy = *createdByID
	u.CreatedAt = local_now
	u.RoleID = createUserReq.RoleID

	return nil
}
