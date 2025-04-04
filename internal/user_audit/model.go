package user_audit

import (
	model "admin-phone-shop-api/pkg/model"
	custom_validator "admin-phone-shop-api/pkg/validator"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

)

type AuditShowRequest struct {
	PageOptions model.PagingOption   `json:"paging_options" query:"paging_options" validate:"required"`
	Sorts       []model.Sort   `json:"sorts,omitempty" query:"sorts"`
	Filters     []model.Filter `json:"filters,omitempty" query:"filters"`
}

func (ua *AuditShowRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {

	if err := c.QueryParser(ua); err != nil {
		return err
	}
	//Fix bug `Filter.Value` nil when http query params failed parse to json type `interface{}`
	for i := range ua.Filters {
		value := c.Query(fmt.Sprintf("filters[%d][value]", i))
		if intValue, err := strconv.Atoi(value); err == nil {
			ua.Filters[i].Value = intValue
		} else if boolValue, err := strconv.ParseBool(value); err == nil {
			ua.Filters[i].Value = boolValue
		} else {
			ua.Filters[i].Value = value
		}
	}

	if err := v.Validate(ua); err != nil {
		return err
	}
	return nil
}

type UserAuditResponse struct {
	UserAudits []UserAudit `json:"user_audits"`
	Total      int         `json:"-"`
}

type UserAudit struct {
	ID               *int    `db:"id" json:"-"`
	UserID           int64   `db:"user_id"`
	UserAuditContext string  `db:"user_audit_context"`
	UserAuditDesc    string  `db:"user_audit_desc"`
	AuditTypeID      int64   `db:"audit_type_id"`
	UserAgent        string  `db:"user_agent"`
	Operator         string  `db:"operator"`
	IP               string  `db:"ip"`
	StatusID         int64   `db:"status_id"`
	Order            int     `db:"order"`
	CreatedBy        int64   `db:"created_by" json:"-"`
	CreatedAt        string  `db:"created_at"`
	UpdatedBy        *int64  `db:"updated_by" json:"-"`
	UpdatedAt        *string `db:"updated_at" json:"-"`
	DeletedBy        *int64  `db:"deleted_by" json:"-"`
	DeletedAt        *string `db:"deleted_at" json:"-"`
}

type TotalRecord struct {
	Total int
}
