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
