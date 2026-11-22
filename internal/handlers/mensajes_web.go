package handlers

import (
	"net/http"
	"strconv"
	sqlc "foro/internal/db/sqlc"
	"foro/internal/logic"
	"foro/views"
	"database/sql"
	"strings"
)

func MensajesPageHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Obtiene la lista con ListUsers y renderiza la vista completa

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	mensajes, err := queries.ListMensajes(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener mensajes", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.MensajesPage(mensajes).Render(r.Context(), w); err != nil {
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

	//parseamos sin preguntar si es vacio ya que es requerido en el formulario
	usuario, err := strconv.ParseInt(r.FormValue("usuario"), 10, 32)
	if err != nil {
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusOK)
			views.Alerta("ID de usuario invalido").Render(r.Context(), w)
			return
	}

	//parsear el campo respuesta, que es opcional
	var respuesta *int32
	if r.FormValue("respuesta") != "" {
		valor, err := strconv.ParseInt(r.FormValue("respuesta"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			mensajes, _ := queries.ListMensajes(r.Context())
			views.MensajesPage(mensajes).Render(r.Context(), w)
			return
		}
		respuestaInt := int32(valor)
		respuesta = &respuestaInt
	}

	respuestaSQL := sql.NullInt32{Valid: false}
	if respuesta != nil {
		respuestaSQL = sql.NullInt32{
			Int32: *respuesta,
			Valid: true,
		}
	}

	req := logic.NuevoMensaje{
		IdUsuario:           int32(usuario),
		IdRespuestaAMensaje: respuesta,
		Texto:               r.FormValue("texto"),
	}

	if err := logic.ValidarMensaje(req); err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("Error al crear mensaje").Render(r.Context(), w)
		return
	}

	_, err = queries.CreateMensaje(r.Context(), sqlc.CreateMensajeParams{
		IDUsuario: sql.NullInt32{Int32: req.IdUsuario, Valid: true,},
		IDRespuestaAMensaje: respuestaSQL,
		Texto: req.Texto,
	})
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("Error al crear mensaje").Render(r.Context(), w)
		return
	}

	mensajes, err := queries.ListMensajes(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	views.MensajeList(mensajes).Render(r.Context(), w)
}

func EditarMensajeWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	//Actualiza un usuario y redirige a GET
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseInt(r.FormValue("id_mensaje"), 10, 32)
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("ID de mensaje invalido").Render(r.Context(), w)
		return
	}

	var usuario *int32
	if r.FormValue("id_usuario") != "" {
		valor, err := strconv.ParseInt(r.FormValue("id_usuario"), 10, 32)
		if err != nil {
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusOK)
			views.Alerta("ID de usuario invalido").Render(r.Context(), w)
			return
		}
		usuarioInt := int32(valor)
		usuario = &usuarioInt
	}

	usuarioSQL := sql.NullInt32{Valid: false}
	if usuario != nil {
		usuarioSQL = sql.NullInt32{
			Int32: *usuario,
			Valid: true,
		}
	}

	var respuesta *int32
	if r.FormValue("id_respuesta_a_mensaje") != "" {
		valor, err := strconv.ParseInt(r.FormValue("id_respuesta_a_mensaje"), 10, 32)
		if err != nil {
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusOK)
			views.Alerta("ID de respuesta a mensaje invalido").Render(r.Context(), w)
			return
		}
		respuestaInt := int32(valor)
		respuesta = &respuestaInt
	}

	respuestaSQL := sql.NullInt32{Valid: false}
	if respuesta != nil {
		respuestaSQL = sql.NullInt32{
			Int32: *respuesta,
			Valid: true,
		}
	}

	var meGusta *int32
	if r.FormValue("me_gusta") != "" {
		valor, err := strconv.ParseInt(r.FormValue("me_gusta"), 10, 32)
		if err != nil {
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusOK)
			views.Alerta("Cantidad de me gusta invalido").Render(r.Context(), w)
			return
		}
		meGustaInt := int32(valor)
		meGusta = &meGustaInt
	}

	meGustaSQL := sql.NullInt32{Valid: false}
	if meGusta != nil {
		meGustaSQL = sql.NullInt32{
			Int32: *meGusta,
			Valid: true,
		}
	}

	texto := r.FormValue("texto")

	reqUpdate := logic.EditarMensaje{
		IdUsuario: usuario,
		IdRespuestaAMensaje: respuesta,
		Texto: &texto,
		MeGusta: meGusta,
	}

	if err := logic.ValidarEdicionMensaje(reqUpdate); err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("Error al actualizar mensaje").Render(r.Context(), w)
		return
	}

	err = queries.UpdateMensaje(r.Context(), sqlc.UpdateMensajeParams{
		IDMensaje: int32(id64),
		IDUsuario: usuarioSQL,
		IDRespuestaAMensaje: respuestaSQL,
		Texto: texto,
		MeGusta: meGustaSQL,
	})

	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("Error al actualizar mensaje").Render(r.Context(), w)
		return
	}
	
	mensajes, err := queries.ListMensajes(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	views.MensajeList(mensajes).Render(r.Context(), w)
}

func EliminarMensajeWebHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	// Elimina el mensaje

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

	if err := queries.DeleteMensaje(r.Context(), int32(id64)); err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusOK)
		views.Alerta("Error al eliminar mensaje").Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}