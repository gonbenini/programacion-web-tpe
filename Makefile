.PHONY: generate up down clean # Dado que estos son los comandos que queremos automatizar nosotros para operar la app, esta linea se agrega para que no interprete que estos comandos viviran en carpetas. Sino que estaran defindos ahi. En caso que no este este comando y existiera una carpeta clean, al llaar a make, se iria a buscar ahi a esa carpeta, que si tiene otra cosa daria error.

# 1.	Genera el código de Go leyendo el schema y queries antes de que Docker Compose se levante, los contenedores dependen de estos assets.
generate:
	docker run --rm -v "$(shell pwd)/ws:/src" -w /src sqlc/sqlc generate
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
	docker compose down -v
	# si la intencion es solo detener la instancia podriamos hacer `docker compose stop` o docker compose down` sin eliminar los volumenes que contiene la informacion generada

# 4. 	Elimina los archivos autogenerados para limpiar el proyecto local, la idea es borrar todos los assets generados.
clean:
	rm -f ws/db/sqlc/*.go
	#  tambien podria ser implementado como un comando de git para resetear a como viene el la carpeta pero perderiamos codigo que estemos desarrollando en vez de solo los assets generados por la instancia.