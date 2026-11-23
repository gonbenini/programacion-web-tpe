package main

import (
	"fmt"
	"net/http"
	"os"

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
		os.Exit(1)
	}
	//cerramos una vez terminada la ejecución del main
	defer db_conecction.Close()

	queries := sqlc.New(db_conecction)

	port := ":8080"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// "/" atrapa cualquier ruta no registrada: solo la raíz muestra la página de usuarios
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.UsuariosPageHandler(w, r, queries)
	})

	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.UsuariosPageHandler(w, r, queries)
		case http.MethodPost:
			handlers.CrearUsuarioWebHandler(w, r, queries)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handlers.EliminarUsuarioWebHandler(w, r, queries)
		case http.MethodPut:
			handlers.EditarUsuarioWebHandler(w, r, queries)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/mensajes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.MensajesPageHandler(w, r, queries)
		case http.MethodPost:
			handlers.CrearMensajeWebHandler(w, r, queries)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/mensajes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handlers.EliminarMensajeWebHandler(w, r, queries)
		case http.MethodPut:
			handlers.EditarMensajeWebHandler(w, r, queries)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
