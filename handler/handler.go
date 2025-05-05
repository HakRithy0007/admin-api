package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	auth "admin-phone-shop-api/internal/auth"
	user "admin-phone-shop-api/internal/users"
	middleware "admin-phone-shop-api/pkg/middleware"
	admin "admin-phone-shop-api/internal/admin"

)

type ServiceHandlers struct {
	Fronted *FrontService
}

type FrontService struct {
	AuthHandler *auth.AuthRoute
	UserHandler *user.UserRoute
	AdminHandler *admin.AdminRoute
}

func NewFrontService(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *FrontService {

	// Authentication
	auth := auth.NewAuthRoute(app, db_pool, redis).RegisterAuthRoute()

	// Middleware
	middleware.NewJwtMinddleWare(app, db_pool, redis)

	// User
	user := user.NewUserRoute(app, db_pool).RegisterUserRoute()

	// Admin
	admin := admin.NewAdminRoute(app, db_pool).RegisterAdminRoute()

	return &FrontService{
		AuthHandler: auth,
		UserHandler: user,
		AdminHandler: admin,
	}
}

func NewServiceHandlers(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *ServiceHandlers {

	return &ServiceHandlers{
		Fronted: NewFrontService(app, db_pool, redis),
	}
}
