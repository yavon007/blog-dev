package core

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	seocore "github.com/yavon007/blog-dev/backend/internal/modules/seo/core"
)

type Repository interface {
	ListPublishedPosts(ctx context.Context, limit int) ([]seocore.PostMeta, error)
	GetSiteSettings(ctx context.Context) (*seocore.SiteSettings, error)
}

type Service struct {
	repo    Repository
	baseURL string
	limit   int
}

func NewService(repo Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: sanitizeBaseURL(baseURL),
		limit:   50,
	}
}

func (s *Service) GenerateFeed(ctx context.Context) (string, error) {
	posts, err := s.repo.ListPublishedPosts(ctx, s.limit)
	if err != nil {
		return "", fmt.Errorf("list posts for rss: %w", err)
	}
	settings, err := s.repo.GetSiteSettings(ctx)
	if err != nil {
		return "", fmt.Errorf("load site settings: %w", err)
	}

	channel := rssChannel{
		Title:       settings.SiteTitle,
		Link:        s.baseURL,
		Description: settings.SiteDescription,
		PubDate:     time.Now().UTC().Format(time.RFC1123Z),
	}

	for _, post := range posts {
		pubDate := post.UpdatedAt
		if post.PublishedAt != nil {
			pubDate = post.PublishedAt.UTC()
		} else {
			pubDate = post.UpdatedAt.UTC()
		}
		desc := post.SEODescription
		if desc == "" {
			desc = post.Summary
		}
		item := rssItem{
			Title:       chooseNonEmpty(post.SEOTitle, settings.DefaultMetaTitle),
			Link:        fmt.Sprintf("%s/post/%s", s.baseURL, post.Slug),
			Description: desc,
			GUID:        fmt.Sprintf("%s/post/%s", s.baseURL, post.Slug),
			PubDate:     pubDate.Format(time.RFC1123Z),
		}
		channel.Items = append(channel.Items, item)
	}

	doc := rssDocument{
		Version: "2.0",
		Channel: channel,
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return "", fmt.Errorf("encode rss xml: %w", err)
	}
	return buf.String(), nil
}

type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func sanitizeBaseURL(base string) string {
	base = strings.TrimSpace(base)
	base = strings.TrimRight(base, "/")
	if base == "" {
		return "http://localhost:8080"
	}
	return base
}

func chooseNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
