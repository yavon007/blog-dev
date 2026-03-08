package core

import (
	"context"
	"fmt"

	"github.com/yavon007/blog-dev/backend/internal/pkg/pagination"
	sherrors "github.com/yavon007/blog-dev/backend/internal/modules/shared/errors"
)

type Repository interface {
	List(ctx context.Context, filter ListFilter, p pagination.Params) ([]*Post, int64, error)
	FindBySlug(ctx context.Context, slug string) (*Post, error)
	FindByID(ctx context.Context, id int64) (*Post, error)
	Create(ctx context.Context, req CreatePostRequest, authorID int64) (*Post, error)
	Update(ctx context.Context, id int64, req UpdatePostRequest) (*Post, error)
	UpdateStatus(ctx context.Context, id int64, status PostStatus) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPublic(ctx context.Context, filter ListFilter, p pagination.Params) ([]*Post, int64, error) {
	filter.Status = StatusPublished
	return s.repo.List(ctx, filter, p)
}

func (s *Service) ListAdmin(ctx context.Context, filter ListFilter, p pagination.Params) ([]*Post, int64, error) {
	return s.repo.List(ctx, filter, p)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Post, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreatePostRequest, authorID int64) (*Post, error) {
	if req.Status == "" {
		req.Status = StatusDraft
	}
	post, err := s.repo.Create(ctx, req, authorID)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return post, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdatePostRequest) (*Post, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return nil, sherrors.ErrNotFound
	}
	post, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	return post, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, req UpdateStatusRequest) error {
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
