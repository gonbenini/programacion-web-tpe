package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func abrirDB() (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		"postgres://dbuser:dbpw@localhost:5432/foro?sslmode=disable",
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
	abrirDB() // Testea que abra la db
	// db, err := abrirDB() // en un futuro usamos esta linea para pasarle a los handles

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	port := ":8080"

	http.Handle("/", fileServer)
	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
