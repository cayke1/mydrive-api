package files

import "net/http"

func (fc *FileController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /file/{folder_id}", fc.GetFilesByFolderIDHandler)
	mux.HandleFunc("POST /files/upload", fc.UploadFileHandler)
	mux.HandleFunc("GET /files/{file_id}/download", fc.DownloadFileHandler)
	mux.HandleFunc("DELETE /files/{file_id}", fc.DeleteFileHandler)
}
