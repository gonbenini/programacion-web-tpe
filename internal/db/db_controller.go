package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

// cadenaConexion arma la URL de conexión a PostgreSQL a partir de variables de entorno
func cadenaConexion() (string, error) {
	usuario := os.Getenv("POSTGRES_USER")
	contrasenia := os.Getenv("POSTGRES_PASSWORD")
	base := os.Getenv("POSTGRES_DB")
	host := os.Getenv("DB_HOST")

	if usuario == "" || contrasenia == "" || base == "" || host == "" {
		return "", fmt.Errorf("faltan variables de entorno de la base (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, DB_HOST): revisar el archivo .env")
	}

	// url.URL escapa caracteres especiales de la contraseña (por ejemplo @ o /)
	conexion := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(usuario, contrasenia),
		Host:     host + ":5432",
		Path:     base,
		RawQuery: "sslmode=disable",
	}
	return conexion.String(), nil
}

func AbrirDB() (*sql.DB, error) {
	conexion, err := cadenaConexion()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", conexion)
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
