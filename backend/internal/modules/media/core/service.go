package core

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, file *MediaFile) error
	FindByID(ctx context.Context, id int64) (*MediaFile, error)
	List(ctx context.Context, filter ListMediaFilter) ([]*MediaFile, int64, error)
	SoftDelete(ctx context.Context, id int64) error
}

type StorageBackend interface {
	Save(ctx context.Context, filename string, reader io.Reader) (path string, url string, err error)
	Delete(ctx context.Context, path string) error
}

type Service struct {
	repo    Repository
	storage StorageBackend
	baseURL string
}

func NewService(repo Repository, storage StorageBackend, baseURL string) *Service {
	return &Service{repo: repo, storage: storage, baseURL: baseURL}
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func (s *Service) Upload(ctx context.Context, file *multipart.FileHeader, altText string, uploadedBy int64) (*UploadResult, error) {
	// Validate MIME type
	if !allowedMimeTypes[file.Header.Get("Content-Type")] {
		return nil, fmt.Errorf("unsupported file type: %s", file.Header.Get("Content-Type"))
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return nil, fmt.Errorf("file too large: max 5MB")
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Save to storage
	path, url, err := s.storage.Save(ctx, filename, f)
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	// Create media record
	media := &MediaFile{
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		Size:         file.Size,
		AltText:      altText,
		Storage:      "local",
		Path:         path,
		URL:          url,
		UploadedBy:   &uploadedBy,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, media); err != nil {
		_ = s.storage.Delete(ctx, path)
		return nil, fmt.Errorf("create media record: %w", err)
	}

	return &UploadResult{
		ID:       media.ID,
		URL:      media.URL,
		Filename: media.OriginalName,
	}, nil
}

func (s *Service) List(ctx context.Context, filter ListMediaFilter) ([]*MediaFile, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete in database
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// Delete from storage
	_ = s.storage.Delete(ctx, file.Path)
	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*MediaFile, error) {
	return s.repo.FindByID(ctx, id)
}
