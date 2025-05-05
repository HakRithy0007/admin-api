package users

import (
	custom_validator "admin-phone-shop-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
)


type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"reqired"`
	LastName string `json:"last_name" validate:"reqired"`
	UserName string `json:"username" validate:"reqired"`
	Password string `json:"password" validate:"reqired, min=6"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (u CreateUserRequest) bind (c *fiber.Ctx, v *custom_validator.Validator) error {
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
