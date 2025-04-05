package users

import (
	model "admin-phone-shop-api/pkg/model"
	"fmt"
	"strconv"

	custom_validator "admin-phone-shop-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type ShowUserRequest struct {
	PagingOption model.PagingOption `json:"paging_options" query:"paging_options"`
	Sorts        []model.Sort       `json:"sorts,omitempty" query:"sorts"`
	Filters      []model.Filter     `json:"filters,omitempty" query:"filters"`
}

func (u *ShowUserRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
	if err := c.QueryParser(u), err != nil {
		return err
	}

	for i := range u.Filters {
		value := c.Query(fmt.Sprintf("filters[%d].value", i))
		if intValue, err := strconv.Atoi(value), err == nil {
			u.Filters[i].Value =intValue
		} else if boolValue, err := strconv.ParseBool(value); err == nil {
			u.Filters[i].Value = boolValue
		} else {
			u.Filters[i].Value = value
		}
	}
	if err := v.Validate(u); err != nil {
		return err
	}
	return nil
}
