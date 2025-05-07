package admin

import (
	custom_models "admin-phone-shop-api/pkg/model"
	error_response "admin-phone-shop-api/pkg/utils/error"

	"github.com/jmoiron/sqlx"
)

type AdminCreator interface {
	ShowAll(adminReqeust AdminShowRequest) (*AdminShowResponse, *error_response.ErrorResponse)
	ShowOne(id int) (*AdminResponse, *error_response.ErrorResponse) 
	CreateNewAdmin(crreq CreateAdminRequest) (*CreateAdminResponse, *error_response.ErrorResponse)
}

type AdminService struct {
	adminCtx  *custom_models.AdminContext
	db_pool  *sqlx.DB
	adminRepo AdminRepo
}

func NewAdminService(aCtx *custom_models.AdminContext, db_pool *sqlx.DB) *AdminService {
	repo := NewAdminRepoImpl(aCtx, db_pool)
	return &AdminService{
		adminCtx:  aCtx,
		db_pool:  db_pool,
		adminRepo: repo,
	}
}

// Show All
func (u *AdminService) ShowAll(adminReqeust AdminShowRequest) (*AdminShowResponse, *error_response.ErrorResponse) {
	return u.adminRepo.ShowAll(adminReqeust)
}

// Show One
func (u *AdminService) ShowOne(id int) (*AdminResponse, *error_response.ErrorResponse) {
	return  u.adminRepo.ShowOne(id)
}

// Create new admin
func (u *AdminService) CreateNewAdmin(crreq CreateAdminRequest) (*CreateAdminResponse, *error_response.ErrorResponse) {
	return u.adminRepo.CreateNewAdmin(crreq)
}