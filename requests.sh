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


test_edicion_mensajes() {
  local BASE_URL="http://localhost:8080/api"
  local TIMESTAMP=$(date +%s)

  # Crear usuarios de prueba
  local RES_U1=$(curl -s -X POST $BASE_URL/usuarios -H "Content-Type: application/json" -d '{"nombre": "Gon Edit","mail": "gon_edit_'$TIMESTAMP'@test.com","contrasenia": "123"}')
  local ID_U1=$(echo "$RES_U1" | grep -o '"id_usuario":[0-9]*' | tr -dc '0-9')

  local RES_U2=$(curl -s -X POST $BASE_URL/usuarios -H "Content-Type: application/json" -d '{"nombre": "Tomi Edit","mail": "tomi_edit_'$TIMESTAMP'@test.com","contrasenia": "123"}')
  local ID_U2=$(echo "$RES_U2" | grep -o '"id_usuario":[0-9]*' | tr -dc '0-9')

  # Crear mensajes de prueba
  local RES_M1=$(curl -s -X POST $BASE_URL/mensajes -H "Content-Type: application/json" -d '{"id_usuario": '$ID_U1', "texto": "Mensaje original."}')
  local ID_M1=$(echo "$RES_M1" | grep -o '"id_mensaje":[0-9]*' | tr -dc '0-9')

  local RES_M2=$(curl -s -X POST $BASE_URL/mensajes -H "Content-Type: application/json" -d '{"id_usuario": '$ID_U2', "texto": "Mensaje que sera respondido."}')
  local ID_M2=$(echo "$RES_M2" | grep -o '"id_mensaje":[0-9]*' | tr -dc '0-9')

  # Testear PUT: Solo Texto
  echo -e "\nResultado de editar solo el texto (valido):"
  curl -s -X PUT $BASE_URL/mensajes/$ID_M1 -H "Content-Type: application/json" -d '{"texto": "Texto editado por PUT"}' -i -w "\n"

  # Testear PUT: Solo Me Gusta
  echo -e "\nResultado de editar solo los me gusta (valido):"
  curl -s -X PUT $BASE_URL/mensajes/$ID_M1 -H "Content-Type: application/json" -d '{"me_gusta": 99}' -i -w "\n"

  # Testear PUT: Solo id_usuario
  echo -e "\nResultado de cambiar solo el autor (valido):"
  curl -s -X PUT $BASE_URL/mensajes/$ID_M1 -H "Content-Type: application/json" -d '{"id_usuario": '$ID_U2'}' -i -w "\n"

  # Testear PUT: Solo id_respuesta_a_mensaje
  echo -e "\nResultado de asignar respuesta a otro mensaje (valido):"
  curl -s -X PUT $BASE_URL/mensajes/$ID_M1 -H "Content-Type: application/json" -d '{"id_respuesta_a_mensaje": '$ID_M2'}' -i -w "\n"

  # Verificar el resultado final
  echo -e "\nResultado final con todos los campos actualizados (valido):"
  curl -s -X GET $BASE_URL/mensajes/$ID_M1 -i -w "\n"
}

test_edicion_mensajes