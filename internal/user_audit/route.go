package user_audit

import (

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserAuditRoute struct {
	app     *fiber.App
	db_pool *sqlx.DB
	handler *UserAuditHandler
}

func NewRoute(app *fiber.App, db_pool *sqlx.DB) *UserAuditRoute {
	handler := NewHandler(db_pool)
	return &UserAuditRoute{
		app:     app,
		db_pool: db_pool,
		handler: handler,
	}
}

func (u *UserAuditRoute) RegisterUserRoute() *UserAuditRoute {
	v1 := u.app.Group("/api/v1")
	user_audit := v1.Group("/useraudit")
	user_audit.Get("/", u.handler.Show)
	return u
}
