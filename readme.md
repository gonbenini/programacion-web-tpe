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

### 1. Registro de usuarios (`/register`)

Permite la creación de un nuevo usuario en la plataforma mediante una solicitud asíncrona desde el navegador.

#### Flujo de trabajo:
1. **Frontend (`index.html` & `index.js`)**:
   - El usuario ingresa `nombre`, `mail` y `contrasenia` en el formulario HTML.
   - JavaScript intercepta el evento de envío (`submit`) previniendo la recarga predeterminada de la página (`event.preventDefault()`).
   - Envía una petición `POST` mediante la **Fetch API** codificando los datos del formulario (`URLSearchParams`).
2. **Backend (`main.go` -> `handleRegister`)**:
   - Valida la ruta y procesa los datos ingresados con `r.ParseForm()`.
   - Llama al método `queries.CreateUser(ctx, params)` generado por `sqlc`.
   - Almacena el nuevo usuario en PostgreSQL.
   - Responde con un JSON indicando el estado del proceso (`{"mensaje": "Se envió correctamente"}`).
3. **Respuesta en interfaz**:
   - JavaScript recibe la respuesta JSON y muestra un mensaje de confirmación verde o de error rojo en el elemento `<p id="mensaje-alerta">`.

---

### 2. Listado de usuarios (`/usuarios`)

Permite la consulta de todos los usuarios registrados en el sistema.

#### Flujo de trabajo:
1. **Frontend (`index.html`)**:
   - Contiene un botón *"Ver usuarios"* que redirige al navegador directamente a `/usuarios`.
2. **Backend (`main.go` -> `handleUsuarios`)**:
   - El handler instancia la estructura `Queries` de `sqlc`.
   - Ejecuta `queries.ListUsers(ctx)` para traer todos los registros de la tabla `Usuario`.
   - Configura la cabecera `Content-Type: text/html; charset=utf-8`.
   - Construye y renderiza dinámicamente un documento HTML con una lista no ordenada (`<ul>`) mostrando `ID`, `Nombre` y `Mail` de cada usuario.

---

## Instancia de la app

Tenemos la arquitectura definida en el archivo `docker-compose.yml`.  
Para el despliegue usamos la herramienta `Makefile`, en esta definimos comandos para operar la app. Siendo operar tareas como mantener, gestionar o desplegar la app según los comandos que definimos.

---

## Requisitos

Se debe clonar el repositorio y contar con una instalación previa de **Docker**, con el plugin **docker-compose** que habilite al comando `docker compose`.  
Además se debe contar con los puertos `8080` y `5432` libres.

---

## Despliegue

Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make up` en la terminal.

---

## Testeo


### Testeo persistencia
Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make test` en la terminal.

### Testeo api
Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml` y `Makefile`), hay que ejecutar el comando `make up` en la terminal. Una vez que el contenedor este funcionando, correr el script `requests.sh`