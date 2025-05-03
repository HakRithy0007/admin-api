package users
 
 import (
 	model "admin-phone-shop-api/pkg/model"
 	"fmt"
 	"strconv"
 	"time"
 
 	custom_validator "admin-phone-shop-api/pkg/validator"
 
 	"github.com/gofiber/fiber/v2"
 )
 
 type User struct {
 	ID           int        `db:"id"`
 	FirstName    string     `db:"first_name"`
 	LastName     string     `db:"last_name"`
 	User_Name     string     `db:"user_name"`
 	Email        string     `db:"email"`
 	Phone        string     `db:"phone"`
 	StatusID     int        `db:"status_id"`
 	CreatedAt    time.Time  `db:"created_at"`
 	DeletedAt    *time.Time `db:"deleted_at"` // pointer because nullable
 	CreatedBy    int        `db:"created_by"`
 	RoleID       int        `db:"role_id"`
 	UserRoleName string     `db:"user_role_name"`
 	Operator     string     `db:"operator"`
 }
 
 type ShowUserRequest struct {
 	PagingOption model.PagingOption `json:"paging_options" query:"paging_options"`
 	Sorts        []model.Sort       `json:"sorts,omitempty" query:"sorts"`
 	Filters      []model.Filter     `json:"filters,omitempty" query:"filters"`
 }
 
 type ShowUserResponse struct {
 	Users []User `json:"users"`
 	Total int    `json:"total"`
 }
 
 type TotalRecord struct {
 	Total int `db:"total"`
 }
 
 func (u *ShowUserRequest) bind(c *fiber.Ctx, v *custom_validator.Validator) error {
 	if err := c.QueryParser(u); err != nil {
 		return err
 	}
 
 	for i := range u.Filters {
 		value := c.Query(fmt.Sprintf("filters[%d][value]", i))
 		if intValue, err := strconv.Atoi(value); err == nil {
 			u.Filters[i].Value = intValue
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