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

-- name: UpdateUserPartial :exec
UPDATE Usuario SET nombre = $2, mail = $3 WHERE id_usuario = $1;

-- name: DeleteUser :exec
DELETE FROM Usuario WHERE id_usuario = $1;



-- name: CreateMensaje :one
INSERT INTO Mensaje (id_usuario, id_respuesta_a_mensaje, texto)
VALUES ($1, $2, $3)
RETURNING id_mensaje, id_usuario, id_respuesta_a_mensaje, texto, me_gusta;

-- name: GetMensajeById :one
SELECT id_mensaje, id_usuario, id_respuesta_a_mensaje, texto, me_gusta
FROM Mensaje
WHERE id_mensaje = $1;

-- name: ListMensajes :many
SELECT id_mensaje, id_usuario, id_respuesta_a_mensaje, texto, me_gusta
FROM Mensaje
ORDER BY id_mensaje;

-- name: UpdateMensaje :exec
UPDATE Mensaje 
SET texto = $2, me_gusta = $3 
WHERE id_mensaje = $1;

-- name: DeleteMensaje :exec
DELETE FROM Mensaje 
WHERE id_mensaje = $1;