-- name: CreateUser :one
INSERT INTO Usuario (nombre, mail, contrasenia)
VALUES ($1, $2, $3)
RETURNING id_usuario, nombre, mail, fecha_creacion;

-- name: GetUserById :one
SELECT nombre, mail
FROM Usuario
WHERE id_usuario = $1;

-- name: ListUsers :many
SELECT id_usuario, nombre, mail
FROM Usuario
ORDER BY id_usuario;

-- name: UpdateUser :exec
UPDATE Usuario SET nombre = $2, mail = $3, contrasenia = $4 WHERE id_usuario = $1;

-- name: DeleteUser :exec
DELETE FROM Usuario WHERE id_usuario = $1;