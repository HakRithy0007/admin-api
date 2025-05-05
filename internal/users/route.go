package users

import (
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
	user.Post("/create-user", u.handler.CreateUser)
	// user.Get("/show-all", u.handler.ShowAll)
	// user.Get("/show-one", u.handler.ShowOne)
	// user.Put("/update-user", u.handler.UpdateUser)
	// user.Delete("/delete-user", u.handler.DeleteUser)
	// user.Put("/change-password", u.handler.ChangePassword)
	// user.Get("/get-form-create", u.handler.GetFormCreate)
	// user.Get("/get-form-update", u.handler.GetFormUpdate)
	// user.Get("/get-basic-info", u.handler.GetBasicInfo)
	return u
}
