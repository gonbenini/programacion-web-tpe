package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "foro/db/sqlc"
)

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
	//abrirDB() // Testea que abra la db
	db, err := abrirDB() // en un futuro usamos esta linea para pasarle a los handles
	if err != nil {
		fmt.Printf("Error al abrir la base de datos: %s\n", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Nombre:   "1",
			Mail: "John Doe",
			Contrasenia: "password123",
	})
	if err != nil {
		fmt.Printf("Error al crear usuario: %s\n", err)
		return
	}
	fmt.Printf("Usuario creado: %v\n", createdUser)


	getUsers, err := queries.GetUsers(ctx)
	#if err != nil {
		fmt.Printf("Error al obtener usuarios: %s\n", err)
		return
	}
	fmt.Printf("Usuarios obtenidos: %v\n", getUsers)
	

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	port := ":8080"

	http.Handle("/", fileServer)
	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
