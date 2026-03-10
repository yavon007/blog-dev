package core

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	GetSiteSettings(ctx context.Context) (*SiteSettings, error)
	UpsertSiteSettings(ctx context.Context, req UpdateMetaRequest) (*SiteSettings, error)
	ListPublishedPosts(ctx context.Context, limit int) ([]PostMeta, error)
}

type Service struct {
	repo    Repository
	baseURL string
}

func NewService(repo Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: sanitizeBaseURL(baseURL),
	}
}

func (s *Service) GetMeta(ctx context.Context) (*SiteSettings, error) {
	settings, err := s.repo.GetSiteSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	return settings, nil
}

func (s *Service) UpdateMeta(ctx context.Context, req UpdateMetaRequest) (*SiteSettings, error) {
	settings, err := s.repo.UpsertSiteSettings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update site settings: %w", err)
	}
	return settings, nil
}

func (s *Service) GenerateSitemap(ctx context.Context) (string, error) {
	posts, err := s.repo.ListPublishedPosts(ctx, 0)
	if err != nil {
		return "", fmt.Errorf("list posts for sitemap: %w", err)
	}

	urlset := urlSet{
		XMLName: xml.Name{Local: "urlset"},
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []urlEntry{
			{
				Loc:     s.baseURL,
				LastMod: time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	for _, post := range posts {
		lastMod := post.UpdatedAt.UTC()
		if post.PublishedAt != nil {
			lastMod = post.PublishedAt.UTC()
		}
		urlset.URLs = append(urlset.URLs, urlEntry{
			Loc:     fmt.Sprintf("%s/post/%s", s.baseURL, post.Slug),
			LastMod: lastMod.Format(time.RFC3339),
		})
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(urlset); err != nil {
		return "", fmt.Errorf("encode sitemap xml: %w", err)
	}
	return buf.String(), nil
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

type urlEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

func sanitizeBaseURL(base string) string {
	base = strings.TrimSpace(base)
	base = strings.TrimRight(base, "/")
	if base == "" {
		return "http://localhost:8080"
	}
	return base
}
