package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sqlc "foro/db/sqlc"
	"foro/logic"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type Respuesta struct {
	Mensaje string `json:"mensaje"`
}

func abrirDB() (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		"postgres://dbuser:dbpw@db:5432/foro?sslmode=disable",
	)
	if err != nil {
		return nil, err
	}

	// Verifica que la base de datos responda
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Configura el tamaño máximo del pool
	db.SetMaxOpenConns(25)

	return db, nil
}

func handleUsuarios(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	ctx := context.Background()

	if r.URL.Path != "/usuarios" {
		http.NotFound(w, r)
		return
	}

	users, err := queries.ListUsers(ctx)
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Usuarios</title></head><body><h1>Lista de Usuarios</h1><ul>")
	for _, user := range users {
		fmt.Fprintf(w, "<li>ID: %d, Nombre: %s, Mail: %s</li>", user.IDUsuario, user.Nombre, user.Mail)
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

func handleRegister(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {

	ctx := context.Background()
	if r.URL.Path != "/register" {
		http.NotFound(w, r)
		return
	}

	//1. Parsear los datos del formulario
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear", http.StatusBadRequest)
		return
	}
	username := r.FormValue("nombre")
	mail := r.FormValue("mail")
	password := r.FormValue("contrasenia")

	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Nombre:      username,
			Mail:        mail,
			Contrasenia: password,
		})
	if err != nil {
		fmt.Printf("Error al crear usuario: %s\n", err)
		return
	}
	fmt.Printf("Usuario creado: %v\n", createdUser)

	// 1. Configuramos la cabecera para responder en formato JSON
	w.Header().Set("Content-Type", "application/json")

	// 2. Creamos la respuesta de éxito
	respuestaExito := Respuesta{
		Mensaje: "Se envió correctamente",
	}

	// 3. Enviamos el JSON de vuelta a JavaScript
	json.NewEncoder(w).Encode(respuestaExito)
}

func usuariosAPIHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	switch r.Method {
	case http.MethodPost:
		createUsuario(w, r, queries)

	case http.MethodGet:
		getUsuarios(w, r, queries)

	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func usuarioAPIHandler(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	//extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "URL invalido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUsuario(w, r, queries, int32(id))

	case http.MethodPut:
		updateUsuario(w, r, queries, int32(id))

	case http.MethodDelete:
		deleteUsuario(w, r, queries, int32(id))

	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func createUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	var req logic.NuevoUsuario

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al procesar JSON", http.StatusBadRequest)
		return
	}

	err = logic.ValidarUsuario(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Nombre:      req.Nombre,
			Mail:        req.Mail,
			Contrasenia: req.Contrasenia,
		})
	if err != nil {
		fmt.Printf("Error al crear usuario: %s\n", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func getUsuarios(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	ctx := context.Background()

	users, err := queries.ListUsers(ctx) // Buscamos todos los registros de usuario
	if err != nil {
		http.Error(w, "Error al listar usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos array JSON
	json.NewEncoder(w).Encode(users)
}

func getUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()

	user, err := queries.GetUserById(ctx, id) // Buscamos el registro
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound) // Si no existe devolvemos un 404 Not Found
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos JSON con estado 200
	json.NewEncoder(w).Encode(user)
}

func updateUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {

}

func deleteUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()
	err := queries.DeleteUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Testeamos que abra la db
	db, err := abrirDB()
	if err != nil {
		fmt.Printf("Error al abrir la base de datos: %s\n", err)
	}
	//cerramos una vez terminada la ejecución del main
	defer db.Close()

	queries := sqlc.New(db)

	port := ":8080"

	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		usuariosAPIHandler(w, r, queries)
	})

	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		usuarioAPIHandler(w, r, queries)
	})

	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
