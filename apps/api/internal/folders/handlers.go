package folders

import (
	"encoding/json"
	"net/http"

	"github.com/cayke1/mydrive-api/internal/auth"
)

type FolderController struct {
	service *FolderService
}

func NewFolderController(service *FolderService) *FolderController {
	return &FolderController{service: service}
}

func (fc *FolderController) GetFoldersHandler(w http.ResponseWriter, r *http.Request) {
	ownerId := r.Context().Value(auth.UserIdKey).(string)
	folders, err := fc.service.GetFolders(r.Context(), ownerId)
	if err != nil {
		http.Error(w, "Failed to retrieve folders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(folders)
}

func (fc *FolderController) GetFolderByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing folder ID", http.StatusBadRequest)
		return
	}

	folder, err := fc.service.GetFolderByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve folder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(folder)
}

func (fc *FolderController) CreateFoldersHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateFolderInput

	ownerID := r.Context().Value(auth.UserIdKey).(string)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	folder, err := fc.service.CreateFolder(r.Context(), &input, ownerID)
	if err != nil {
		http.Error(w, "Failed to create folder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func (fc *FolderController) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing folder ID", http.StatusBadRequest)
		return
	}
	err := fc.service.DeleteFolder(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete folder: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (fc *FolderController) GetFolderWithContentsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing folder ID", http.StatusBadRequest)
		return
	}
	folderWithContents, err := fc.service.GetFolderWithContents(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve folder with contents: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(folderWithContents)
}
