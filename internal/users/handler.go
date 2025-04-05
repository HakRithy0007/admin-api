package users

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	model "admin-phone-shop-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	dbPool	 *sqlx.DB
	userService func(*fiber.Ctx) UserCreator
}

func NewUserHandler(dbPool *sqlx.DB, redisClient *redis.Client) *UserHandler {
	return &UserHandler{
		dbPool:     dbPool,
		userService: func(c *fiber.Ctx) UserCreator{
			context := c.Locals("UserContext")
			var uCtx model.UserContext
			if contextMap, ok := context.(model.UserContext); ok {
				uCtx = contextMap
			}else{
				custom_log.NewCustomLog("get_user_context_failed", "Failed to cast UserContext to map[string]interface{}", "warn")
				uCtx = model.UserContext{}
			}

			return NewUserService(&uCtx, dbPool)
		},
	}
}

// Show all Users
func (u *UserHandler) ShowAllUsers(c *fiber.Ctx) error {
	var userRequest ShowUserRequest

	
}
