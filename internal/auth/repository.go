package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"admin-phone-shop-api/pkg/custom_log"
	redis_util "admin-phone-shop-api/pkg/redis"
	env "admin-phone-shop-api/pkg/utils/env"
	error_response "admin-phone-shop-api/pkg/utils/error"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	Login(admin_name, password string) (*AdminData, string, *error_response.ErrorResponse)
	CheckSession(loginSession string, adminID int) (bool, *error_response.ErrorResponse)
	Logout(adminID int) (bool, *error_response.ErrorResponse)
}

type authRepositoryImpl struct {
	dbPool *sqlx.DB
	redis  *redis.Client
}

func NewAuthRepository(dbPool *sqlx.DB, redisClient *redis.Client) AuthRepository {
	return &authRepositoryImpl{
		dbPool: dbPool,
		redis:  redisClient,
	}
}

// Login
func (a *authRepositoryImpl) Login(admin_name, password string) (*AdminData, string, *error_response.ErrorResponse) {
	var admin AdminData
	msg := error_response.ErrorResponse{}

	query := `
		SELECT
			id, 
			admin_name,
			email,
			password
		FROM tbl_admin 
		WHERE admin_name = $1 AND password = $2
	`

	err := a.dbPool.Get(&admin, query, admin_name, password)
	if err != nil {
		custom_log.NewCustomLog("admin_not_found", err.Error(), "error")
		return nil, "", msg.NewErrorResponse("admin_not_found", fmt.Errorf("admin not found. Please check the provided information"))
	}

	hours := env.GetenvInt("JWT_EXP_HOUR", 7)
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)
	loginSession, err := uuid.NewV7()

	if err != nil {
		custom_log.NewCustomLog("uuid_generate_failed", err.Error(), "error")
		return nil, "", msg.NewErrorResponse("uuid_generate_failed", fmt.Errorf("failed to generate UUID. Please try again later"))
	}

	claims := jwt.MapClaims{
		"admin_id":     admin.ID,
		"admin_name":    admin.Admin_name,
		"login_session": loginSession.String(),
		"exp":           expirationTime.Unix(),
	}

	// Set Redis Data
	key := fmt.Sprintf("admin: %d", admin.ID)
	redisUtil := redis_util.NewRedisUtil(a.redis)
	redisUtil.SetCacheKey(key, claims, context.Background())

	_ = godotenv.Load()
	secretKey := os.Getenv("JWT_SECRET_KEY")

	updateQuery := `	
		UPDATE tbl_admin
		SET login_session = $1
		WHERE id = $2
	`
	_, err = a.dbPool.Exec(updateQuery, loginSession.String(), admin.ID)
	if err != nil {
		custom_log.NewCustomLog("session_update_failed", err.Error(), "error")
		return nil, "", msg.NewErrorResponse("session_update_failed", fmt.Errorf("cannot update session"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		custom_log.NewCustomLog("jwt_failed", err.Error(), "error")
		return nil, "", msg.NewErrorResponse("jwt_failed", fmt.Errorf("failed to get jwt"))
	}

	return &admin, tokenString, nil
}

// CheckSession
func (a *authRepositoryImpl) CheckSession(loginSession string, adminID int) (bool, *error_response.ErrorResponse) {
	msg := error_response.ErrorResponse{}

	key := fmt.Sprintf("admin: %d", int(adminID))
	redisUtil := redis_util.NewRedisUtil(a.redis)

	keyData, err := redisUtil.GetCacheKey(key, context.Background())
	if err == nil {
		if keyData.LoginSession == loginSession {
			return true, nil
		}
	}

	var storedLoginSession string

	query := `
		SELECT login_session
		FROM tbl_admin
		WHERE login_session = $1
	`
	err = a.dbPool.Get(&storedLoginSession, query, loginSession)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			custom_log.NewCustomLog("invalid_session_id", "invalid login session: "+loginSession, "warn")
			return false, msg.NewErrorResponse("invalid_session_id", fmt.Errorf("invalid login session"))
		}
		custom_log.NewCustomLog("query_data_failed", err.Error(), "error")
		return false, msg.NewErrorResponse("query_data_failed", fmt.Errorf("database query error"))
	}

	if storedLoginSession != loginSession {
		return false, msg.NewErrorResponse("invalid_session_id", fmt.Errorf("invalid login session"))
	}
	return true, nil
}

// Logout
func (a *authRepositoryImpl) Logout(adminID int) (bool, *error_response.ErrorResponse) {
	msg := error_response.ErrorResponse{}

	// Delete session from redis
	key := fmt.Sprintf("admin: %d", int(adminID))
	redis_util := redis_util.NewRedisUtil(a.redis)
	if err := redis_util.DeleteCache(key); err != nil {
		custom_log.NewCustomLog("redis_delete_failed", err.Error(), "error")
	}

	// Remove session from database
	updateQuery :=
		`
		UPDATE tbl_admin
		SET login_session = NULL
		WHERE id = $1
		`
	result, err := a.dbPool.Exec(updateQuery, adminID)
	if err != nil {
		custom_log.NewCustomLog("logout_failed", err.Error(), "error")
		return false, msg.NewErrorResponse("logout_failed", fmt.Errorf("could not log out"))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		custom_log.NewCustomLog("logout_failed", "no rows affected", "warn")
		return false, msg.NewErrorResponse("logout_failed", fmt.Errorf("logout failed"))
	}

	return true, nil
}