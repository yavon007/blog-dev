package core

import "time"

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

type Post struct {
	ID                int64
	Title             string
	Slug              string
	Summary           string
	ContentMD         string
	ContentHTMLCached string
	CoverURL          string
	Status            PostStatus
	PublishedAt       *time.Time
	CategoryID        *int64
	CategoryName      string
	AuthorID          int64
	Tags              []Tag
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Tag struct {
	ID   int64
	Name string
	Slug string
}

type CreatePostRequest struct {
	Title      string     `json:"title" binding:"required,max=255"`
	Slug       string     `json:"slug" binding:"required,max=255"`
	Summary    string     `json:"summary"`
	ContentMD  string     `json:"content_md" binding:"required"`
	CoverURL   string     `json:"cover_url"`
	Status     PostStatus `json:"status" binding:"omitempty,oneof=draft published"`
	CategoryID *int64     `json:"category_id"`
	TagIDs     []int64    `json:"tag_ids"`
}

type UpdatePostRequest struct {
	Title      string     `json:"title" binding:"required,max=255"`
	Slug       string     `json:"slug" binding:"required,max=255"`
	Summary    string     `json:"summary"`
	ContentMD  string     `json:"content_md" binding:"required"`
	CoverURL   string     `json:"cover_url"`
	Status     PostStatus `json:"status" binding:"omitempty,oneof=draft published"`
	CategoryID *int64     `json:"category_id"`
	TagIDs     []int64    `json:"tag_ids"`
}

type UpdateStatusRequest struct {
	Status PostStatus `json:"status" binding:"required,oneof=draft published"`
}

type ListFilter struct {
	Category string
	Tag      string
	Query    string
	Status   PostStatus // empty = all (admin), "published" (public)
}
