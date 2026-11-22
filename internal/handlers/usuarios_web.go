package handlers

import (
	"net/http"
	"strconv"
	sqlc "foro/internal/db/sqlc"
	"foro/internal/logic"
	"foro/views"
	"strings"
)

func UsuariosPageHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Obtiene la lista con ListUsers y renderiza la vista completa

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.UsuariosPage(usuarios, "").Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar vista: "+err.Error(), http.StatusInternalServerError)
	}
}

func CrearUsuarioWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Crea el usuario con CreateUser y redirige a GET

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	nuevo := logic.NuevoUsuario{
		Nombre:      r.FormValue("nombre"),
		Mail:        r.FormValue("mail"),
		Contrasenia: r.FormValue("contrasenia"),
	}

	if err := logic.ValidarUsuario(nuevo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "Usuario inválido - " + err.Error()).Render(r.Context(), w)
		return
	}

	_, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Nombre:      nuevo.Nombre,
		Mail:        nuevo.Mail,
		Contrasenia: nuevo.Contrasenia,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "Error al crear usuario - " + err.Error()).Render(r.Context(), w)
		return
	}
	users, err := queries.ListUsers(r.Context())
	views.UsuarioList(users).Render(r.Context(), w)
}

func EditarUsuarioWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	//Actualiza un usuario y redirige a GET
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseInt(r.FormValue("id_usuario"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "ID de usuario inválido - " + err.Error()).Render(r.Context(), w)
		return
	}

	nombre := r.FormValue("nombre")
	mail := r.FormValue("mail")
	contrasenia := r.FormValue("contrasenia")

	reqUpdate := logic.UpdateUsuario{
		Nombre: &nombre,
		Mail:   &mail,
	}
	if contrasenia != "" {
		reqUpdate.Contrasenia = &contrasenia
	}

	if err := logic.ValidarUpdateUsuario(reqUpdate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "Usuario inválido - " + err.Error()).Render(r.Context(), w)
		return
	}

	// Si no ingresó contraseña nueva en el form de edición, usamos UpdateUserPartial
	if contrasenia == "" {
		err = queries.UpdateUserPartial(r.Context(), sqlc.UpdateUserPartialParams{
			IDUsuario: int32(id64),
			Nombre:    nombre,
			Mail:      mail,
		})
	} else {
		err = queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
			IDUsuario:   int32(id64),
			Nombre:      nombre,
			Mail:        mail,
			Contrasenia: contrasenia,
		})
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "Error al actualizar usuario - " + err.Error()).Render(r.Context(), w)
		return
	}

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	views.UsuarioList(users).Render(r.Context(), w)
}

func EliminarUsuarioWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Elimina el usuario y redirige a GET

	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id64 ,err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := queries.DeleteUser(r.Context(), int32(id64)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		usuarios, _ := queries.ListUsers(r.Context())
    	views.UsuariosPage(usuarios, "Error al eliminar usuario - " + err.Error()).Render(r.Context(), w)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
