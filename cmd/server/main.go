package main

import (
	"fmt"
	"net/http"

	"foro/internal/db"
	"foro/internal/handlers"

	sqlc "foro/internal/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type Respuesta struct {
	Mensaje string `json:"mensaje"`
}

func main() {
	// Testeamos que abra la db
	db_conecction, err := db.AbrirDB()
	if err != nil {
		fmt.Printf("Error al abrir la base de datos: %s\n", err)
	}
	//cerramos una vez terminada la ejecución del main
	defer db_conecction.Close()

	queries := sqlc.New(db_conecction)

	port := ":8080"

	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		handlers.UsuariosAPIHandler(w, r, queries)
	})

	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		handlers.UsuarioAPIHandler(w, r, queries)
	})

	http.HandleFunc("/api/mensajes", func(w http.ResponseWriter, r *http.Request) {
		handlers.MensajesAPIHandler(w, r, queries)
	})

	http.HandleFunc("/api/mensajes/", func(w http.ResponseWriter, r *http.Request) {
		handlers.MensajeAPIHandler(w, r, queries)
	})

	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
