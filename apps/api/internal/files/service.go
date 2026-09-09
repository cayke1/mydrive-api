package files

import "context"

type FileService struct {
	repo *FileRepository
}

func NewFileService(repo *FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) GetFilesByFolderID(ctx context.Context, folderID string) ([]File, error) {
	return s.repo.GetByFolderID(ctx, folderID)
}
