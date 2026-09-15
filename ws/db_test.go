package main
import (
	"testing"
	"context"
	_ "github.com/lib/pq"
	sqlc "foro/db/sqlc"
)

func TestDBusers(t *testing.T) {

	//1. abrimos conexcion con la DB.
	db ,err := abrirDB()
	if err != nil {
		t.Fatalf("Error al abrir la base de datos: %s\n", err)
	}
	defer db.Close()


	//2. creamos instancias de sqlc.
	queries := sqlc.New(db)
	ctx := context.Background()


	//3. creamos un usuario y verificamos que se creo correctamente
		//usuario a insertar
		user := sqlc.CreateUserParams{
			Nombre:   "Alberto Gonzalez",
			Mail: "alberto.gonzalez@example.com",
			Contrasenia: "password123",
		}

	createdUser, err := queries.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("Error al crear usuario: %s\n", err)
	}
	if createdUser.Nombre != user.Nombre || createdUser.Mail != user.Mail {
		t.Errorf("El usuario creado no coincide con los parametros pasados")
	}


	//4. recuperamos el usuario por su id y volvemos a verificar que coinciden los datos
	userGet, err := queries.GetUserById(ctx, createdUser.IDUsuario)
	if err != nil {
		t.Errorf("Error al obtener usuario: %s\n", err)
	}
	if user.Nombre != userGet.Nombre || user.Mail != userGet.Mail {
		t.Errorf("El usuario recuperado no coincide con los parametros pasados")
	}


	//5. actualizamos el usuario
		//parametros de actualizacion
		updateParams := sqlc.UpdateUserParams{
			IDUsuario: createdUser.IDUsuario,
			Nombre:    "Alberto Gonzalez Updated",
			Mail: "alberto.gonzalez.updated@example.com",
		}
	err = queries.UpdateUser(ctx, updateParams)
	if err != nil {
		t.Errorf("Error al actualizar usuario: %s\n", err)
	}


	//6. listamos usuarios y verificamos que el usuario actualizado este
	usersList, err := queries.ListUsers(ctx)
	if err != nil {
		t.Errorf("Error al listar usuarios: %s\n", err)
	}
	found := false
	for _, i := range usersList {
		
		if i.IDUsuario == createdUser.IDUsuario || i.Nombre == updateParams.Nombre || i.Mail == updateParams.Mail {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("El usuario actualizado no se encuentra en la lista de usuarios")
	}


	//7. eliminamos el usuario y verificamos que ya no existe
	err = queries.DeleteUser(ctx, createdUser.IDUsuario)
	if err != nil {
		t.Errorf("Error al eliminar usuario: %s\n", err)
	}
	_, err = queries.GetUserById(ctx, createdUser.IDUsuario)
	if err == nil {
		t.Errorf("El usuario aun sigue estando en la DB.")
	}
}