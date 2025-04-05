package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserRoute struct {
	app     *fiber.App
	dbPool  *sqlx.DB
	handler *UserHandler
}

func NewUserRoute(app *fiber.App, dbPool *sqlx.DB) *UserRoute {
	handler := NewUserHandler(dbPool)
	return &UserRoute{
		app: app,
		dbPool: dbPool,
		handler: handler,
	}
}

func (u *UserRoute) RegisterUserRoute() *UserRoute {
	v1 := u.app.Group("api/v1")
	users := v1.Group("/admin/users")

	// Show all users
	users.Get("/show", u.handler.ShowAllUsers)

	return u
}

// GET /admin/users → Get all users

// GET /admin/users/{id} → Get details of a specific user

// PUT /admin/users/{id} → Update user details

// DELETE /admin/users/{id} → Delete a user

// POST /admin/users/{id}/role → Assign a role to a user (e.g., customer, admin)
