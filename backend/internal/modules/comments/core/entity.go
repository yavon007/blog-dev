package core

import "time"

type CommentStatus string

const (
	StatusPending  CommentStatus = "pending"
	StatusApproved CommentStatus = "approved"
	StatusRejected CommentStatus = "rejected"
)

type Comment struct {
	ID              int64
	PostID          int64
	ParentCommentID *int64
	AuthorName      string
	AuthorEmail     string
	Body            string
	Status          CommentStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateCommentRequest struct {
	AuthorName      string `json:"author_name" binding:"required,max=100"`
	AuthorEmail     string `json:"author_email" binding:"required,email"`
	Body            string `json:"body" binding:"required,min=2"`
	ParentCommentID *int64 `json:"parent_comment_id"`
}

type UpdateCommentStatusRequest struct {
	Status CommentStatus `json:"status" binding:"required,oneof=approved rejected"`
}

type ListCommentsFilter struct {
	PostID int64
	Status CommentStatus
}
