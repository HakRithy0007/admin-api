package admin

import (
	model "admin-phone-shop-api/pkg/model"
	custom_validator "admin-phone-shop-api/pkg/validator"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
	AdminRoleName string  `db:"admin_role_name" json:"admin_role_name"`
	Operator     string  `json:"operator" db:"operator"`
}

type AdminShowResponse struct {
	Admin []Admin `json:"admins"`
	Total int `json:"-"`
}

type TotalRecord struct {
	Total int `db:"total"`
}

type CreateUserRequest struct{
	FirstName string `json:"first_name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	AdminName string `json:"admin_name" validate:"required"`
	Password string `json:"password" validate:"required, min=6"`
	PasswordConfirm string `json:"password_confirm" validate:"required, min=6"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone"`
	StatusID int `json:"status_id"`
	RoleID int `json:"role_id"`
}

func(u *CreateUserRequest) bind (c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.BodyParser(u); err != nil {
		return err
	}
	if err := v.Validate(u); err != nil{
		return err
	}
	return nil
}