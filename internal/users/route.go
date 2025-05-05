package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserRoute struct {
	app     *fiber.App
	db_pool *sqlx.DB
	handler *UserHandler
}

func NewUserRoute(app *fiber.App, db_pool *sqlx.DB) *UserRoute {
	handler := NewHandler(db_pool)
	return &UserRoute{
		app:     app,
		db_pool: db_pool,
		handler: handler,
	}
}

func (u *UserRoute) RegisterUserRoute() *UserRoute {
	v1 := u.app.Group("/api/v1")
	user := v1.Group("/user")
	// Feature	Description
	// GET /users	List all users (customers)
	// GET /users/:id	View customer details (orders, email, status)
	// PUT /users/:id	Update user info (maybe update name, email, etc.)
	// DELETE /users/:id	Soft delete or deactivate user (e.g., ban)
	// GET /users/:id/orders	View all orders of a specific user
	fmt.Println(user)
	return u
}
