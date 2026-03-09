package core

import "time"

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

type Post struct {
	ID                int64      `json:"id"`
	Title             string     `json:"title"`
	Slug              string     `json:"slug"`
	Summary           string     `json:"summary"`
	ContentMD         string     `json:"content_md"`
	ContentHTMLCached string     `json:"content_html_cached"`
	CoverURL          string     `json:"cover_url"`
	Status            PostStatus `json:"status"`
	PublishedAt       *time.Time `json:"published_at"`
	CategoryID        *int64     `json:"category_id"`
	CategoryName      string     `json:"category_name"`
	AuthorID          int64      `json:"author_id"`
	Tags              []Tag      `json:"tags"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
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

// ArchiveItem represents a year-month archive group
type ArchiveItem struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Count int `json:"count"`
}
