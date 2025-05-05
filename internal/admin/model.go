package admin

import "time"

type AdminResponse struct {
	Admin AdminOne `json:"admin"`
}

type AdminOne struct {
	ID          int        `json:"-" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	AdminName    string     `json:"admin_name" db:"admin_name"`
	Email       string     `json:"email" db:"email"`
	PhoneNumber string     `json:"phone" db:"phone"`
	StatusID    int        `json:"status_id" db:"status_id"`
	RoleID      int        `json:"role_id" db:"role_id"`
	CreatedBy   int        `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
	DeletedBy   *int       `json:"-" db:"deleted_by"`
}
