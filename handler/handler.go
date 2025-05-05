package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	auth "admin-phone-shop-api/internal/auth"
	user "admin-phone-shop-api/internal/users"
	middleware "admin-phone-shop-api/pkg/middleware"
)

type ServiceHandlers struct {
	Fronted *FrontService
}

type FrontService struct {
	AuthHandler *auth.AuthRoute
	UserHandler *user.UserRoute
}

func NewFrontService(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *FrontService {

	// Authentication
	auth := auth.NewAuthRoute(app, db_pool, redis).RegisterAuthRoute()

	// Middleware
	middleware.NewJwtMinddleWare(app, db_pool, redis)

	// User
	user := user.NewUserRoute(app, db_pool).RegisterUserRoute()

	return &FrontService{
		AuthHandler: auth,
		UserHandler: user,
	}
}

func NewServiceHandlers(app *fiber.App, db_pool *sqlx.DB, redis *redis.Client) *ServiceHandlers {

	return &ServiceHandlers{
		Fronted: NewFrontService(app, db_pool, redis),
	}
}
