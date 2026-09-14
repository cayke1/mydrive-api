package files

import "time"

type File struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	FolderID   string    `json:"folder_id"`
	OwnerID    string    `json:"owner_id"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	StorageKey string    `json:"storage_key"`
	Checksum   string    `json:"checksum"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateFileInput struct {
	Name     string `json:"name"`
	FolderID string `json:"folder_id"`
	OwnerID  string `json:"owner_id"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

type UploadFileInput struct {
	FolderID string `json:"folder_id"`
	OwnerID  string `json:"owner_id"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	Filename string `json:"filename"`
}
