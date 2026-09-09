package files

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *File) error {
	query := `INSERT INTO files (id, name, folder_id, owner_id, size, 
	mime_type, storage_key, checksum, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, query,
		file.ID,
		file.Name,
		file.FolderID,
		file.OwnerID,
		file.Size,
		file.MimeType,
		file.StorageKey,
		file.Checksum,
		file.CreatedAt,
		file.UpdatedAt,
	)
	return err
}

func (r *FileRepository) GetByFolderID(ctx context.Context, folderID string) ([]File, error) {
	query := `SELECT id, name, folder_id, owner_id, size, mime_type, 
				storage_key, checksum, created_at, updated_at
				FROM files WHERE folder_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.FolderID,
			&file.OwnerID,
			&file.Size,
			&file.MimeType,
			&file.StorageKey,
			&file.Checksum,
			&file.CreatedAt,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}
