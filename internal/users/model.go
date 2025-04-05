package users

import (
	model "admin-phone-shop-api/pkg/model"
)

type ShowUserRequest struct {
	PagingOption model.PagingOption `json:"paging_options" query:"paging_options"`
	Sorts        []model.Sort       `json:"sorts,omitempty" query:"sorts"`
	Filters      []model.Filter     `json:"filters,omitempty" query:"filters"`
}
