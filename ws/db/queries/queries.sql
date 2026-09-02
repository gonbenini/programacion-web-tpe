-- name: CreateUser :one
INSERT INTO Usuario (nombre, mail, contrasenia)
VALUES ($1, $2, $3)
RETURNING id_usuario, nombre, mail, fecha_creacion;

-- name: ListUsers :many
SELECT id_usuario, nombre, mail, fecha_creacion
FROM Usuario
ORDER BY id_usuario;