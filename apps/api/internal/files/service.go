package files

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cayke1/mydrive-api/internal/auth"
	"github.com/cayke1/mydrive-api/internal/storage"
	"github.com/google/uuid"
)

type FileService struct {
	repo    *FileRepository
	storage storage.Storage
}

func NewFileService(repo *FileRepository, storage storage.Storage) *FileService {
	return &FileService{repo: repo, storage: storage}
}

func (s *FileService) GetFilesByFolderID(ctx context.Context, folderID string) ([]File, error) {
	return s.repo.GetByFolderID(ctx, folderID)
}

func (s *FileService) GetFileByID(ctx context.Context, fileID string) (*File, error) {
	return s.repo.GetByID(ctx, fileID)
}

func (s *FileService) UpdateFileMetadata(ctx context.Context, fileID string, metadata string) error {
	file, err := s.verifyOwnership(ctx, fileID)
	if err != nil {
		return err
	}

	file.Metadata = metadata
	file.UpdatedAt = time.Now()

	return s.repo.UpdateMetadata(ctx, file.ID, file.Metadata)
}

func (s *FileService) verifyOwnership(ctx context.Context, fileID string) (*File, error) {
	userID, ok := ctx.Value(auth.UserIdKey).(string)
	if !ok {
		return nil, errors.New("Unauthorized")
	}

	file, err := s.repo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if file.OwnerID != userID {
		return nil, errors.New("Unauthorized")
	}

	return file, nil
}

func (s *FileService) UploadFile(ctx context.Context, file UploadFileInput, fileStream io.Reader) (*File, error) {
	storageKey := uuid.New().String()

	hash := md5.New()
	reader := io.TeeReader(fileStream, hash)

	if err := s.storage.Put(ctx, storageKey, reader, file.Size, file.MimeType); err != nil {
		return nil, err
	}

	newFile := &File{
		ID:         uuid.New().String(),
		Name:       file.Filename,
		FolderID:   file.FolderID,
		OwnerID:    file.OwnerID,
		Size:       file.Size,
		MimeType:   file.MimeType,
		StorageKey: storageKey,
		Checksum:   fmt.Sprintf("%x", hash.Sum(nil)),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, newFile); err != nil {
		s.storage.Delete(ctx, storageKey)
		return nil, err
	}

	return newFile, nil
}

func (s *FileService) DownloadFile(ctx context.Context, fileID string) (*File, io.ReadCloser, error) {
	file, err := s.verifyOwnership(ctx, fileID)
	if err != nil {
		return nil, nil, err
	}

	reader, err := s.storage.Get(ctx, file.StorageKey)
	if err != nil {
		return nil, nil, err
	}

	return file, reader, nil
}

func (s *FileService) DeleteFile(ctx context.Context, fileID string) error {
	file, err := s.verifyOwnership(ctx, fileID)
	if err != nil {
		return err
	}

	// Delete main file
	if err := s.storage.Delete(ctx, file.StorageKey); err != nil {
		return err
	}

	// Delete thumbnail if exists
	thumbnailKey := fmt.Sprintf("%s-thumb.jpg", fileID)
	if err := s.storage.Delete(ctx, thumbnailKey); err != nil {
		// Log but don't fail - thumbnail might not exist
		fmt.Printf("[FILES] Warning: failed to delete thumbnail %s: %v\n", thumbnailKey, err)
	}

	// Delete from database
	if err := s.repo.Delete(ctx, fileID); err != nil {
		return err
	}

	return nil
}
