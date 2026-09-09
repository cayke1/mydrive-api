package folders

import (
	"time"

	"github.com/cayke1/mydrive-api/internal/files"
)

type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ParentID  *string   `json:"parent_id,omitempty"`
	OwnerID   string    `json:"owner_id"`
}

type CreateFolderInput struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
	OwnerID  string  `json:"owner_id"`
}

type FolderWithContent struct {
	Folder   Folder       `json:"folder"`
	Children []Folder     `json:"children"`
	Files    []files.File `json:"files"`
}
