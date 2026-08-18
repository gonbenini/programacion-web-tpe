package main

import (
	"fmt"
	"net/http"
)

func main() {
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	port := ":8080"

	http.Handle("/", fileServer)
	fmt.Printf("Servidor con formulario escuchando en http://localhost%s\n", port)
	
	err := http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
