package core

import (
	"context"
	"fmt"

	sherrors "github.com/yourblog/backend/internal/modules/shared/errors"
)

type Repository interface {
	ListCategories(ctx context.Context) ([]*Category, error)
	FindCategoryByID(ctx context.Context, id int64) (*Category, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error)
	UpdateCategory(ctx context.Context, id int64, req UpdateCategoryRequest) (*Category, error)
	DeleteCategory(ctx context.Context, id int64) error

	ListTags(ctx context.Context) ([]*Tag, error)
	FindTagByID(ctx context.Context, id int64) (*Tag, error)
	CreateTag(ctx context.Context, req CreateTagRequest) (*Tag, error)
	UpdateTag(ctx context.Context, id int64, req UpdateTagRequest) (*Tag, error)
	DeleteTag(ctx context.Context, id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListCategories(ctx context.Context) ([]*Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	c, err := s.repo.CreateCategory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id int64, req UpdateCategoryRequest) (*Category, error) {
	if _, err := s.repo.FindCategoryByID(ctx, id); err != nil {
		return nil, sherrors.ErrNotFound
	}
	return s.repo.UpdateCategory(ctx, id, req)
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	if _, err := s.repo.FindCategoryByID(ctx, id); err != nil {
		return sherrors.ErrNotFound
	}
	return s.repo.DeleteCategory(ctx, id)
}

func (s *Service) ListTags(ctx context.Context) ([]*Tag, error) {
	return s.repo.ListTags(ctx)
}

func (s *Service) CreateTag(ctx context.Context, req CreateTagRequest) (*Tag, error) {
	return s.repo.CreateTag(ctx, req)
}

func (s *Service) UpdateTag(ctx context.Context, id int64, req UpdateTagRequest) (*Tag, error) {
	if _, err := s.repo.FindTagByID(ctx, id); err != nil {
		return nil, sherrors.ErrNotFound
	}
	return s.repo.UpdateTag(ctx, id, req)
}

func (s *Service) DeleteTag(ctx context.Context, id int64) error {
	if _, err := s.repo.FindTagByID(ctx, id); err != nil {
		return sherrors.ErrNotFound
	}
	return s.repo.DeleteTag(ctx, id)
}
