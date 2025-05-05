package auth

import (
	"fmt"
	audit "admin-phone-shop-api/pkg/utils/audit"
	error_response "admin-phone-shop-api/pkg/utils/error"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// AuthService defines the service layer for authentication
type AuthService interface {
	Login(admin_name, password string) (*AuthResponse, *error_response.ErrorResponse)
	CheckSession(loginSession string, adminID float64) (bool, *error_response.ErrorResponse)
	Logout(adminID float64) (bool, *error_response.ErrorResponse)
}


// authServiceImpl implements AuthService
type authServiceImpl struct {
	repo AuthRepository
	dbPool *sqlx.DB
}

func NewAuthService(dbPool *sqlx.DB, redisClient *redis.Client) AuthService {
	repo := NewAuthRepository(dbPool, redisClient)
	return &authServiceImpl{
		repo: repo,
		dbPool: dbPool,
	}
}

// Login
func (a *authServiceImpl) Login(admin_name, password string) (*AuthResponse, *error_response.ErrorResponse) {
	admin, tokenString, err := a.repo.Login(admin_name, password)
	if err != nil {
		return nil, err
	}
	
	// Create response
	var res AuthResponse
	res.Auth.Token = tokenString
	res.Auth.TokenType = "jwt"

	// Add audit log
	auditDesc := fmt.Sprintf(`Admin : %s has been login to the system`, admin_name)
	_, auditErr := audit.AddMemeberAuditLog(float64(admin.ID), "Login", auditDesc, 1, "adminAgent", admin.AdminName, "ip", float64(admin.ID), a.dbPool)
	if auditErr != nil {
		custom_log.NewCustomLog("add_audit_log_failed", auditErr.Error(), "error")
		errResponse := &error_response.ErrorResponse{}
		return nil, errResponse.NewErrorResponse("add_audit_log_failed", fmt.Errorf("cannot insert data to audit log"))
	}
	
	return &res, nil
}

// Check session
func (a *authServiceImpl) CheckSession(loginSession string, adminID float64) (bool, *error_response.ErrorResponse) {
	return a.repo.CheckSession(loginSession, adminID)
}

// Logout
func (a *authServiceImpl) Logout(adminID float64) (bool, *error_response.ErrorResponse) {
	success, err := a.repo.Logout(adminID)
	if err != nil {
		return false, err
	}

	// Add audit 
	auditDesc := fmt.Sprintf(`Admin ID: %f has logged out`, adminID)
	_, auditErr := audit.AddMemeberAuditLog(adminID, "Logout", auditDesc, 1, "adminAgent", "", "ip", adminID, a.dbPool)
	if auditErr != nil {
		custom_log.NewCustomLog("add_audit_log_failed", auditErr.Error(), "error")
		errResponse := &error_response.ErrorResponse{}
		return true, errResponse.NewErrorResponse("add_audit_log_failed", fmt.Errorf("cannot insert logout audit log"))
	}

	return success, nil
}