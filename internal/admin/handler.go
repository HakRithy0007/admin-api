package admin

import (
	constants "admin-phone-shop-api/pkg/constants"
	custom_log "admin-phone-shop-api/pkg/custom_log"
	custom_models "admin-phone-shop-api/pkg/model"
	response "admin-phone-shop-api/pkg/utils/response"
	translate "admin-phone-shop-api/pkg/utils/translate"
	custom_validator "admin-phone-shop-api/pkg/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AdminHandler struct {
	db_pool      *sqlx.DB
	adminService func(*fiber.Ctx) AdminCreator
}

func NewHandler(db_pool *sqlx.DB) *AdminHandler {
	return &AdminHandler{
		db_pool: db_pool,
		adminService: func(c *fiber.Ctx) AdminCreator {
			context := c.Locals("AdminContext")

			var aCtx custom_models.AdminContext
			if contextMap, ok := context.(custom_models.AdminContext); ok {
				aCtx = contextMap
			} else {
				custom_log.NewCustomLog("get_admin_context_failed", "Failed to cast AdminContext to map[string]interface{}", "warn")
				aCtx = custom_models.AdminContext{}
			}

			return NewAdminService(&aCtx, db_pool)
		},
	}
}

// Show all
func (u *AdminHandler) ShowAll(c *fiber.Ctx) error {
	var AdminReq AdminShowRequest

	v := custom_validator.NewValidator()
	if err := AdminReq.bind(c, v); err != nil {
		msg, err_msg := translate.TranslateWithError(c, "invalid_request")
		if err_msg != nil{
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
	service := u.adminService(c)
	success, err := service.ShowAll(AdminReq)
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
				constants.Show_admin_Failed,
				err.Err,
			),
		)
	} else {
		msg, err_msg := translate.TranslateWithError(c, "show_admn_success")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewResponse(
				err_msg.ErrorString(),
				constants.Translate_failed,
				err_msg.Err,
			))
		}
		return c.Status(fiber.StatusOK).JSON(
			response.NewResponseWithPaing(
				msg,
				constants.Show_admin_success,
				success,
				AdminReq.PageOption.Page,
				AdminReq.PageOption.PerPage,
				success.Total,
			))
	}
}

// ShowOne
func (u *AdminHandler) ShowOne(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
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
	service := u.adminService(c)
	ressuccess, err_service := service.ShowOne(id)
	if err_service != nil {
		msg, err_msg := translate.TranslateWithError(c, err_service.MessageID)
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
				constants.Get_current_admin_failed,
				err_service.Err,
			),
		)
	} else {
		msg, err_msg := translate.TranslateWithError(c, "show_admin_success")
		if err_msg != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.NewResponseError(
				err_msg.ErrorString(),
				constants.Translate_failed,
				err_msg.Err,
			))
		}
		return c.Status(fiber.StatusOK).JSON(
			response.NewResponse(
				msg,
				constants.Get_current_admin_success,
				ressuccess,
			))
	}
}
