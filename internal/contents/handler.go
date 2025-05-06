package users

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	db_pool     *sqlx.DB
	userService func(*fiber.Ctx) UserCreator
}

func NewHandler(db_pool *sqlx.DB) *UserHandler {
	return &UserHandler{
		db_pool: db_pool,
		userService: func(c *fiber.Ctx) UserCreator {
			context := c.Locals("AdminContext")

			var uCtx custom_models.AdminContext
			if contextMap, ok := context.(custom_models.AdminContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog("get_user_context_failed", "Failed to cast AdminContext to map[string]interface{}", "warn")
				uCtx = custom_models.AdminContext{}
			}

			return NewUserService(&uCtx, db_pool)
		},
	}
}
