package jobs

import (
	"context"
	"errors"
	"time"
)

type JobType string

const (
	JobTypeThumbnail   JobType = "thumbnail"
	JobTypeMetadata    JobType = "metadata"
	JobTypeProcessing  JobType = "processing"
	JobTypeTranscoding JobType = "transcoding"
)

type JobPayload interface {
	Validate() error
}

type Job struct {
	ID         string          `json:"id"`
	Type       JobType         `json:"type"`
	Payload    JobPayload      `json:"payload"`
	Ctx        context.Context `json:"-"`
	RetryCount int             `json:"retry_count"`
	MaxRetries int             `json:"max_retries"`
	CreatedAt  time.Time       `json:"created_at"`
}

type BasePayload struct {
	FileID     string `json:"file_id"`
	MimeType   string `json:"mime_type"`
	StorageKey string `json:"storage_key"`
	UserID     string `json:"user_id"`
}

type ThumbnailPayload struct {
	*BasePayload
}

type MetadataPayload struct {
	*BasePayload
}

func (p *BasePayload) Validate() error {
	if p.FileID == "" {
		return errors.New("file_id is required")
	}
	if p.MimeType == "" {
		return errors.New("mime_type is required")
	}
	if p.StorageKey == "" {
		return errors.New("storage_key is required")
	}
	if p.UserID == "" {
		return errors.New("user_id is required")
	}
	return nil
}
