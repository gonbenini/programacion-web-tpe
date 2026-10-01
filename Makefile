.PHONY: generate up down clean # Dado que estos son los comandos que queremos automatizar nosotros para operar la app, esta linea se agrega para que no interprete que estos comandos viviran en carpetas. Sino que estaran defindos ahi. En caso que no este este comando y existiera una carpeta clean, al llaar a make, se iria a buscar ahi a esa carpeta, que si tiene otra cosa daria error.
# poner @ antes del comando evita que se imprima el comando

# 1.	Genera el código de Go leyendo el schema y queries antes de que Docker Compose se levante, los contenedores dependen de estos assets.
generate:
	@docker run --rm -v "$(shell pwd):/src" -w /src sqlc/sqlc generate >/dev/null
# --rm (remove): docker borra el contenedor en cuanto termina de ejecutar el comando.
# -v (mapear volumen): pasamos codigo que queremos ejecutar adentro del contenedor.
# -w /src (workdir): carpeta donde estara parada la terminal cuando empiece a ejecutar el contendor, no es necesario pero facilita ubicarnos.
# imagen: usamos la imagen sqlc en un contenedor donde queda aislada esta libreria que usamos solo para generar el codigo. que se llame sqlc/sqlc es por el estandar de dockerhub <usuario que creo la immamgen>/<nombre de la imagen>
# comando generate: cuando se levanta el contenedor, la tarea que se le pasa generate. este comando es interpretado por el entrypoint del contenedor.

# 2. 	Genera los archivos (item 1.) y luego levanta la BD y el webserver (los contenedores que presentan las capas de la app.)
up: generate
	docker compose up
	# podriammos optimizar el despliegue sacando el generate de .PHONY y por lo tanto obteniendo un despliegue mas rapido. En este caso como el desarrollo es continuo y es una instancia livianta, elegimos matentenerlo para reducir la complejidad, y evitar que el codigo Go y la base de datos nunca queden desfasados.

# 3. 	Frena los contenedores y destruye los volúmenes (reinicia la BD)
down:
	@docker compose down -v
# si la intencion es solo detener la instancia podriamos hacer `docker compose stop` o docker compose down` sin eliminar los volumenes que contiene la informacion generada

# 4. 	Elimina los archivos autogenerados para limpiar el proyecto local, la idea es borrar todos los assets generados.
# Comando git clean -ndx: -n no borra, solo lista; -d muestra directorios; x muestra tambien lo escondido por el .gitignore
#					-fdx: en cambio poner -f en vez de -n, confirma el borrado. git por defecto pide confirmacion para esto, sino se niega.
# Borro los archivos desde adentro de un contenedor para no tener que hacer sudo
clean:
	echo "Se van a borrar los siguientes archivos:"
	@git clean -ndx | awk '{print $3}'
	@docker compose run --rm webserver rm -rf internal/db/sqlc >/dev/null 2>&1
	@-git clean -fdx >/dev/null 2>&1
# Dato importante, el make se ejecuta en sh, no en bash. Por lo tanto no existe `&>/dev/null` para redirigir error y salida estandar al mismo archivo
# Si hacemos `&>/dev/null` manda el proceso a segundo plano y sigue con el siguiente. En cambio hay que usar `>/dev/null 2>&1` esta se llama "direccioon estandar POSIX"

# 5.    Ejecuta los tests.
test: generate
	docker compose run --rm webserver go test -v ./internal/db/tests
# Usamos 'docker compose run' para levantar un contenedor efímero basado en la configuración de 'webserver'.
# Como 'webserver' depende de 'db' (depends_on), Docker Compose se asegurará de que la BD esté levantada antes de correr los tests.
# Al usar --rm, el contenedor efímero donde corrieron los tests se elimina automáticamente al terminar.
# Respecto al comando `go test -v ./`, es una herramienta nativa de go que ejecuta los test.
# Lo que hace es buscar archivos cuyos nombres terminen en _test.go y ejecuta las funciones que tenga. -v es verbose, el primer `./` es para que empiece a buscar de la carpeta actual.

test-api: generate
	@docker compose up -d 
	# esperando inicializacion de servicios
	@until curl -fs http://localhost:8080/api/usuarios >/dev/null; do sleep 1; done
	@./requests.sh
	@docker compose stop >/dev/null 2>&1
	@$(MAKE) clean >/dev/null 2>&1
	@$(MAKE) down >/dev/null 2>&1
# el make clean quizas deberia ir dentro de un trap por si falla algo para que se ejecute igual.
# Es buena practica en vez de hacer `@make clean &>/dev/null`, usar `$(MAKE)`