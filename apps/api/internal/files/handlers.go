package files

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cayke1/mydrive-api/internal/auth"
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

func (fc *FileController) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	folderID := r.FormValue("folder_id")
	if folderID == "" {
		http.Error(w, "Missing folder_id", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(auth.UserIdKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileInput := UploadFileInput{
		FolderID: folderID,
		OwnerID:  userID,
		Filename: header.Filename,
		MimeType: header.Header.Get("Content-Type"),
		Size:     header.Size,
	}

	uploadedFile, err := fc.service.UploadFile(
		r.Context(),
		fileInput,
		file,
	)

	if err != nil {
		log.Printf("Upload error %w", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uploadedFile)
}

func (fc *FileController) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("file_id")
	if fileID == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
	}

	log.Printf("Downloading file: %s", fileID)

	file, reader, err := fc.service.DownloadFile(r.Context(), fileID)
	if err != nil {
		log.Printf("❌ Download error: %v", err)
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file.Name))

	log.Printf("📤 Streaming %s (%d bytes)", file.Name, file.Size)

	if _, err := io.Copy(w, reader); err != nil {
		log.Printf("Copy error: %v", err)
		return
	}

	log.Printf("Download complete: %s", file.Name)
}

func (fc *FileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("file_id")
	if fileID == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	if err := fc.service.DeleteFile(r.Context(), fileID); err != nil {
		log.Printf("❌ Delete error: %v", err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
		"file_id": fileID,
	})
}
