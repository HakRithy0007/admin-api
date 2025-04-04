package user_audit

import (
	"admin-phone-shop-api/pkg/constants"
	"admin-phone-shop-api/pkg/utils/response"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	model "admin-phone-shop-api/pkg/model"
	"admin-phone-shop-api/pkg/utils/translate"
	custom_validator "admin-phone-shop-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type UserAuditHandler struct {
	db_pool          *sqlx.DB
	userAuditService func(*fiber.Ctx) UserAuditCreator
}

func NewHandler(db_pool *sqlx.DB) *UserAuditHandler {
	return &UserAuditHandler{
		db_pool: db_pool,
		userAuditService: func(c *fiber.Ctx) UserAuditCreator {
			userContext := c.Locals("userContext")
			var uCtx model.UserContext

			if contextMap, ok := userContext.(model.UserContext); ok {
				uCtx = contextMap
			} else {
				custom_log.NewCustomLog(
					"user_context_cast_failed",
					"Failed to cast UserContext to share.UserContext",
					"warn",
				)
				uCtx = model.UserContext{}
			}

			// Pass uCtx to NewAuthService if needed
			return NewUserAuditService(&uCtx, db_pool)
		},
	}
}

// Show audit funtion
func (ua *UserAuditHandler) Show(c *fiber.Ctx) error {
	auditReq := AuditShowRequest{}
	v := custom_validator.NewValidator()

	// Bind request data and validate
	if err := auditReq.bind(c, v); err != nil {
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

	// UserAudit Service
	service := ua.userAuditService(c)
	responses, err := service.Show(auditReq)
	if err != nil {
		msg, err_msg := translate.TranslateWithError(c, "show_user_audit_failed")
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
				constants.Show_user_audit_Failed,
				err.Err,
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponseWithPaing(
			"success",
			constants.Show_user_success,
			responses,
			auditReq.PageOptions.Page,
			auditReq.PageOptions.PerPage,
			responses.Total,
		),
	)
}
