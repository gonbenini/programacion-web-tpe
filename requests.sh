#!/bin/bash

RANDOM_MAIL="usuario_$(date +%s)@test.com"


# Testeamos el POST
RESPUESTA_POST_1=$(curl -s \
  -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Gon",
    "mail": "'"$RANDOM_MAIL"'",
    "contrasenia": "123456"
  }' \
  -i \
  -w "\n")

echo "$RESPUESTA_POST_1"

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


USER1_ID=$(echo "$RESPUESTA_POST_1" | grep -o '"id_usuario":[0-9]*' | tr -dc '0-9') # ID del usuario recién creado


# Testeamos el GET usuarios/{id}
curl -X GET http://localhost:8080/api/usuarios/$USER1_ID -i -w "\n"

RANDOM_MAIL_2="usuario2_$(date +%s)@test.com"

RESPUESTA_POST_2=$(curl -s \
  -X POST http://localhost:8080/api/usuarios \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "tomig",
    "mail": "'"$RANDOM_MAIL_2"'",
    "contrasenia": "654321"
  }' \
  -i \
  -w "\n")
USER2_ID=$(echo "$RESPUESTA_POST_2" | grep -o '"id_usuario":[0-9]*' | tr -dc '0-9')
  

# Testeamoso el GET usuarios/ (trae todos)
curl -X GET http://localhost:8080/api/usuarios -i -w "\n"


# Testeamos el eliminar de los dos registros
curl -X DELETE http://localhost:8080/api/usuarios/$USER1_ID -i -w "\n"

curl -X DELETE http://localhost:8080/api/usuarios/$USER2_ID -i -w "\n"