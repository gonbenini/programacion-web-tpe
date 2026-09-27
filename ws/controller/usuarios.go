package controller

func createUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	var req logic.NuevoUsuario

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al procesar JSON", http.StatusBadRequest)
		return
	}

	err = logic.ValidarUsuario(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Nombre:      req.Nombre,
			Mail:        req.Mail,
			Contrasenia: req.Contrasenia,
		})
	if err != nil {
		fmt.Printf("Error al crear usuario: %s\n", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func getUsuarios(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries) {
	ctx := context.Background()

	users, err := queries.ListUsers(ctx) // Buscamos todos los registros de usuario
	if err != nil {
		http.Error(w, "Error al listar usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos array JSON
	json.NewEncoder(w).Encode(users)
}

func getUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()

	user, err := queries.GetUserById(ctx, id) // Buscamos el registro
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound) // Si no existe devolvemos un 404 Not Found
		return
	}

	w.Header().Set("Content-Type", "application/json") // devolvemos JSON con estado 200
	json.NewEncoder(w).Encode(user)
}

func updateUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {

}

func deleteUsuario(w http.ResponseWriter, r *http.Request, queries *sqlc.Queries, id int32) {
	ctx := context.Background()
	err := queries.DeleteUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}