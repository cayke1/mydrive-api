package files

import (
	"encoding/json"
	"net/http"
)

type FileController struct {
	service *FileService
}

func NewFileController(service *FileService) *FileController {
	return &FileController{service: service}
}

func (fc *FileController) GetFilesByFolderIDHandler(w http.ResponseWriter, r *http.Request) {
	folderID := r.PathValue("folder_id")
	if folderID == "" {
		http.Error(w, "Missing folder ID", http.StatusBadRequest)
		return
	}
	files, err := fc.service.GetFilesByFolderID(r.Context(), folderID)
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}
