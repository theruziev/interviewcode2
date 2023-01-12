package model

import "time"

type ContentStatus string

const (
	ContentStatusActive = ContentStatus("active")
	ContentStatusHide   = ContentStatus("hide")
)

type ContentType string

const (
	ContentTypeMP4   = ContentType("mp4")
	ContentTypeImage = ContentType("image")
	ContentTypeWebM  = ContentType("webm")
)

type PreviewContentType string

const (
	PreviewContentTypeImage = PreviewContentType("image")
	PreviewContentTypeWebM  = PreviewContentType("webm")
)

type Content struct {
	ID          uint64             `db:"id"`
	Source      string             `db:"source"`
	PreviewURL  string             `db:"preview_url"`
	PreviewType PreviewContentType `db:"preview_type"`
	Description string             `db:"description"`
	MediaURL    string             `db:"media_url"`
	Tags        []string           `db:"tags"`
	Type        ContentType        `db:"content_type"`
	Status      ContentStatus      `db:"status"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
	Views       uint64             `db:"views"`
}

type ContentListResponse struct {
	Result []*Content
	LastID uint64
}
