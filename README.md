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

**Entrega actual: TP6** (rama `tp6`). La aplicación es un foro con interfaz web renderizada del lado del servidor (SSR) con [templ](https://templ.guide), sobre Go, PostgreSQL y sqlc. En esta entrega se incorpora [HTMX](https://htmx.org): crear, editar y borrar ya **no recargan la página**; el servidor devuelve solo el fragmento de HTML que cambió y HTMX lo reemplaza en el lugar correspondiente.

---

## 📑 Índice

- [Ejecución rápida](#-ejecución-rápida)
- [Requisitos](#requisitos)
- [Comandos disponibles](#comandos-disponibles)
- [Uso de la aplicación](#uso-de-la-aplicación)
  - [Pantalla general](#pantalla-general)
  - [Usuarios](#usuarios)
  - [Mensajes](#mensajes)
  - [Errores](#errores)
  - [Recorrido de ejemplo](#recorrido-de-ejemplo)
- [Introducción](#introducción)
  - [Dominio elegido](#dominio-elegido)
  - [Modelo de datos](#modelo-de-datos)
- [Arquitectura](#arquitectura)
- [Interfaz dinámica con HTMX (TP6)](#interfaz-dinámica-con-htmx-tp6)
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
git checkout tp6
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
- Puerto **8080** libre (la base de datos no se expone fuera de Docker).
- Conexión a internet en el navegador: HTMX y Pico CSS se cargan desde un CDN (el proyecto no incluye archivos `.js`).

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

## Uso de la aplicación

Con la aplicación levantada (`make up`), abrir **http://localhost:8080** en el navegador.

### Pantalla general

```text
┌──────────────────────────────────────────────────────────────────┐
│ Foro - Trabajo de Cursada                  Usuarios    Mensajes  │  ← barra de navegación
├──────────────────────────────────────────────────────────────────┤
│ [ recuadro rojo con el error, solo si algo falló ]               │
│                                                                  │
│ ┌─ Listado ─────────────────────────────┐  ┌─ Crear ──────────┐  │
│ │ ID │ ...datos... │ [Editar] [Borrar]  │  │ campos           │  │
│ │ ID │ ...datos... │ [Editar] [Borrar]  │  │ [Crear]          │  │
│ └───────────────────────────────────────┘  └──────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

- La **barra de navegación** permite pasar de una sección a otra. La página de inicio es **Usuarios**.
- Cada sección tiene el **listado** a la izquierda y el **formulario de alta** a la derecha. En pantallas angostas (celulares) el formulario aparece debajo del listado.
- Después de crear, editar o borrar, **la página no se recarga**: solo se actualiza el listado (o desaparece la fila borrada).
- Los datos se guardan en la base mientras la aplicación está levantada. Al detenerla con `Ctrl+C` se borran, y la próxima ejecución arranca vacía.

### Usuarios

**Crear un usuario**

1. En el panel **Crear Usuario**, completar **Nombre**, **Email** y **Contraseña**. Los tres campos son obligatorios y el email tiene que tener formato válido.
2. Presionar **Crear**. El usuario aparece en el listado con su **ID**, que se usa después para crear mensajes.

No puede haber dos usuarios con el mismo email.

**Editar un usuario**

1. En la fila del usuario, presionar **Editar**. Se despliega un formulario con los datos actuales.
2. Modificar **nombre** y/o **mail**. El campo **Nueva contraseña** es opcional: si se deja vacío, la contraseña no cambia.
3. Presionar **Guardar**. El listado se actualiza sin recargar la página. Para cerrar el formulario sin guardar, presionar **Editar** de nuevo.

**Borrar un usuario**

1. Presionar **Borrar** en la fila del usuario.
2. El navegador pide confirmación (*"Esta seguro que quiere eliminar?"*). Al aceptar, la fila desaparece de la tabla.

- Los mensajes que escribió ese usuario **no se borran**: quedan sin autor (la columna **Autor** aparece vacía al volver a cargar la sección Mensajes).

### Mensajes

**Crear un mensaje**

1. Ir a **Mensajes** en la barra de navegación.
2. En el panel **Crear Mensaje**, completar:
   - **ID usuario**: el **ID** del autor, tal como figura en el listado de Usuarios. Obligatorio.
   - **ID mensaje respuesta**: el **ID** del mensaje al que se responde. Se deja **vacío** si el mensaje no responde a otro.
   - **Texto**: el contenido del mensaje. Obligatorio.
3. Presionar **Crear**.

**Leer el listado**

| Columna | Significado |
| :--- | :--- |
| **ID** | Identificador del mensaje. Se usa para responderlo. |
| **Autor** | ID del usuario que lo escribió. Vacío si ese usuario fue borrado. |
| **Respuesta a** | ID del mensaje al que responde. Vacío si es un mensaje inicial. |
| **Texto** | Contenido del mensaje. |
| **Me gusta** | Cantidad de "me gusta". |

Para reconstruir una conversación, seguir la columna **Respuesta a**: un mensaje con `Respuesta a = 1` es una respuesta al mensaje con `ID = 1`.

**Editar un mensaje**

1. En la fila del mensaje, presionar **Editar**. Se despliega un formulario con los datos actuales, en este orden: **Usuario** (ID del autor), **Respuesta a** (ID), **Texto** y **Me gusta**.
2. Modificar los campos necesarios:
   - Dejar **Respuesta a** vacío convierte el mensaje en un mensaje inicial.
   - Dejar **Usuario** vacío deja el mensaje sin autor.
   - **Me gusta** se edita escribiendo el número; no puede ser negativo.
3. Presionar **Guardar**.

**Borrar un mensaje**

1. Presionar **Borrar** en la fila del mensaje y aceptar la confirmación.
2. La fila desaparece de la tabla.

- En la base se borran también **todas sus respuestas**, y las respuestas de esas respuestas (borrado en cascada). Como solo se quita la fila presionada, las respuestas dejan de verse en el listado al volver a cargar la sección **Mensajes**.

### Errores

Si una acción no se puede realizar, aparece un **recuadro rojo** arriba del contenido con el motivo, **los datos no se guardan** y el listado queda como estaba. La próxima acción exitosa borra el recuadro.

Casos habituales:

| Situación | Qué pasa |
| :--- | :--- |
| Campo obligatorio vacío | El navegador no deja enviar el formulario. Si igual llega vacío, se muestra el error de validación (por ejemplo, *"el texto del mensaje es obligatorio"*). |
| Email repetido al crear o editar un usuario | Se muestra un error y el usuario no se guarda. |
| ID de usuario o de mensaje que no existe | Se muestra *"el usuario o el mensaje respondido no existen"* y el mensaje no se guarda. |
| ID no numérico | Se muestra *"ID de usuario invalido"* (o el campo correspondiente). |
| "Me gusta" negativo | Se muestra *"la cantidad de me gusta no puede ser negativa"*. |

Como ya no hay recargas de página, después de un error se puede corregir el formulario y volver a enviarlo directamente.

### Recorrido de ejemplo

1. En **Usuarios**, crear a *Ana* (`ana@mail.com`) y a *Bruno* (`bruno@mail.com`). Quedan con los ID 1 y 2. Cada alta aparece en la tabla sin recargar la página.
2. En **Mensajes**, crear un mensaje con **ID usuario** `1`, **ID mensaje respuesta** vacío y el texto *"¿Alguien usó templ?"*. Queda con ID 1.
3. Crear otro con **ID usuario** `2`, **ID mensaje respuesta** `1` y el texto *"Sí, en el TP5"*. En el listado aparece con **Respuesta a** = 1.
4. Intentar crear un mensaje con **ID usuario** `99`: aparece el recuadro rojo y no se guarda nada.
5. Presionar **Editar** en el mensaje 2, cambiar **Me gusta** a `3` y **Guardar**.
6. Volver a **Usuarios** y borrar a *Bruno* (aceptando la confirmación). En **Mensajes**, el mensaje 2 sigue existiendo, pero con **Autor** vacío.
7. Borrar el mensaje 1. Al volver a cargar **Mensajes**, ya no aparece tampoco su respuesta, el mensaje 2.

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
Navegador (HTML + HTMX) ──HTTP──▶ webserver (Go, :8080) ──SQL──▶ db (PostgreSQL, :5432)
                                     │
                                     ├── handlers  → leen el formulario y validan (logic)
                                     ├── sqlc      → ejecutan las consultas
                                     └── views     → renderizan HTML con templ
                                                     (página completa o solo un fragmento)
```

La aplicación sigue un modelo de **renderizado del lado del servidor (SSR)**: el servidor arma el HTML con los datos de la base y lo envía al navegador.

- La **primera carga** de cada sección (`GET /usuarios`, `GET /mensajes`) devuelve la página completa (layout + listado + formulario).
- Las **acciones** (crear, editar, borrar) las envía HTMX por AJAX y el servidor responde **solo con el fragmento de HTML que cambió**. HTMX lo inserta en la página sin recargarla.

No se escribe JavaScript propio: todo el comportamiento dinámico se declara con atributos `hx-*` en los templates. HTMX se carga desde un CDN en el `<head>` del layout.

Los estilos se toman de [Pico CSS](https://picocss.com) (desde CDN) más algunas clases propias definidas en el layout.

---

## Interfaz dinámica con HTMX (TP6)

### Inclusión de la librería

En `views/layout.templ` se agrega, dentro del `<head>`:

```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.11/dist/htmx.min.js"
        integrity="sha384-..." crossorigin="anonymous"></script>
```

Se usa la versión 2 de HTMX con verificación de integridad (SRI).

### Atributos usados

| Dónde | Atributos | Efecto |
| :--- | :--- | :--- |
| Formulario de alta (`usuario_form.templ`, `mensaje_form.templ`) | `hx-post="/usuarios"` (o `/mensajes`), `hx-target="#tabla-usuarios"` (o `#tabla-mensajes`), `hx-swap="outerHTML"` | Envía el formulario por AJAX y reemplaza **todo** el contenedor del listado por el que devuelve el servidor. |
| Formulario de edición en cada fila (`*_list.templ`) | `hx-put="/usuarios/{id}"` (o `/mensajes/{id}`), mismo `hx-target` y `hx-swap="outerHTML"` | Envía la edición con el método `PUT` y reemplaza el listado actualizado. |
| Botón **Borrar** en cada fila (`*_list.templ`) | `hx-delete="/usuarios/{id}"` (o `/mensajes/{id}`), `hx-target="closest tr"`, `hx-swap="outerHTML"`, `hx-confirm="..."` | Pide confirmación, envía `DELETE` y reemplaza la fila (`<tr>`) por la respuesta vacía del servidor, es decir, la elimina. |

El contenedor del listado (`<section id="tabla-usuarios">` / `<section id="tabla-mensajes">`) existe siempre, aunque la lista esté vacía, para que HTMX siempre tenga un destino donde insertar la respuesta.

### Flujo de una alta

```text
1. El usuario presiona "Crear".
2. HTMX envía  POST /usuarios  (nombre, mail, contrasenia) por AJAX.
3. El handler valida, inserta con sqlc y vuelve a consultar ListUsers.
4. Responde 200 con dos fragmentos:
     <div id="alerta" hx-swap-oob="true"></div>      ← limpia la alerta
     <section id="tabla-usuarios">...</section>      ← listado actualizado
5. HTMX reemplaza #tabla-usuarios (hx-target + outerHTML) y,
   por "out-of-band", también reemplaza #alerta.
```

### Alerta de errores con Out-of-Band swap

El layout tiene un `<div id="alerta">` fijo arriba del contenido. El componente `Alerta(err)` renderiza ese mismo `div` con el atributo `hx-swap-oob="true"`, así que HTMX lo reemplaza **aunque no sea el `hx-target`** de la petición:

- **Acción exitosa:** se envía `Alerta("")` (vacía) junto al fragmento principal, para borrar un error anterior. El CSS `#alerta:empty { display: none; }` la oculta.
- **Error:** el handler responde `200` con el header `HX-Reswap: none` y solo `Alerta("<mensaje>")`. `HX-Reswap: none` le indica a HTMX que **no toque el target** (el listado y la fila quedan como estaban) y solo se procesa el swap out-of-band de la alerta.

Se responde `200` y no `4xx` en los errores esperados porque HTMX, por defecto, no procesa el cuerpo de respuestas con código de error.

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
│   │   ├── usuarios_web.go       # Handlers web de usuarios (página, alta, edición, borrado)
│   │   └── mensajes_web.go       # Handlers web de mensajes (página, alta, edición, borrado)
│   └── <strong>logic</strong>
│       ├── usuarios.go           # Validaciones de negocio de usuarios
│       └── mensajes.go           # Validaciones de negocio de mensajes
├── <strong>views</strong>                     # Templates templ (los *_templ.go se generan)
│   ├── layout.templ              # Layout común + componente Alerta
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

- **`cmd/server`**: archivo principal del servidor; abre la conexión y asocia cada ruta y método HTTP con su handler.
- **`internal/db`**: todo lo asociado a la base de datos: conexión, esquema, consultas, código generado y tests de persistencia.
- **`internal/handlers`**: reciben cada petición HTTP, leen el formulario o la URL, llaman a la validación y a sqlc, y responden con la página completa o con un fragmento de HTML para HTMX.
- **`internal/logic`**: reglas de negocio, como la verificación de campos vacíos o de valores negativos.
- **`views`**: templates de templ que generan el HTML.

Los archivos generados (`internal/db/sqlc/*.go` y `views/*_templ.go`) están en `.gitignore` y se crean con `make generate`, que `make up` ejecuta automáticamente.

---

## Rutas de la aplicación

Los `POST` y `PUT` reciben los datos como formulario (`application/x-www-form-urlencoded`), enviados por HTMX. Cualquier otra ruta devuelve `404`, y un método no soportado en una ruta existente devuelve `405`.

### Usuarios

| Método | Ruta | Handler | Qué hace | Respuesta |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/` | `UsuariosPageHandler` | Página de inicio: muestra la gestión de usuarios. | Página completa `UsuariosPage` |
| `GET` | `/usuarios` | `UsuariosPageHandler` | Lista los usuarios y muestra el formulario de alta. | Página completa `UsuariosPage` |
| `POST` | `/usuarios` | `CrearUsuarioWebHandler` | Crea un usuario (`nombre`, `mail`, `contrasenia`). | Fragmento `UsuarioList` + `Alerta` vacía (OOB) |
| `PUT` | `/usuarios/{id}` | `EditarUsuarioWebHandler` | Edita nombre y mail. Si la contraseña viene vacía, no se modifica. | Fragmento `UsuarioList` + `Alerta` vacía (OOB) |
| `DELETE` | `/usuarios/{id}` | `EliminarUsuarioWebHandler` | Elimina el usuario. Sus mensajes quedan sin autor. | `200` sin contenido principal (la fila desaparece) + `Alerta` vacía (OOB) |

### Mensajes

| Método | Ruta | Handler | Qué hace | Respuesta |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/mensajes` | `MensajesPageHandler` | Lista los mensajes y muestra el formulario de alta. | Página completa `MensajesPage` |
| `POST` | `/mensajes` | `CrearMensajeWebHandler` | Crea un mensaje (`usuario`, `texto` y opcionalmente `respuesta`). | Fragmento `MensajeList` + `Alerta` vacía (OOB) |
| `PUT` | `/mensajes/{id}` | `EditarMensajeWebHandler` | Edita autor, mensaje al que responde, texto y me gusta. Un campo vacío se guarda como `NULL`. | Fragmento `MensajeList` + `Alerta` vacía (OOB) |
| `DELETE` | `/mensajes/{id}` | `EliminarMensajeWebHandler` | Elimina el mensaje y, en cascada, sus respuestas. | `200` sin contenido principal (la fila desaparece) + `Alerta` vacía (OOB) |

Ante un error de validación o de base de datos, `POST`, `PUT` y `DELETE` responden `200` con el header `HX-Reswap: none` y solo el fragmento `Alerta` con el mensaje (ver [Manejo de errores](#manejo-de-errores)).

En `DELETE` el `{id}` se toma de la URL. En `PUT` el formulario de edición envía además el ID en un campo oculto (`id_usuario` o `id_mensaje`), que es el que usa el handler.

---

## Vistas (templates templ)

Las vistas están en la carpeta `views/`. Cada archivo `.templ` define un componente que `make generate` convierte en una función Go (por ejemplo, `UsuariosPage(...)`). Las páginas se arman por composición: una **página** usa el **layout** y dentro coloca una **lista** y un **formulario**. Las listas también se renderizan **solas**, como fragmento, en las respuestas a HTMX.

```text
UsuariosPage                         MensajesPage
└── Layout                           └── Layout
    ├── UsuarioList                      ├── MensajeList
    └── UsuarioForm                      └── MensajeForm
```

| Template | Componente | ¿Para qué sirve? | ¿Qué rutas satisface? |
| :--- | :--- | :--- | :--- |
| `layout.templ` | `Layout(title)` | Estructura HTML común: `<head>` con Pico CSS, estilos propios y el script de **HTMX**; barra de navegación (Usuarios / Mensajes) y el contenedor `<div id="alerta">`. El contenido de cada página se inserta con `children`. | Todas las páginas, a través de `UsuariosPage` y `MensajesPage`. |
| `layout.templ` | `Alerta(err)` | Renderiza `<div id="alerta" hx-swap-oob="true">` con el mensaje de error, o vacío. HTMX lo usa para actualizar la alerta "fuera de banda". | Todas las respuestas de `POST`, `PUT` y `DELETE`. |
| `usuarios_page.templ` | `UsuariosPage(usuarios)` | Página completa de gestión de usuarios: título y grilla con el listado y el formulario de alta, dentro del layout. | `GET /` y `GET /usuarios`. |
| `usuario_list.templ` | `UsuarioList(usuarios)` | `<section id="tabla-usuarios">` con la tabla de usuarios (ID, nombre, mail). Cada fila tiene un formulario desplegable para **editar** (`hx-put`) y un botón para **borrar** (`hx-delete`). | Dentro de `GET /` y `GET /usuarios`; sola como respuesta de `POST /usuarios` y `PUT /usuarios/{id}`. |
| `usuario_form.templ` | `UsuarioForm()` | Formulario de alta de usuario (nombre, email, contraseña) con `hx-post`. | Se muestra en `GET /` y `GET /usuarios`. Envía a `POST /usuarios`. |
| `mensajes_page.templ` | `MensajesPage(mensajes)` | Página completa de gestión de mensajes, dentro del layout. | `GET /mensajes`. |
| `mensaje_list.templ` | `MensajeList(mensajes)` | `<section id="tabla-mensajes">` con la tabla de mensajes (ID, autor, respuesta a, texto, me gusta). Los valores `NULL` se muestran como celda vacía. Cada fila tiene un formulario desplegable para **editar** y un botón para **borrar**. | Dentro de `GET /mensajes`; sola como respuesta de `POST /mensajes` y `PUT /mensajes/{id}`. |
| `mensaje_form.templ` | `MensajeForm()` | Formulario de alta de mensaje con `hx-post`: ID del autor, ID del mensaje al que responde (opcional) y texto. | Se muestra en `GET /mensajes`. Envía a `POST /mensajes`. |

---

## Manejo de errores

- **Validación de negocio** (`internal/logic`): campos obligatorios vacíos, `me_gusta` negativo, identificadores menores o iguales a 0, etc. El mensaje de la validación se muestra al usuario.
- **Validación de formato** (handlers): IDs que no son números.
- **Integridad de los datos**: las restricciones de la base impiden guardar un mensaje con un usuario o un mensaje respondido inexistentes (claves foráneas) y un usuario con un email repetido (`UNIQUE`).
- **Cómo se informa**: el handler responde `200` con el header `HX-Reswap: none` y el componente `Alerta(mensaje)`, que HTMX inserta en `#alerta` por *out-of-band swap* sin modificar el listado ni la fila. Ver [Alerta de errores con Out-of-Band swap](#alerta-de-errores-con-out-of-band-swap).
- **Errores inesperados** (por ejemplo, que falle la consulta del listado): se responde `500 Internal Server Error`.
- **Arranque**: si el servidor no puede conectarse a la base, termina con código de salida `1`.

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

Con la aplicación levantada (`make up`), en **http://localhost:8080**. El paso a paso detallado está en [Recorrido de ejemplo](#recorrido-de-ejemplo).

1. En **Usuarios**, crear uno o más usuarios. Verificar que aparecen en la tabla **sin que la página se recargue**.
2. En **Mensajes**, crear un mensaje indicando el ID de un usuario existente.
3. Crear otro mensaje completando "ID mensaje respuesta" con el ID del anterior.
4. Usar **Editar** en cada tabla: el listado se actualiza en el lugar.
5. Usar **Borrar**: aparece la confirmación y, al aceptar, la fila desaparece.
6. Probar casos inválidos, como crear un mensaje con un ID de usuario inexistente. Aparece el recuadro rojo y el listado no cambia.

Para comprobar que no hay recargas completas, abrir las herramientas de desarrollo del navegador (pestaña **Red / Network**): cada acción genera una única petición `fetch/xhr` cuya respuesta es solo un fragmento de HTML.
