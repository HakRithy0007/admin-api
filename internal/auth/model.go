package auth

import (
	custom_validator "admin-phone-shop-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthLoginRequest represents the login request payload
type AuthLoginRequest struct {
	Auth struct {
		Admin_name string `json:"admin_name" db:"admin_name" validate:"required"`
		Password   string `json:"password" db:"password" validate:"required"`
	} `json:"auth"`
}

// Bind request payload
func (r *AuthLoginRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

type AuthResponse struct {
	Auth struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type"`
	} `json:"auths"`
}

type AdminData struct {
	ID         int    `db:"id" json:"id"`
	Admin_name string `db:"admin_name" json:"admin_name"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password" json:"password"`
}

type RedisSession struct {
	LoginSession string `json:"login_session"`
}

type LogoutRequest struct {
	AdminID int `json:"admin_id" validate:"required"`
}

func (l *LogoutRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.BodyParser(l); err != nil {
		return err
	}
	return v.Validate(l)
}
