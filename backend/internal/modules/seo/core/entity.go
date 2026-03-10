package core

import "time"

type SiteSettings struct {
	SiteTitle              string    `json:"site_title"`
	SiteDescription        string    `json:"site_description"`
	DefaultMetaTitle       string    `json:"default_meta_title"`
	DefaultMetaDescription string    `json:"default_meta_description"`
	OgImageURL             string    `json:"og_image_url"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type UpdateMetaRequest struct {
	SiteTitle              string `json:"site_title" binding:"required,max=255"`
	SiteDescription        string `json:"site_description" binding:"max=2000"`
	DefaultMetaTitle       string `json:"default_meta_title" binding:"required,max=255"`
	DefaultMetaDescription string `json:"default_meta_description" binding:"max=2000"`
	OgImageURL             string `json:"og_image_url" binding:"omitempty,max=512"`
}

type PostMeta struct {
	Slug           string
	SEOTitle       string
	SEODescription string
	OGImageURL     string
	Summary        string
	PublishedAt    *time.Time
	UpdatedAt      time.Time
}

func (s *SiteSettings) MergeMeta(post PostMeta, baseURL string) map[string]string {
	title := post.SEOTitle
	if title == "" {
		title = s.DefaultMetaTitle
	}
	desc := post.SEODescription
	if desc == "" {
		desc = post.Summary
	}
	if desc == "" {
		desc = s.DefaultMetaDescription
	}
	return map[string]string{"title": title, "description": desc, "url": baseURL}
}
