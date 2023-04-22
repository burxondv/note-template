package models

import "time"

type User struct {
	ID          int        `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	Email       string     `json:"email"`
	ImageURL    string     `json:"image_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	ImageURL    string `json:"image_url"`
}

type GetAllUserParams struct {
	Limit  int32  `json:"limit" binding:"required" default:"10"`
	Page   int32  `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
}

type GetAllUsersResponse struct {
	Users []*User `json:"users"`
	Count int     `json:"count"`
}
