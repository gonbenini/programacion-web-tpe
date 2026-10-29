package handlers

import (
	"net/http"

	sqlc "foro/internal/db/sqlc"
	"foro/views"
)

func UsuariosPageHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.UsuariosPage(usuarios).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar vista: "+err.Error(), http.StatusInternalServerError)
	}
}
