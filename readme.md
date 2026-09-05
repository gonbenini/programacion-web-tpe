# 🌐 Trabajo Práctico Especial — Programación Web

<p align="center">
  <strong>Ingeniería de Sistemas · UNICEN</strong><br>
  Trabajo Práctico Especial de la cursada de Programación Web
</p>

Este repositorio corresponde al Trabajo Práctico Especial de la cursada de Programación Web de la carrera de Ingeniería de Sistemas de la UNICEN.

El objetivo es desarrollar una aplicación web a partir de un dominio simple, aplicando los conceptos y tecnologías vistos durante la cursada.


## 📑 Índice

- [Introducción](#introducción)
- [Instalación](#instalación)
- [Despliegue](#despliegue)

---

## Introducción



###  Dominio elegido

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

## Instancia de la app

Tenemos la arquitectura definida en el archivo `docker-compose.up`.  
Para el despliegue usamos la herramienta `Makefile`, en esta definimos comandos para operar la app. Siendo operar tareas como mantener, gestionar o desplegar la app segun los comandos que definimos.

---

## Requisitos
Se debe clonar el repositorio y contar con una instalación previa de **docker**, con el plugin **docker-compose** que habilite al comando `docker compose`.  
Ademas se debe contar con los puertos `8080` y `5432` libres.

---

## Despliegue

Dentro de la carpeta raíz (la que contiene el archivo `docker-compose.yml`), hay que ejecutar el comando `docker compose up` en la terminal.
