package custom_models

import "time"

type AdminContext struct {
	AdminID     float64   `json:"admin_id"`
	Admin_Name     string    `json:"admin_name"`
	LoginSession string    `json:"login_session"`
	Exp          time.Time `json:"exp"`
	AdminAgent    string    `json:"admin_agent"`
	Ip           string    `json:"ip"`
	MembershipId float64   `json:"membership_id"`
	RoleID       int       `json:"role_id"`
}
type Token struct {
	Id       float64 `json:"id"`
	Admin_name string  `json:"admin_name"`
}
type PagingOption struct {
	PerPage int `json:"perpage" query:"per_page" validate:"required"`
	Page    int `json:"page" query:"page" validate:"required"`
}
type Filter struct {
	Property string      `json:"property" query:"property"`
	Value    interface{} `json:"value" query:"value"`
}
type Sort struct {
	Property  string `json:"property" query:"property"`
	Direction string `json:"direction" query:"direction"`
}
