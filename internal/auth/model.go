package auth

import (
	custom_validator "admin-phone-shop-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthLoginRequest represents the login request payload
type AuthLoginRequest struct {
	Auth struct {
		Username string `json:"username" db:"username" validate:"required"`
		Password string `json:"password" db:"password" validate:"required"`
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

type UserData struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type RedisSession struct {
	LoginSession string `json:"login_session"`
}