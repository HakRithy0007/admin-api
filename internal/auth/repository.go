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
	audit "admin-phone-shop-api/pkg/utils/audit"
	env "admin-phone-shop-api/pkg/utils/env"
	"admin-phone-shop-api/pkg/utils/error"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	Login(admin_name, password string) (*AuthResponse, *error.ErrorResponse)
	CheckSession(loginSession string, adminID float64) (bool, *error.ErrorResponse)
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
func (a *authRepositoryImpl) Login(admin_name, password string) (*AuthResponse, *error.ErrorResponse) {
	var admin AdminData
	msg := error.ErrorResponse{}

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
		return nil, msg.NewErrorResponse("admin_not_found", fmt.Errorf("admin not found. Please check the provided information"))
	}

	var res AuthResponse

	hours := env.GetenvInt("JWT_EXP_HOUR", 7)
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)
	loginSession, err := uuid.NewV7()

	if err != nil {
		custom_log.NewCustomLog("uuid_generate_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("uuid_generate_failed", fmt.Errorf("failed to generate UUID. Please try again later"))
	}

	claims := jwt.MapClaims{
		"player_id":     admin.ID,
		"admin_name":      admin.Admin_name,
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
		return nil, msg.NewErrorResponse("session_update_failed", fmt.Errorf("cannot update session"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		custom_log.NewCustomLog("jwt_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("jwt_failed", fmt.Errorf("failed to get jwt"))
	}

	res.Auth.Token = tokenString
	res.Auth.TokenType = "jwt"

	auditDesc := fmt.Sprintf(`Admin : %s has been login to the system`, admin_name)
	_, err = audit.AddMemeberAuditLog(admin.ID, "Login", auditDesc, 1, "adminAgent", admin.Admin_name, "ip", admin.ID, a.dbPool)
	if err != nil {
		custom_log.NewCustomLog("add_audit_log_failed", err.Error(), "error")
		return nil, msg.NewErrorResponse("add_audit_log_failed", fmt.Errorf("cannot insert data to audit log"))
	}

	return &res, nil
}

// CheckSession
func (a *authRepositoryImpl) CheckSession(loginSession string, adminID float64) (bool, *error.ErrorResponse) {
	msg := error.ErrorResponse{}

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