#!/bin/bash

RANDOM_MAIL="usuario_$(date +%s)@test.com"

curl \
  -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Gon",
    "mail": "'"$RANDOM_MAIL"'",
    "contrasenia": "123456"
  }' \
  -i \
  -w "\n"

# -X POST: 
#    Define explícitamente el método HTTP de la petición como POST 
#
# -H "Content-Type: application/json": 
#    Cabecera HTTP (Header), le explica al servidor con que formato esta estructurado el cuerpo de la peticion.
#
# -d '{...}': 
#    Cuerpo (body) de la petición que contiene los datos del objeto que se envia en formato JSON.
#
# -i: 
#    Le dice al curl que incluya las cabeceras (headers) de la respuesta HTTP en la salida estándar. Lo usamos para ver el estado, que si fue exitoso es 201 Created.
#
# -w "\n": 
#    Define un formato de salida posterior a la ejecución. Ahora no tiene sentido pero a medida que agreguemos endpoints al test ayuda a la claridad.