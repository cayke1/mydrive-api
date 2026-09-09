package files

import "net/http"

func (fc *FileController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /file/{folder_id}", fc.GetFilesByFolderIDHandler)
}
