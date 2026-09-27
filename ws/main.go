package main

import (
	"database/sql"
	"fmt"
	"net/http"

	sqlc "foro/db/sqlc"

	"foro/controller"

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
		controller.UsuariosAPIHandler(w, r, queries)
	})

	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		controller.UsuarioAPIHandler(w, r, queries)
	})

	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
