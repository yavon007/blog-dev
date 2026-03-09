package core

import "time"

type MediaFile struct {
	ID           int64      `json:"id"`
	Filename     string     `json:"filename"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	Size         int64      `json:"size"`
	Width        *int       `json:"width"`
	Height       *int       `json:"height"`
	AltText      string     `json:"alt_text"`
	Storage      string     `json:"storage"`
	Path         string     `json:"path"`
	URL          string     `json:"url"`
	UploadedBy   *int64     `json:"uploaded_by"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type ListMediaFilter struct {
	MimeType string
	Page     int
	PageSize int
}

type UploadResult struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Filename string `json:"filename"`
}
