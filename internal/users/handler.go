package users

import (
	constants "admin-phone-shop-api/pkg/constants"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	response "admin-phone-shop-api/pkg/utils/response"
	translate "admin-phone-shop-api/pkg/utils/translate"
	custom_validator "admin-phone-shop-api/pkg/validator"

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


// Create User
func (u *UserHandler) CreateUser(c *fiber.Ctx) error {
	var createUserReq = CreateUserRequest{}
	v := custom_validator.NewValidator()

	if err := createUserReq.bind(c, v); err != nil {
		msg, err_msg := translate.TranslateWithError(c, "invalid_request")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					err_msg.ErrorString(),
					constants.Translate_failed,
					err_msg.Err,
				),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(
				msg,
				constants.Invalid_request,
				err,
			),
		)
	}

	// Service
	service := u.userService(c)
	user, err := service.CreateUser(createUserReq)
	if err != nil {
		msg, err_msg := translate.TranslateWithError(c, err.MessageID)
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					err_msg.ErrorString(),
					constants.Translate_failed,
					err_msg.Err,
				),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(
				msg,
				constants.Created_user_failed,
				err.Err,
			),
		)
	} else {

		msg, err_msg := translate.TranslateWithError(c, "user_create_success")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewResponse(
				err_msg.ErrorString(),
				constants.Translate_failed,
				err_msg.Err,
			))
		}
		return c.Status(fiber.StatusOK).JSON(
			response.NewResponse(
				msg,
				constants.Created_user_success,
				user,
			))
	}

}