package core

import "time"

type Category struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	PostCount   int64
	CreatedAt   time.Time
}

type Tag struct {
	ID        int64
	Name      string
	Slug      string
	PostCount int64
	CreatedAt time.Time
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Slug        string `json:"slug" binding:"required,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest = CreateCategoryRequest

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=100"`
	Slug string `json:"slug" binding:"required,max=100"`
}

type UpdateTagRequest = CreateTagRequest
