package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AdminRoute struct {
	app     *fiber.App
	db_pool *sqlx.DB
	handler *AdminHandler
}

func NewAdminRoute(app *fiber.App, db_pool *sqlx.DB) *AdminRoute {
	handler := NewHandler(db_pool)
	return &AdminRoute{
		app:     app,
		db_pool: db_pool,
		handler: handler,
	}
}

func (u *AdminRoute) RegisterAdminRoute() *AdminRoute {
	v1 := u.app.Group("/api/v1")
	admin := v1.Group("/admin")

	admin.Get("/", u.handler.ShowAll)
	admin.Get("/:id", u.handler.ShowOne)

	// POST	/admin/logout	(Optional) Logout admin
	// POST	/admins	Create new admin account
	// GET	/admins	List all admin users
	// GET	/admins/:id	Get details of an admin
	// PUT	/admins/:id	Update admin info or password
	// DELETE	/admins/:id	Soft delete or deactivate admin

	// Print the registered routes in a more detailed way
	return u
}
