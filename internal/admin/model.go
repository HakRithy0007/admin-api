package admin

import (
	"admin-phone-shop-api/pkg/custom_log"
	model "admin-phone-shop-api/pkg/model"
	sql "admin-phone-shop-api/pkg/sql"
	custom_validator "admin-phone-shop-api/pkg/validator"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AdminResponse struct {
	AdminInfo AdminOne `json:"admin-info"`
}

type AdminOne struct {
	ID          int        `json:"-" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	AdminName   string     `json:"admin_name" db:"admin_name"`
	Email       string     `json:"email" db:"email"`
	PhoneNumber string     `json:"phone" db:"phone"`
	StatusID    int        `json:"status_id" db:"status_id"`
	RoleID      int        `json:"role_id" db:"role_id"`
	CreatedBy   int        `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
	DeletedBy   *int       `json:"-" db:"deleted_by"`
}

type AdminShowRequest struct {
	PageOption model.PagingOption `json:"paging_options" query:"paging_options" validate:"required"`
	Sort       []model.Sort       `json:"sorts,omitempty" query:"sorts"`
	Filters    []model.Filter     `json:"filters,omitempty" query:"filters"`
}

func (u *AdminShowRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.QueryParser(u); err != nil {
		return err
	}

	for i := range u.Filters {
		value := c.Query(fmt.Sprintf("filters[%d][value]", i))
		if intValue, err := strconv.Atoi(value); err == nil {
			u.Filters[i].Value = intValue
		} else if boolValue, err := strconv.ParseBool(value); err == nil {
			u.Filters[i].Value = boolValue
		} else {
			u.Filters[i].Value = value
		}
	}
	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
}

type Admin struct {
	ID            int     `db:"id" json:"id"`
	FirstName     string  `db:"first_name" json:"first_name"`
	LastName      string  `db:"last_name" json:"last_name"`
	AdminName      string  `db:"admin_name" json:"admin_name"`
	Email         string  `db:"email" json:"email"`
	LoginSession  *string `db:"login_session" json:"-"`
	Phone         string  `db:"phone" json:"phone"`
	Password      string  `db:"password" json:"-"`
	StatusID      int     `db:"status_id" json:"status_id"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
	CreatedBy     int     `db:"created_by" json:"created_by"`
	DeletedAt     *string `db:"deleted_at" json:"-"`
	RoleID        int     `db:"role_id" json:"role_id"`
	AdminRoleName string  `db:"admin_role_name" json:"admin_role_name"`
	Operator      string  `json:"operator" db:"operator"`
}

type AdminShowResponse struct {
	Admin []Admin `json:"admins"`
	Total int     `json:"-"`
}

type TotalRecord struct {
	Total int `db:"total"`
}

type CreateAdminRequest struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	AdminName       string `json:"admin_name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone"`
	RoleID          int    `json:"role_id"`
}

func (u *CreateAdminRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.BodyParser(u); err != nil {
		return err
	}
	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
}

type CreateAdminResponse struct {
	Admin Admin `json:"admins"`
}

type NewAdmin struct {
	ID           int       `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	AdminName    string    `db:"admin_name"`
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

func (u *NewAdmin) New(createAdminReq *CreateAdminRequest, uCtx *model.AdminContext, db_pool *sqlx.DB) error {
	// Check for missing role_id in the request
	if createAdminReq.RoleID == 0 {
		// Default to a standard user role (adjust the value based on your system)
		createAdminReq.RoleID = 2 // Assuming 2 is a standard user role
	}

	if uCtx.RoleID > createAdminReq.RoleID {
		return fmt.Errorf("failed: your role cannot create this admin")
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

	is_adminName, err := sql.IsExits("tbl_admin", "admin_name", createAdminReq.AdminName, db_pool)
	if err != nil {
		return err
	} else {
		if is_adminName {
			return fmt.Errorf("%s", fmt.Sprintf("admin name:`%s` already exists", createAdminReq.AdminName))
		}
	}

	createdByID, err := sql.GetAdminIdByField("tbl_admin", "admin_name", uCtx.Admin_Name, db_pool)
	if err != nil {
		return err
	}

	orderValue, err := sql.GetSeqNextVal("tbl_admin_id_seq", db_pool)
	if err != nil {
		return fmt.Errorf("failed to generate order value: %w", err)
	}

	u.FirstName = createAdminReq.FirstName
	u.LastName = createAdminReq.LastName
	u.AdminName = createAdminReq.AdminName
	u.Password = createAdminReq.Password
	u.Email = createAdminReq.Email
	u.LoginSession = &sessionString
	u.StatusID = 1
	u.OrderBy = *orderValue
	u.Phone = createAdminReq.Phone
	u.CreatedBy = *createdByID
	u.CreatedAt = local_now
	u.RoleID = createAdminReq.RoleID

	return nil
}