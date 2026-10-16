package jobs

import "context"

type FileService interface {
	GetFileByID(ctx context.Context, fileID string) (interface{}, error)
	UpdateFileMetadata(ctx context.Context, fileID string, metadata string) error
}
