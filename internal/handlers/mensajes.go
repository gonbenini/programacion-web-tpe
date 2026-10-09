package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sqlc "foro/internal/db/sqlc"
	"foro/internal/logic"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

func MensajesAPIHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	switch r.Method {
	case http.MethodPost:
		createMensaje(w, r, queries)

	case http.MethodGet:
		getMensajes(w, r, queries)

	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func MensajeAPIHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	//extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "URL invalido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getMensaje(w, r, queries, int32(id))

	case http.MethodPut:
		updateMensaje(w, r, queries, int32(id))

	case http.MethodDelete:
		deleteMensaje(w, r, queries, int32(id))

	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func createMensaje(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	var req logic.NuevoMensaje

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al procesar JSON", http.StatusBadRequest)
		return
	}

	err = logic.ValidarMensaje(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var respuestaID sql.NullInt32
	if req.IdRespuestaAMensaje != nil {
		respuestaID = sql.NullInt32{Int32: *req.IdRespuestaAMensaje, Valid: true}
	}

	ctx := context.Background()
	createdMensaje, err := queries.CreateMensaje(ctx,
		sqlc.CreateMensajeParams{
			IDUsuario:           sql.NullInt32{Int32: req.IdUsuario, Valid: true},
			IDRespuestaAMensaje: respuestaID,
			Texto:               req.Texto,
		})
	if err != nil {
		fmt.Printf("Error al crear mensaje: %s\n", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMensaje)
}

func getMensajes(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	ctx := context.Background()

	mensajes, err := queries.ListMensajes(ctx) // Buscamos todos los registros de mensaje
	if err != nil {
		http.Error(w, "Error al listar mensajes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos array JSON
	json.NewEncoder(w).Encode(mensajes)
}

func getMensaje(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()

	mensaje, err := queries.GetMensajeById(ctx, id) // Buscamos el registro
	if err != nil {
		http.Error(w, "Mensaje no encontrado", http.StatusNotFound) // Si no existe devolvemos un 404 Not Found
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos JSON con estado 200
	json.NewEncoder(w).Encode(mensaje)
}

func updateMensaje(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	var req logic.EditarMensaje

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al procesar JSON", http.StatusBadRequest)
		return
	}

	err = logic.ValidarEdicionMensaje(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// 1. Buscamos el mensaje original
	mensajeActual, err := queries.GetMensajeById(ctx, id)
	if err != nil {
		http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
		return
	}

	// 2. Validamos que las nuevas dependencias existan en la BD (si es que las enviaron)
	if req.IdUsuario != nil {
		_, err := queries.GetUserById(ctx, *req.IdUsuario)
		if err != nil {
			http.Error(w, "El id_usuario especificado no existe", http.StatusBadRequest)
			return
		}
	}

	if req.IdRespuestaAMensaje != nil {
		_, err := queries.GetMensajeById(ctx, *req.IdRespuestaAMensaje)
		if err != nil {
			http.Error(w, "El id_respuesta_a_mensaje especificado no existe", http.StatusBadRequest)
			return
		}
	}

	// 3. Preparamos las variables con los datos viejos
	nuevoIdUsuario := mensajeActual.IDUsuario
	nuevoIdRespuesta := mensajeActual.IDRespuestaAMensaje
	nuevoTexto := mensajeActual.Texto
	nuevoMeGusta := mensajeActual.MeGusta

	// 4. Reemplazamos SOLO si en el JSON enviaron un dato nuevo
	if req.IdUsuario != nil {
		nuevoIdUsuario = sql.NullInt32{Int32: *req.IdUsuario, Valid: true}
	}
	if req.IdRespuestaAMensaje != nil {
		nuevoIdRespuesta = sql.NullInt32{Int32: *req.IdRespuestaAMensaje, Valid: true}
	}
	if req.Texto != nil {
		nuevoTexto = *req.Texto
	}
	if req.MeGusta != nil {
		nuevoMeGusta = sql.NullInt32{Int32: *req.MeGusta, Valid: true}
	}

	// 5. Ejecutamos el Update
	err = queries.UpdateMensaje(ctx, sqlc.UpdateMensajeParams{
		IDMensaje:           id,
		IDUsuario:           nuevoIdUsuario,
		IDRespuestaAMensaje: nuevoIdRespuesta,
		Texto:               nuevoTexto,
		MeGusta:             nuevoMeGusta,
	})
	if err != nil {
		fmt.Printf("Error al actualizar mensaje: %s\n", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteMensaje(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()
	err := queries.DeleteMensaje(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
