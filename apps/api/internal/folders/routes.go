package folders

import "net/http"

func (fc *FolderController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /folders", fc.CreateFoldersHandler)
	mux.HandleFunc("GET /folders", fc.GetFoldersHandler)
	mux.HandleFunc("GET /folders/{id}", fc.GetFolderByIDHandler)
	mux.HandleFunc("DELETE /folders/{id}", fc.DeleteFolderHandler)
	mux.HandleFunc("GET /folders/{id}/contents", fc.GetFolderWithContentsHandler)
}
