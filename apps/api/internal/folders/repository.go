package folders

import (
	"context"
	"errors"

	"github.com/cayke1/mydrive-api/internal/files"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FolderRepository struct {
	db *pgxpool.Pool
}

func NewFolderRepository(db *pgxpool.Pool) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) GetByID(ctx context.Context, id string) (*Folder, error) {
	query := `SELECT id, name, created_at, updated_at, parent_id, owner_id
				FROM folders WHERE id = $1`
	var folder Folder
	err := r.db.QueryRow(ctx, query, id).Scan(
		&folder.ID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
		&folder.ParentID,
		&folder.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepository) GetByIDWithContents(ctx context.Context, id string) (*FolderWithContent, error) {
	folder, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, nil
	}
	children, err := r.GetByParentID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		children = []Folder{}
	}
	fr := files.NewFileRepository(r.db)
	folderFiles, err := fr.GetByFolderID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(folderFiles) == 0 {
		folderFiles = []files.File{}
	}
	return &FolderWithContent{
		Folder:   *folder,
		Children: children,
		Files:    folderFiles,
	}, nil
}

func (r *FolderRepository) GetByParentID(ctx context.Context, parentID string) ([]Folder, error) {
	query := `SELECT id, name, created_at, updated_at, parent_id, owner_id
				FROM folders WHERE parent_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var folder Folder
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt, &folder.ParentID, &folder.OwnerID); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, rows.Err()
}

func (r *FolderRepository) Create(ctx context.Context, folder *Folder) error {
	query := `INSERT INTO folders (id, name, created_at, updated_at, parent_id, owner_id)
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query,
		folder.ID,
		folder.Name,
		folder.CreatedAt,
		folder.UpdatedAt,
		folder.ParentID,
		folder.OwnerID,
	)
	return err
}

func (r *FolderRepository) GetByOwner(ctx context.Context, ownerID string) ([]Folder, error) {
	query := `SELECT id, name, created_at, updated_at, parent_id, owner_id
				FROM folders WHERE owner_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var folder Folder
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt, &folder.ParentID, &folder.OwnerID); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, rows.Err()
}

func (r *FolderRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM folders WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
