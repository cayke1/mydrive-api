package folders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FolderRepository struct {
	db *pgxpool.Pool
}

func NewFolderRepository(db *pgxpool.Pool) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) GetAll(ctx context.Context) ([]Folder, error) {
	query := `SELECT id, name, created_at, updated_at, parent_id, owner_id
                FROM folders ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
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
