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
- [Estructura de proyecto](#estructura-de-proyecto)
- [Persistencia y Generación de Código con sqlc](#persistencia-y-generación-de-código-con-sqlc)
  - [Estructura del código generado](#estructura-del-código-generado)
  - [Consultas SQL mapeadas](#consultas-sql-mapeadas)
- [Endpoints de la API REST](#endpoints-de-la-API-REST)
- [Instancia de la app](#instancia-de-la-app)
- [Requisitos](#requisitos)
- [Despliegue](#despliegue)

---

## Introducción

### Dominio elegido

El dominio de la aplicación será una página tipo foro, donde los usuarios podrán participar escribiendo mensajes en respuesta a otros mensajes o por su cuenta.

Cada mensaje tiene asociado un usuario y puede llegar a tener un mensaje al que responde.

### Ejemplo de la estructura:

```text
💬 Mensaje inicial (👤 Usuario A)
│
├── 📝 Respuesta (👤 Usuario B)
│   │
│   ├── 📝 Respuesta (👤 Usuario C)
│   │   │
│   │   ├── 📝 Respuesta (👤 Usuario A)
│   │   │
│   │   └── 📝 Otra respuesta (👤 Usuario D)
│   │
│   └── 📝 Otra respuesta (👤 Usuario C)
│
├── 📝 Otra respuesta (👤 Usuario D)
│   │
│   └── 📝 Respuesta (👤 Usuario B)
│
└── 📝 Otra respuesta (👤 Usuario C)
```

### Modelo de datos

- **Usuario:** `id_usuario`, `nombre`, `fecha_creacion`, `contrasenia`, `mail`
- **Mensaje:** `id_usuario`, `id_respuesta_a_mensaje`, `texto`, `me_gusta`

---

## Estructura de proyecto
<pre>
.
├── <strong><span style="color: orange">cmd</span></strong>
│   └── <strong><span style="color: orange">server</span></strong>
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── <strong><span style="color: orange">internal</span></strong>
│   ├── <strong><span style="color: orange">db</span></strong>
│   │   ├── db_controller.go
│   │   ├── <strong><span style="color: orange">queries</span></strong>
│   │   │   └── queries.sql
│   │   ├── <strong><span style="color: orange">schema</span></strong>
│   │   │   └── schema.sql
│   │   ├── <strong><span style="color: orange">sqlc</span></strong>
│   │   │   ├── db.go
│   │   │   ├── models.go
│   │   │   └── queries.sql.go
│   │   └── <strong><span style="color: orange">tests</span></strong>
│   │       ├── db_mensaje_test.go
│   │       └── db_usuarios_test.go
│   ├── <strong><span style="color: orange">handlers</span></strong>
│   │   ├── mensajes.go
│   │   └── usuarios.go
│   └── <strong><span style="color: orange">logic</span></strong>
│       ├── mensajes.go
│       └── usuarios.go
├── Makefile
├── readme.md
├── requests.sh
└── sqlc.yaml
</pre>

<strong>Reestructuramos el arbol de archivos para:
- Contener en <strong><span style="color: orange">cmd/server</span></strong> el archivo principal del servidor que contiene los endpoints de la API para cada entidad.

- Contener en <strong><span style="color: orange">internal/db</span></strong> toda la logica asociada al manejo de la base de datos incluyendo el controlador y sus respectivos tests de persistencia.

- Contener en <strong><span style="color: orange">internal/handlers</span></strong> toda la logica asociada al manejo de consultas a la API para las entidades usuarios y mensajes por separado.

- Contener en <strong><span style="color: orange">internal/logic</span></strong> toda la logica asociada a los chequeos de negocio como la verificacion de campos vacios.

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


### Conexión a la Base de Datos (`main.go`)

La función `abrirDB()` gestiona la conexión con la base de datos PostgreSQL utilizando la URL de conexión del contenedor Docker:
- Cadena de conexión: `postgres://dbuser:dbpw@db:5432/foro?sslmode=disable`
- Configura un pool de conexiones máximo mediante `db.SetMaxOpenConns(25)`.
- Verifica la disponibilidad mediante `db.Ping()`.

---

### Endpoints de la API REST

El sistema expone una API RESTful

#### 1. Gestión de Usuarios

| Método HTTP | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/usuarios` | Devuelve el listado completo de usuarios en formato JSON. |
| `POST` | `/api/usuarios` | Registra un nuevo usuario validando los campos obligatorios del payload. |
| `GET` | `/api/usuarios/{id}` | Obtiene los detalles de un usuario específico por su ID. |
| `PUT` | `/api/usuarios/{id}` | Actualiza la información de un usuario existente. |
| `DELETE` | `/api/usuarios/{id}` | Elimina un usuario. |

#### 2. Gestión de Mensajes

| Método HTTP | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/mensajes` | Devuelve el listado de todos los mensajes. |
| `POST` | `/api/mensajes` | Crea un nuevo mensaje. Valida la existencia del `id_usuario` y opcionalmente del mensaje al que responde. |
| `GET` | `/api/mensajes/{id}` | Obtiene un mensaje específico. |
| `PUT` | `/api/mensajes/{id}` | Actualización de los mensajes. |
| `DELETE` | `/api/mensajes/{id}` | Elimina un mensaje por su ID. |

La intencion de esto es dividir aquellas consultas que no necesitan de, en este caso, un identificador como `create` o `getAll` de las que si lo necesitan como `get`, `update` y `delete`.

Asi, dada la ruta especificada, construimos un `handler` para cada tipo de direccion para cada entidad. Por ejemplo, para usuarios, tenemos `UsuariosAPIHandler` asociado a `api/usuarios` que se encarga de redirigir a la operacion de creacion u obtencion total segun el metodo `http` provisto (POST o GET). Luego, tenemos `UsuarioAPIHandler` asociado a `api/usuarios/` que se encarga de extraer el id de la URL y redirigir a la operacion de actualizacion, obtencion por id o eliminacion segun el metodo `http` provisto (GET, UPDATE o DELETE).

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