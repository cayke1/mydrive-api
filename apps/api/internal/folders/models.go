package folders

import "time"

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
