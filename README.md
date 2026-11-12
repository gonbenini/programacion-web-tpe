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

**Entrega actual: TP5** (rama `tp5`). La aplicación es un foro con interfaz web renderizada del lado del servidor (SSR) con [templ](https://templ.guide), sobre Go, PostgreSQL y sqlc.

---

## 📑 Índice

- [Ejecución rápida](#-ejecución-rápida)
- [Requisitos](#requisitos)
- [Comandos disponibles](#comandos-disponibles)
- [Introducción](#introducción)
  - [Dominio elegido](#dominio-elegido)
  - [Modelo de datos](#modelo-de-datos)
- [Arquitectura](#arquitectura)
- [Estructura de proyecto](#estructura-de-proyecto)
- [Rutas de la aplicación](#rutas-de-la-aplicación)
- [Vistas (templates templ)](#vistas-templates-templ)
- [Manejo de errores](#manejo-de-errores)
- [Persistencia y generación de código con sqlc](#persistencia-y-generación-de-código-con-sqlc)
- [Testeo](#testeo)

---

## 🚀 Ejecución rápida

```bash
git clone <url-del-repositorio>
cd <carpeta-del-repositorio>
git checkout tp5
make up
```

`make up` realiza todas las tareas necesarias: genera el código, borra contenedores y volúmenes anteriores, compila, levanta la base de datos, espera a que esté lista y ejecuta la aplicación.

Cuando la consola muestre `Servidor con formulario escuchando en http://localhost:8080`, abrir **http://localhost:8080** en el navegador.

Para terminar, presionar **`Ctrl+C`**. Al detenerse, `make up` borra automáticamente los contenedores y volúmenes. No hace falta ejecutar ningún otro comando.

> Al cortar con `Ctrl+C`, `make` muestra `make: *** [Makefile:...: up] Error 130`. Es normal: indica que el comando se interrumpió con `Ctrl+C`, y la limpieza ya se realizó.

> La primera ejecución puede tardar más porque Docker descarga las imágenes (`postgres`, `golang`, `sqlc`, `templ`) y Go descarga las dependencias.

> Cada ejecución de `make up` arranca con la base de datos vacía: los datos cargados no se conservan entre ejecuciones.

---

## Requisitos

- **Docker** con el plugin **docker compose** (comando `docker compose`).
- **make**.
- Puertos **8080** (aplicación) y **5432** (base de datos) libres.

No hace falta instalar Go, sqlc ni templ en la máquina: todo se ejecuta dentro de contenedores.

---

## Comandos disponibles

Todos los comandos se ejecutan desde la carpeta raíz del proyecto (la que contiene `Makefile` y `docker-compose.yml`).

| Comando | Qué hace |
| :--- | :--- |
| `make up` | **Ejecuta la aplicación con todas sus tareas previas y posteriores** (ver detalle abajo). Al cortar con `Ctrl+C` borra contenedores y volúmenes. |
| `make down` | Detiene los contenedores y **borra los volúmenes** (la base vuelve a quedar vacía). `make up` ya lo hace solo; sirve para limpiar a mano si algo quedó levantado. |
| `make generate` | Genera el código Go a partir de `schema.sql`/`queries.sql` (**sqlc**) y de los archivos `.templ` (**templ**). Ambas herramientas corren en contenedores efímeros. |
| `make test` | Genera el código, corre los tests de persistencia contra la base y al terminar borra contenedores y volúmenes. |
| `make clean` | Borra los archivos autogenerados (código de sqlc y de templ). |

### Qué ocurre al ejecutar `make up`

**Tareas previas**

1. **Generación de código:**
   - **sqlc** lee `internal/db/schema/schema.sql` y `internal/db/queries/queries.sql` y genera el paquete `internal/db/sqlc`.
   - **templ** compila cada `views/*.templ` a un archivo `views/*_templ.go`.
2. **Limpieza previa:** `docker compose down -v --remove-orphans` borra contenedores y volúmenes de una ejecución anterior. Así la base siempre se crea a partir del `schema.sql` actual, que PostgreSQL solo ejecuta cuando el volumen está vacío.
3. **Compilación:** `go build ./...` se ejecuta en un contenedor efímero de Go. Si hay un error de compilación, `make` se detiene y no levanta nada.
4. **Levantar la base:** se levanta solo el contenedor `db` (PostgreSQL), que al arrancar con el volumen vacío ejecuta `schema.sql` para crear las tablas.
5. **Esperar a la base:** se repite `pg_isready` cada segundo hasta que PostgreSQL acepta conexiones. Recién entonces se levanta el servidor, así nunca intenta conectarse a una base que todavía no está lista.

**Ejecución**

6. Se levanta el contenedor `webserver` (imagen `golang`), que ejecuta `cmd/server/main.go` con `go run`. El servidor queda escuchando en el puerto `8080` y sus logs se ven en la consola.

**Tareas posteriores**

7. Al presionar `Ctrl+C` (o si el servidor termina por cualquier motivo), `docker compose down -v` borra los contenedores y volúmenes.

---

## Introducción

### Dominio elegido

El dominio de la aplicación es una página tipo foro, donde los usuarios pueden participar escribiendo mensajes en respuesta a otros mensajes o por su cuenta.

Cada mensaje tiene asociado un usuario y puede tener un mensaje al que responde.

Ejemplo de la estructura:

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

- **Usuario:** `id_usuario`, `nombre`, `mail` (único), `contrasenia`, `fecha_creacion`
- **Mensaje:** `id_mensaje`, `id_usuario`, `id_respuesta_a_mensaje`, `texto`, `me_gusta`

Relaciones:

- `Mensaje.id_usuario` → `Usuario.id_usuario`. Si se borra el usuario, el mensaje se conserva sin autor (`ON DELETE SET NULL`).
- `Mensaje.id_respuesta_a_mensaje` → `Mensaje.id_mensaje`. Si se borra un mensaje, se borran también sus respuestas (`ON DELETE CASCADE`).

---

## Arquitectura

```text
Navegador ──HTTP──▶ webserver (Go, :8080) ──SQL──▶ db (PostgreSQL, :5432)
                       │
                       ├── handlers  → leen el formulario y validan (logic)
                       ├── sqlc      → ejecutan las consultas
                       └── views     → renderizan el HTML con templ
```

La aplicación sigue un modelo de **renderizado del lado del servidor (SSR)**: el servidor arma el HTML completo con los datos de la base y lo envía al navegador. No se usa JavaScript.

Los formularios siguen el patrón **PRG (Post / Redirect / Get)**: cada `POST` modifica la base y responde con una redirección `303 See Other` a la página de listado (`GET`). Así, recargar la página no reenvía el formulario.

Los estilos se toman de [Pico CSS](https://picocss.com) (desde CDN) más algunas clases propias definidas en el layout.

---

## Estructura de proyecto

<pre>
.
├── <strong>cmd</strong>
│   └── <strong>server</strong>
│       └── main.go               # Punto de entrada: conexión a la BD y registro de rutas
├── <strong>internal</strong>
│   ├── <strong>db</strong>
│   │   ├── db_controller.go      # AbrirDB(): conexión a PostgreSQL
│   │   ├── <strong>queries</strong>
│   │   │   └── queries.sql       # Consultas SQL que lee sqlc
│   │   ├── <strong>schema</strong>
│   │   │   └── schema.sql        # Definición de tablas
│   │   ├── <strong>sqlc</strong>              # (generado por sqlc, no versionado)
│   │   │   ├── db.go
│   │   │   ├── models.go
│   │   │   └── queries.sql.go
│   │   └── <strong>tests</strong>
│   │       ├── db_mensaje_test.go
│   │       └── db_usuarios_test.go
│   ├── <strong>handlers</strong>
│   │   ├── usuarios_web.go       # Handlers de las páginas web (usuarios y mensajes)
│   │   ├── usuarios.go           # Handlers de la API JSON del TP3 (sin rutas registradas)
│   │   └── mensajes.go           # Handlers de la API JSON del TP3 (sin rutas registradas)
│   └── <strong>logic</strong>
│       ├── usuarios.go           # Validaciones de negocio de usuarios
│       └── mensajes.go           # Validaciones de negocio de mensajes
├── <strong>views</strong>                     # Templates templ (los *_templ.go se generan)
│   ├── layout.templ
│   ├── usuarios_page.templ
│   ├── usuario_list.templ
│   ├── usuario_form.templ
│   ├── mensajes_page.templ
│   ├── mensaje_list.templ
│   └── mensaje_form.templ
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── sqlc.yaml
</pre>

- **`cmd/server`**: archivo principal del servidor; abre la conexión y asocia cada ruta con su handler.
- **`internal/db`**: todo lo asociado a la base de datos: conexión, esquema, consultas, código generado y tests de persistencia.
- **`internal/handlers`**: reciben cada petición HTTP, leen el formulario, llaman a la validación y a sqlc, y responden con una vista o una redirección.
- **`internal/logic`**: reglas de negocio, como la verificación de campos vacíos o de valores negativos.
- **`views`**: templates de templ que generan el HTML.

Los archivos generados (`internal/db/sqlc/*.go` y `views/*_templ.go`) están en `.gitignore` y se crean con `make generate`, que `make up` ejecuta automáticamente.

---

## Rutas de la aplicación

Todas las rutas devuelven HTML o redirigen. Los `POST` reciben los datos como formulario (`application/x-www-form-urlencoded`).

### Usuarios

| Método | Ruta | Handler | Qué hace | Respuesta |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/` | `UsuariosPageHandler` | Página de inicio: muestra la gestión de usuarios. | Vista `UsuariosPage` |
| `GET` | `/usuarios` | `UsuariosPageHandler` | Lista los usuarios y muestra el formulario de alta. | Vista `UsuariosPage` |
| `POST` | `/usuarios` | `CrearUsuarioWebHandler` | Crea un usuario (`nombre`, `mail`, `contrasenia`). | `303` → `/usuarios` |
| `POST` | `/usuarios/editar` | `EditarUsuarioWebHandler` | Edita nombre y mail. Si la contraseña viene vacía, no se modifica. | `303` → `/usuarios` |
| `POST` | `/usuarios/eliminar` | `EliminarUsuarioWebHandler` | Elimina el usuario. Sus mensajes quedan sin autor. | `303` → `/usuarios` |

### Mensajes

| Método | Ruta | Handler | Qué hace | Respuesta |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/mensajes` | `MensajesPageHandler` | Lista los mensajes y muestra el formulario de alta. | Vista `MensajesPage` |
| `POST` | `/mensajes` | `CrearMensajeWebHandler` | Crea un mensaje (`usuario`, `texto` y opcionalmente `respuesta`). Verifica que el usuario exista. | `303` → `/mensajes` |
| `POST` | `/mensajes/editar` | `EditarMensajeWebHandler` | Edita autor, mensaje al que responde, texto y me gusta. Un campo vacío se guarda como `NULL`. | `303` → `/mensajes` |
| `POST` | `/mensajes/eliminar` | `EliminarMensajeWebHandler` | Elimina el mensaje y, en cascada, sus respuestas. | `303` → `/mensajes` |

Las rutas de edición y eliminación usan `POST` porque los formularios HTML solo admiten `GET` y `POST`. El identificador viaja en un campo oculto del formulario (`id_usuario` o `id_mensaje`).

---

## Vistas (templates templ)

Las vistas están en la carpeta `views/`. Cada archivo `.templ` define un componente que `make generate` convierte en una función Go (por ejemplo, `UsuariosPage(...)`). Las páginas se arman por composición: una **página** usa el **layout** y dentro coloca una **lista** y un **formulario**.

```text
UsuariosPage                         MensajesPage
└── Layout                           └── Layout
    ├── UsuarioList                      ├── MensajeList
    └── UsuarioForm                      └── MensajeForm
```

| Template | Componente | ¿Para qué sirve? | ¿Qué rutas satisface? |
| :--- | :--- | :--- | :--- |
| `layout.templ` | `Layout(title, err)` | Estructura HTML común a todas las páginas: `<head>` con Pico CSS y estilos propios, barra de navegación (Usuarios / Mensajes) y, si `err` no está vacío, un recuadro de alerta con el error. El contenido de cada página se inserta con `children`. | Todas las páginas, de forma indirecta, a través de `UsuariosPage` y `MensajesPage`. |
| `usuarios_page.templ` | `UsuariosPage(usuarios, err)` | Página completa de gestión de usuarios: título y grilla con el listado y el formulario de alta, dentro del layout. | `GET /` y `GET /usuarios`. |
| `usuario_list.templ` | `UsuarioList(usuarios)` | Tabla de usuarios (ID, nombre, mail). Cada fila tiene un formulario desplegable para **editar** y un botón para **borrar**. | Se muestra en `GET /` y `GET /usuarios`. Sus formularios envían a `POST /usuarios/editar` y `POST /usuarios/eliminar`. |
| `usuario_form.templ` | `UsuarioForm()` | Formulario de alta de usuario (nombre, email, contraseña). | Se muestra en `GET /` y `GET /usuarios`. Envía a `POST /usuarios`. |
| `mensajes_page.templ` | `MensajesPage(mensajes, err)` | Página completa de gestión de mensajes: título y grilla con el listado y el formulario de alta, dentro del layout. Recibe el texto de error para mostrarlo en la alerta. | `GET /mensajes`. También se renderiza con el error cuando falla `POST /mensajes`, `POST /mensajes/editar` o `POST /mensajes/eliminar`. |
| `mensaje_list.templ` | `MensajeList(mensajes)` | Tabla de mensajes (ID, autor, respuesta a, texto, me gusta). Los valores `NULL` se muestran como celda vacía. Cada fila tiene un formulario desplegable para **editar** y un botón para **borrar**. | Se muestra en `GET /mensajes`. Sus formularios envían a `POST /mensajes/editar` y `POST /mensajes/eliminar`. |
| `mensaje_form.templ` | `MensajeForm()` | Formulario de alta de mensaje: ID del autor, ID del mensaje al que responde (opcional) y texto. | Se muestra en `GET /mensajes`. Envía a `POST /mensajes`. |

---

## Manejo de errores

- **Validación de negocio** (`internal/logic`): campos obligatorios vacíos, `me_gusta` negativo, identificadores menores o iguales a 0, etc.
- **Existencia de datos relacionados**: antes de crear o editar un mensaje se verifica que el usuario indicado exista.
- **Cómo se informa**: en la sección de mensajes, los errores se muestran en la misma página, en un recuadro de alerta del layout (respuesta `400`). En la sección de usuarios, los errores se devuelven como texto plano.

---

## Persistencia y generación de código con sqlc

### Código generado

`sqlc` genera el paquete `db` en `internal/db/sqlc/` a partir del esquema y las consultas (ver `sqlc.yaml`):

1. **`models.go`**: estructuras Go para cada tabla (`Usuario`, `Mensaje`). Los campos que admiten `NULL` usan los tipos `sql.NullInt32`, `sql.NullTime`, etc. Incluye etiquetas JSON.
2. **`db.go`**: la interfaz `DBTX` (compatible con `*sql.DB` y `*sql.Tx`) y la estructura `Queries`, que permite ejecutar consultas simples o dentro de una transacción (`WithTx`).
3. **`queries.sql.go`**: implementación en Go de cada consulta de `queries.sql`.

### Consultas SQL

| Método Go | Entidad | Tipo | Descripción |
| :--- | :--- | :--- | :--- |
| `CreateUser` | Usuario | `:one` | Inserta un usuario y devuelve `id_usuario`, `nombre`, `mail` y `fecha_creacion`. |
| `GetUserById` | Usuario | `:one` | Obtiene nombre y mail de un usuario por ID. |
| `ListUsers` | Usuario | `:many` | Lista todos los usuarios ordenados por ID. |
| `UpdateUser` | Usuario | `:exec` | Actualiza nombre, mail y contraseña. |
| `UpdateUserPartial` | Usuario | `:exec` | Actualiza nombre y mail sin tocar la contraseña. |
| `DeleteUser` | Usuario | `:exec` | Elimina un usuario por ID. |
| `CreateMensaje` | Mensaje | `:one` | Inserta un mensaje y lo devuelve completo. |
| `GetMensajeById` | Mensaje | `:one` | Obtiene un mensaje por ID. |
| `ListMensajes` | Mensaje | `:many` | Lista todos los mensajes ordenados por ID. |
| `UpdateMensaje` | Mensaje | `:exec` | Actualiza autor, mensaje respondido, texto y me gusta. |
| `DeleteMensaje` | Mensaje | `:exec` | Elimina un mensaje por ID. |

### Conexión a la base de datos

`AbrirDB()` (en `internal/db/db_controller.go`) abre la conexión con PostgreSQL:

- Cadena de conexión: `postgres://dbuser:dbpw@db:5432/foro?sslmode=disable` (`db` es el nombre del servicio en `docker-compose.yml`).
- Verifica la disponibilidad con `db.Ping()`.
- Limita el pool a 25 conexiones con `db.SetMaxOpenConns(25)`.

---

## Testeo

### Tests de persistencia

```bash
make test
```

Genera el código, levanta la base y corre `go test -v ./internal/db/tests` en un contenedor efímero. Los tests crean, leen, actualizan y borran usuarios y mensajes directamente con sqlc. Al terminar, borra contenedores y volúmenes.

### Prueba manual de la interfaz

Con la aplicación levantada (`make up`), en **http://localhost:8080**:

1. En **Usuarios**, crear uno o más usuarios.
2. En **Mensajes**, crear un mensaje indicando el ID de un usuario existente.
3. Crear otro mensaje completando "Respuesta a mensaje" con el ID del anterior.
4. Usar **Editar** y **Borrar** en cada tabla. Tras cada acción, la página se recarga con los datos actualizados.
5. Probar casos inválidos, como crear un mensaje con un ID de usuario inexistente. Se muestra un error en lugar de guardar el dato.
