package users

import (
	"admin-phone-shop-api/pkg/constants"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	model "admin-phone-shop-api/pkg/model"
	response "admin-phone-shop-api/pkg/utils/response"
	"admin-phone-shop-api/pkg/utils/translate"
	custom_validator "admin-phone-shop-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	dbPool      *sqlx.DB
	userService func(*fiber.Ctx) UserCreator
}

func NewUserHandler(dbPool *sqlx.DB) *UserHandler {
	return &UserHandler{
		dbPool: dbPool,
		userService: func(c *fiber.Ctx) UserCreator {
			context := c.Locals("UserContext")
			var uCtx *model.UserContext
			
			// Fix the type assertion
			if contextPtr, ok := context.(*model.UserContext); ok {
				// Direct pointer cast was successful
				uCtx = contextPtr
			} else if contextVal, ok := context.(model.UserContext); ok {
				// Value cast was successful, convert to pointer
				uCtx = &contextVal
			} else {
				// Log the failure and create an empty context
				custom_log.NewCustomLog("get_user_context_failed", "Failed to cast UserContext from Locals", "warn")
				uCtx = &model.UserContext{}
			}

			return NewUserService(uCtx, dbPool)
		},
	}
}

// Show all Users
func (u *UserHandler) ShowAllUsers(c *fiber.Ctx) error {
	var userRequest ShowUserRequest

	// Validate request
	v := custom_validator.NewValidator()
	if err := userRequest.bind(c, v); err != nil {
		msg, errMsg := translate.TranslateWithError(c, "invalid_request")
		if errMsg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					errMsg.ErrorString(),
					constants.Translate_failed,
					errMsg.Err,
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

	// Call service
	service := u.userService(c)
	responses, err := service.ShowAllUser(userRequest)
	if err != nil {
		msg, errMsg := translate.TranslateWithError(c, err.MessageID)
		if errMsg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewResponseError(
					errMsg.ErrorString(),
					constants.Translate_failed,
					errMsg.Err,
				),
			)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(
				msg,
				constants.Show_user_Failed,
				err,
			),
		)
	}

	// Success response
	msg, errMsg := translate.TranslateWithError(c, "show_user_success")
	if errMsg != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewResponseError(
				errMsg.ErrorString(),
				constants.Translate_failed,
				errMsg.Err,
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(
			msg,
			constants.Show_user_success,
			responses,
		),
	)
}
