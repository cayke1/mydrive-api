package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"

	"github.com/cayke1/mydrive-api/internal/auth"
	"github.com/cayke1/mydrive-api/internal/storage"
	"github.com/disintegration/imaging"
)

type ThumbnailHandler struct {
	storage storage.Storage
}

type MetadataHandler struct {
	storage     storage.Storage
	fileService FileService
}

func NewThumbnailHandler(storage storage.Storage) *ThumbnailHandler {
	return &ThumbnailHandler{storage: storage}
}

func NewMetadataHandler(storage storage.Storage, fileService FileService) *MetadataHandler {
	return &MetadataHandler{storage: storage, fileService: fileService}
}

const (
	Width  = 200
	Height = 200
)

func isImageMimeType(mimeType string) bool {
	return mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" || mimeType == "image/bmp" || mimeType == "image/webp"
}

func isVideoMimeType(mimeType string) bool {
	return mimeType == "video/mp4" || mimeType == "video/avi" || mimeType == "video/mov" || mimeType == "video/mpeg" || mimeType == "video/webm"
}

func resizeImageFit(img image.Image, maxWidth, maxHeight int) image.Image {
	return imaging.Fit(img, maxWidth, maxHeight, imaging.MitchellNetravali)
}

func (h *ThumbnailHandler) Handle(ctx context.Context, job *Job) error {
	payload, ok := job.Payload.(*ThumbnailPayload)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}

	if err := payload.Validate(); err != nil {
		return err
	}

	if !isImageMimeType(payload.MimeType) {
		log.Printf("[JOBS->Thumbnail] Skipping thumbnail generation for non-image file: %s", payload.FileID)
		return nil // Skip thumbnail generation for non-image files
	}

	reader, err := h.storage.Get(ctx, payload.StorageKey)
	if err != nil {
		return fmt.Errorf("failed to get file from storage: %w", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	thumbnail := resizeImageFit(img, Width, Height)
	thumbnailKey := fmt.Sprintf("%s-thumb.jpg", payload.FileID)
	buff := new(bytes.Buffer)

	if err := jpeg.Encode(buff, thumbnail, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	if err := h.storage.Put(ctx, thumbnailKey, buff, int64(buff.Len()), "image/jpeg"); err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	log.Printf("[JOBS->Thumbnail] Generated thumbnail for %s", payload.FileID)

	return nil
}

func (h *MetadataHandler) Handle(ctx context.Context, job *Job) error {
	payload, ok := job.Payload.(*MetadataPayload)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}
	if err := payload.Validate(); err != nil {
		return fmt.Errorf("payload validation failed: %w", err)
	}

	log.Printf("[JOBS->Metadata] Processing metadata for file: %s (type: %s)", payload.FileID, payload.MimeType)

	metadata := make(map[string]interface{})

	if isImageMimeType(payload.MimeType) {
		imgMetadata, err := h.extractImageMetadata(ctx, payload.StorageKey)
		if err != nil {
			return fmt.Errorf("failed to extract image metadata: %w", err)
		}
		metadata = imgMetadata
	} else if isVideoMimeType(payload.MimeType) {
		log.Printf("[JOBS->Metadata] Video metadata extraction not implemented for file: %s", payload.FileID)
		metadata["type"] = "video"
	} else {
		log.Printf("[JOBS->Metadata] Unsupported file type for metadata extraction: %s", payload.MimeType)
		return nil // Skip metadata extraction for unsupported file types
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	log.Printf("[JOBS->Metadata] Updating metadata in DB for file %s: %s", payload.FileID, string(metadataJSON))

	// Create a new context with userID so UpdateFileMetadata can verify ownership
	// Import auth package at the top of this file
	ctxWithUser := context.WithValue(ctx, auth.UserIdKey, payload.UserID)

	if err := h.fileService.UpdateFileMetadata(ctxWithUser, payload.FileID, string(metadataJSON)); err != nil {
		return fmt.Errorf("failed to update file metadata in DB: %w", err)
	}

	log.Printf("[JOBS->Metadata] ✓ Successfully updated metadata for file: %s", payload.FileID)

	return nil
}

func (h *MetadataHandler) extractImageMetadata(ctx context.Context, storageKey string) (map[string]interface{}, error) {
	reader, err := h.storage.Get(ctx, storageKey)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	return map[string]interface{}{
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
		"type":   "image",
	}, nil
}

// CompositeJobHandler routes jobs to the appropriate handler based on job type
type CompositeJobHandler struct {
	ThumbnailHandler *ThumbnailHandler
	MetadataHandler  *MetadataHandler
}

func (c *CompositeJobHandler) Handle(ctx context.Context, job *Job) error {
	switch job.Type {
	case JobTypeThumbnail:
		return c.ThumbnailHandler.Handle(ctx, job)
	case JobTypeMetadata:
		return c.MetadataHandler.Handle(ctx, job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}
