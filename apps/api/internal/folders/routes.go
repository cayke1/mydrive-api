package folders

import (
	"net/http"

	"github.com/cayke1/mydrive-api/internal/auth"
)

func (fc *FolderController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /folders", auth.Authorize(http.HandlerFunc(fc.CreateFoldersHandler)))
	mux.Handle("GET /folders", auth.Authorize(http.HandlerFunc(fc.GetFoldersHandler)))
	mux.HandleFunc("GET /folders/{id}", fc.GetFolderByIDHandler)
	mux.HandleFunc("DELETE /folders/{id}", fc.DeleteFolderHandler)
	mux.HandleFunc("GET /folders/{id}/contents", fc.GetFolderWithContentsHandler)
}
