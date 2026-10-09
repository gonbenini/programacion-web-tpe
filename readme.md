# 🌐 Trabajo Práctico Especial — Programación Web

<p align="center">
  <strong>Ingeniería de Sistemas · UNICEN</strong><br>
  Trabajo Práctico Especial de la cursada de Programación Web<br>
  Benini Gonzalo David<br>
  Lozano Francisco<br>
  Rezola Tomas Ezequiel<br>
</p>

Este repositorio corresponde al Trabajo Práctico Especial de la cursada de Programación Web de la carrera de Ingeniería de Sistemas de la UNICEN.

El objetivo es desarrollar una aplicación web a partir de un dominio simple, aplicando los conceptos y tecnologías vistos durante la cursada.


## 📑 Índice

- [Introducción](#introducción)
  - [Dominio elegido](#dominio-elegido)
  - [Ejemplo de la estructura](#ejemplo-de-la-estructura)
  - [Modelo de datos](#modelo-de-datos)
- [Persistencia y Generación de Código con sqlc](#persistencia-y-generación-de-código-con-sqlc)
  - [Estructura del código generado](#estructura-del-código-generado)
  - [Consultas SQL mapeadas](#consultas-sql-mapeadas)
- [Funcionalidades Básicas del Sistema](#funcionalidades-básicas-del-sistema)
  - [1. Registro de usuarios (`/register`)](#1-registro-de-usuarios-register)
  - [2. Listado de usuarios (`/usuarios`)](#2-listado-de-usuarios-usuarios)
- [Instancia de la app](#instancia-de-la-app)
- [Requisitos](#requisitos)
- [Despliegue](#despliegue)

---

## Introducción

### Dominio elegido

El dominio de la aplicación será una página tipo foro, donde los usuarios podrán participar en diferentes tableros y discusiones mediante mensajes.

Los tableros son formas de organizar y categorizar las discusiones.

Cada tablero puede tener **cero, uno o múltiples subtableros**, permitiendo construir una estructura de categorías de profundidad variable.

### Ejemplo de la estructura:

```text
📁 Tablero de inicio
│
├── 📁 Subtablero
│   │
│   ├── 📁 Subtablero
│   │   │
│   │   ├── 💬 Discusión
│   │   │   ├── 📝 Mensaje
│   │   │   ├── 📝 Mensaje
│   │   │   └── 📝 Mensaje
│   │   │
│   │   └── 💬 Discusión
│   │       └── 📝 Mensaje
│   │
│   └── 📁 Subtablero
│       └── 💬 Discusión
│
└── 📁 Subtablero
    └── 📁 Subtablero
        └── 💬 Discusión
            └── 📝 Mensaje
```

### Modelo de datos

- **Usuario:** `id_usuario`, `nombre`, `fecha_creacion`, `contrasenia`, `mail`
- **Tablero:** `id_tablero`, `nombre`, `id_tablero_padre`, `descripcion`
- **Discusión:** `nombre`, `id_usuario_creador`, `id_tablero`, `id_primer_mensaje`
- **Mensaje:** `id_usuario`, `id_respuesta_a_mensaje`, `texto`, `me_gusta`, `id_discusion`

---

## Persistencia y Generación de Código con sqlc

### Estructura del código generado

El código auto-generado por `sqlc` se ubica dentro del paquete `db` (`ws/db/sqlc`) y se divide en tres archivos principales:

1. **`models.go`**: Mapea los modelos de la base de datos a estructuras nativas de Go (`Usuario`, `Tablero`, `Discusion`, `Mensaje`).
   - Maneja los campos nulos de SQL utilizando los tipos seguros estándar de Go (`sql.NullInt32`, `sql.NullString`, `sql.NullTime`).
   - Incluye etiquetas JSON para facilitar la serialización/deserialización HTTP.
2. **`db.go`**: Define la interfaz `DBTX` (compatible con `*sql.DB` y `*sql.Tx`) y la estructura `Queries`, permitiendo ejecutar consultas simples o transaccionales (`WithTx`).
3. **`queries.sql.go`**: Contiene la implementación en Go de las sentencias SQL preparadas.

### Consultas SQL mapeadas

A través de `sqlc`, se implementaron las siguientes operaciones sobre la entidad **Usuario**:

| Método Go | Consulta SQL Subyacente | Descripción |
| :--- | :--- | :--- |
| `CreateUser` | `INSERT INTO Usuario (nombre, mail, contrasenia) VALUES ($1, $2, $3) RETURNING id_usuario, nombre, mail, fecha_creacion` | Inserta un nuevo usuario y retorna sus datos generados. |
| `ListUsers` | `SELECT id_usuario, nombre, mail FROM Usuario ORDER BY id_usuario` | Devuelve el listado completo de usuarios ordenados por ID. |
| `GetUserById` | `SELECT nombre, mail FROM Usuario WHERE id_usuario = $1` | Obtiene un usuario específico por su ID. |
| `UpdateUser` | `UPDATE Usuario SET nombre = $2, mail = $3, contrasenia = $4 WHERE id_usuario = $1` | Actualiza la información de un usuario existente. |
| `DeleteUser` | `DELETE FROM Usuario WHERE id_usuario = $1` | Elimina un usuario por su ID. |

---

## Funcionalidades Básicas del Sistema - Interaccion con webservice

El servidor expone endpoints básicos para la interacción con la interfaz de usuario, conectando el frontend con la base de datos mediante la capa generada por `sqlc`.

### Conexión a la Base de Datos (`main.go`)

La función `abrirDB()` gestiona la conexión con la base de datos PostgreSQL utilizando la URL de conexión del contenedor Docker:
- Cadena de conexión: `postgres://dbuser:dbpw@db:5432/foro?sslmode=disable`
- Configura un pool de conexiones máximo mediante `db.SetMaxOpenConns(25)`.
- Verifica la disponibilidad mediante `db.Ping()`.

---

### Endpoints de la API REST

El sistema expone una API RESTful estructurada bajo el patrón de separación de **Colecciones** y **Recursos específicos**. Las rutas son procesadas por los controladores (`handlers`), quienes validan el cuerpo de las peticiones (JSON), interactúan con la capa de lógica y persistencia, y devuelven las respuestas correspondientes con sus códigos de estado HTTP.

#### 1. Gestión de Usuarios

| Método HTTP | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/usuarios` | Devuelve el listado completo de usuarios en formato JSON. |
| `POST` | `/api/usuarios` | Registra un nuevo usuario validando los campos obligatorios del payload. |
| `GET` | `/api/usuarios/{id}` | Obtiene los detalles de un usuario específico por su ID. |
| `PUT` | `/api/usuarios/{id}` | Actualiza la información de un usuario existente. |
| `DELETE` | `/api/usuarios/{id}` | Elimina un usuario del sistema (borrado en cascada/set null según corresponda). |

#### 2. Gestión de Mensajes

Se implementó el CRUD para los mensajes del foro, contemplando el manejo de relaciones (autor y respuesta a otros mensajes) y el tratamiento de valores nulos provenientes de la base de datos.

| Método HTTP | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/mensajes` | Devuelve el listado de todos los mensajes. |
| `POST` | `/api/mensajes` | Crea un nuevo mensaje. Valida la existencia del `id_usuario` y opcionalmente del mensaje al que responde. |
| `GET` | `/api/mensajes/{id}` | Obtiene un mensaje específico. |
| `PUT` | `/api/mensajes/{id}` | Actualización parcial (PATCH lógico). Permite modificar campos específicos (`texto`, `me_gusta`, `id_usuario`) sin sobreescribir ni borrar los demás, gracias al mapeo dinámico mediante punteros en Go. |
| `DELETE` | `/api/mensajes/{id}` | Elimina un mensaje por su ID. |

---

## Instancia de la app

Tenemos la arquitectura definida en el archivo `docker-compose.yml`.  
Para el despliegue usamos la herramienta `Makefile`, en esta definimos comandos para operar la app. Siendo operar tareas como mantener, gestionar o desplegar la app según los comandos que definimos.

---

## Requisitos

Se debe clonar el repositorio y contar con una instalación previa de **Docker**, con el plugin **docker-compose** que habilite al comando `docker compose`.  
Además se debe contar con los puertos `8080` y `5432` libres.
Se debe poder usar el comando `curl`

---

## Despliegue

Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make up` en la terminal.

---

## Testeo


### Testeo persistencia
Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make test` en la terminal.

### Testeo api
Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make test-api` en la terminal.