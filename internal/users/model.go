package users

import (
	model "admin-phone-shop-api/pkg/model"
	"fmt"
	"strconv"

	custom_validator "admin-phone-shop-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

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


type ShowUserRequest struct {
	PagingOption model.PagingOption `json:"paging_options" query:"paging_options"`
	Sorts        []model.Sort       `json:"sorts,omitempty" query:"sorts"`
	Filters      []model.Filter     `json:"filters,omitempty" query:"filters"`
}

type ShowUserResponse struct {
	Users []User `json:"users"`
	Total int         `json:"total"`
}

func (u *ShowUserRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
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