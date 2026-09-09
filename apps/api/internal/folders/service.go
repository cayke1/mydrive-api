package folders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FolderService struct {
	repo *FolderRepository
}

func NewFolderService(repo *FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) GetFolders(ctx context.Context) ([]Folder, error) {
	return s.repo.GetAll(ctx)
}

func (s *FolderService) GetFolderByID(ctx context.Context, id string) (*Folder, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FolderService) GetUserFolders(ctx context.Context, userID string) ([]Folder, error) {
	return s.repo.GetByOwner(ctx, userID)
}

func (s *FolderService) CreateFolder(ctx context.Context, input *CreateFolderInput) (*Folder, error) {
	if input.Name == "" {
		return nil, errors.New("folder name is required")
	}
	if len(input.Name) > 255 {
		return nil, errors.New("folder name exceeds maximum length of 255 characters")
	}
	if input.OwnerID == "" {
		return nil, errors.New("owner ID is required")
	}

	if input.ParentID != nil && *input.ParentID != "" {
		parent, err := s.repo.GetByID(ctx, *input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve parent folder: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent folder not found")
		}
	}

	now := time.Now().UTC()
	folder := &Folder{
		ID:        uuid.New().String(),
		Name:      input.Name,
		ParentID:  input.ParentID,
		OwnerID:   input.OwnerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, folder); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return folder, nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, id string) error {
	folder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve folder: %w", err)
	}
	if folder == nil {
		return fmt.Errorf("folder not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *FolderService) GetFolderWithContents(ctx context.Context, id string) (*FolderWithContent, error) {
	return s.repo.GetByIDWithContents(ctx, id)
}
