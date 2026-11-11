package handlers

import (
	"net/http"
	"strconv"
	sqlc "foro/internal/db/sqlc"
	"foro/internal/logic"
	"foro/views"
	"database/sql"
)

func MensajesPageHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Obtiene la lista con ListUsers y renderiza la vista completa

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := queries.ListMensajes(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.MensajesPage(usuarios).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar vista: "+err.Error(), http.StatusInternalServerError)
	}
}

func CrearMensajeWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Crea el mensaje con CreateMensaje y redirige a GET
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	usuario, err := strconv.ParseInt(r.FormValue("usuario"), 10, 32)
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	var respuesta *int32
	if r.FormValue("respuesta") != "" {
		valor, err := strconv.ParseInt(r.FormValue("respuesta"), 10, 32)
		if err != nil {
			http.Error(w, "Respuesta inválida", http.StatusBadRequest)
			return
		}
		respuestaInt := int32(valor)
		respuesta = &respuestaInt
	}

	req := logic.NuevoMensaje{
		IdUsuario:           int32(usuario),
		IdRespuestaAMensaje: respuesta,
		Texto:               r.FormValue("texto"),
	}

	if err := logic.ValidarMensaje(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respuestaSQL := sql.NullInt32{Valid: false}

	if respuesta != nil {
		respuestaSQL = sql.NullInt32{
			Int32: *respuesta,
			Valid: true,
		}
	}

	_, err = queries.CreateMensaje(r.Context(), sqlc.CreateMensajeParams{
		IDUsuario: sql.NullInt32{Int32: req.IdUsuario, Valid: true,},
		IDRespuestaAMensaje: respuestaSQL,
		Texto: req.Texto,
	})
	if err != nil {
		http.Error(w, "Error al crear mensaje: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/mensajes", http.StatusSeeOther)
}

func UsuariosPageHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Obtiene la lista con ListUsers y renderiza la vista completa

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Nombre:      nuevo.Nombre,
		Mail:        nuevo.Mail,
		Contrasenia: nuevo.Contrasenia,
	})
	if err != nil {
		http.Error(w, "Error al crear usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func EditarUsuarioWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	//Actualiza un usuario y redirige a GET
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseInt(r.FormValue("id_usuario"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Error al actualizar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func EliminarUsuarioWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Elimina el usuario y redirige a GET

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseInt(r.FormValue("id_usuario"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	if err := queries.DeleteUser(r.Context(), int32(id64)); err != nil {
		http.Error(w, "Error al eliminar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
