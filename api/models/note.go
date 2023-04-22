package models

import "time"

type Note struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateNoteRequest struct {
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateNote struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetAllNotesParams struct {
	Limit      int32  `json:"limit" binding:"required" default:"10"`
	Page       int32  `json:"page" binding:"required" default:"1"`
	Search     string `json:"search"`
	UserID     int64  `json:"user_id"`
	SortByData string `json:"sort_by_data" enums:"asc,desc" default:"desc"`
}

type GetAllNotesResponse struct {
	Notes []*Note `json:"notes"`
	Count int32   `json:"count"`
}
