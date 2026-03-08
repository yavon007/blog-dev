package core

import (
	"context"
	"fmt"

	"github.com/yavon007/blog-dev/backend/internal/pkg/pagination"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
)

type Repository interface {
	List(ctx context.Context, filter ListCommentsFilter, p pagination.Params) ([]*Comment, int64, error)
	FindByID(ctx context.Context, id int64) (*Comment, error)
	Create(ctx context.Context, postID int64, req CreateCommentRequest, ipHash, ua string) (*Comment, error)
	UpdateStatus(ctx context.Context, id int64, status CommentStatus) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPublic(ctx context.Context, postID int64, p pagination.Params) ([]*Comment, int64, error) {
	return s.repo.List(ctx, ListCommentsFilter{PostID: postID, Status: StatusApproved}, p)
}

func (s *Service) ListAdmin(ctx context.Context, filter ListCommentsFilter, p pagination.Params) ([]*Comment, int64, error) {
	return s.repo.List(ctx, filter, p)
}

func (s *Service) Create(ctx context.Context, postID int64, req CreateCommentRequest, ipHash, ua string) (*Comment, error) {
	comment, err := s.repo.Create(ctx, postID, req, ipHash, ua)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, req UpdateCommentStatusRequest) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return sherrors.ErrNotFound
	}
	return s.repo.UpdateStatus(ctx, id, req.Status)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return sherrors.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}
