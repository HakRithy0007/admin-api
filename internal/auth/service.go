package auth

import (
	"admin-phone-shop-api/pkg/utils/error"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// AuthService defines the service layer for authentication
type AuthService interface {
	Login(admin_name, password string) (*AuthResponse, *error.ErrorResponse)
	CheckSession(loginSession string, adminID float64) (bool, *error.ErrorResponse)
	Logout(adminID float64, loginSession string) (bool, *error.ErrorResponse)
}

// authServiceImpl implements AuthService
type authServiceImpl struct {
	repo AuthRepository
}

func NewAuthService(dbPool *sqlx.DB, redisClient *redis.Client) AuthService {
	repo := NewAuthRepository(dbPool, redisClient)
	return &authServiceImpl{
		repo: repo,
	}
}

// Login
func (a *authServiceImpl) Login(admin_name, password string) (*AuthResponse, *error.ErrorResponse) {
	return a.repo.Login(admin_name, password)
}

// Check session
func (a *authServiceImpl) CheckSession(loginSession string, adminID float64) (bool, *error.ErrorResponse) {
	return a.repo.CheckSession(loginSession, adminID)
}

// Logout
func (a *authServiceImpl) Logout(adminID float64, loginSession string) (bool, *error.ErrorResponse) {
	return a.repo.Logout(adminID, loginSession)
}