package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
	baseURL   string
}

func NewLocalStorage(uploadDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}
	return &LocalStorage{uploadDir: uploadDir, baseURL: baseURL}, nil
}

func (s *LocalStorage) Save(ctx context.Context, filename string, reader io.Reader) (string, string, error) {
	path := filepath.Join(s.uploadDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return "", "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", "", fmt.Errorf("write file: %w", err)
	}

	url := fmt.Sprintf("%s/uploads/%s", s.baseURL, filename)
	return path, url, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	return os.Remove(path)
}
