package main

import (
	"context"
	"encoding/json"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "foro/db/sqlc"
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

func handleUsuarios(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	queries := sqlc.New(db)
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

func handleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	queries := sqlc.New(db)
	ctx := context.Background()
	if r.URL.Path != "/register" {
			http.NotFound(w, r)
			return
	}
		
		//1. Parsear los datos del formulario (¡Crucial!)
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

func main() {
	// Testeamos que abra la db
	db, err := abrirDB()
	if err != nil {
		fmt.Printf("Error al abrir la base de datos: %s\n", err)
	}
	//cerramos una vez terminada la ejecución del main
	defer db.Close()
	
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	port := ":8080"

	http.Handle("/", fileServer)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, db)
	})
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		handleUsuarios(w, r, db)
	})
	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)
	
	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
