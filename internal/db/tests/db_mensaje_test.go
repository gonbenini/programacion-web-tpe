package test
import (
	"testing"
	"context"
	sqlc "foro/internal/db/sqlc"
	"database/sql"
	"foro/internal/db"
)

func TestDBmessages(t *testing.T) {

	//1. abrimos conexcion con la DB.
	db, err := db.AbrirDB()
	if err != nil {
		t.Fatalf("Error al abrir la base de datos: %s\n", err)
	}
	defer db.Close()


	//2. creamos instancias de sqlc.
	queries := sqlc.New(db)
	ctx := context.Background()


	//3. Creamos usuario y mensaje inicial de prueba
	user := sqlc.CreateUserParams{
		Nombre:      "Alberto Gonzalez",
		Mail:        "alberto.gonzalez@example.com",
		Contrasenia: "password123",
	}
	testUser, err := queries.CreateUser(ctx, user)
	mensaje := sqlc.CreateMensajeParams{
		IDUsuario: sql.NullInt32{Int32: testUser.IDUsuario, Valid: true},
		IDRespuestaAMensaje: sql.NullInt32{Valid: false},
		Texto: "Mensaje de testeo",
	}
	testMessage, err := queries.CreateMensaje(ctx, mensaje)


	//4. creamos un mensaje y verificamos que se creo correctamente
		//mensaje a insertar
		mensaje = sqlc.CreateMensajeParams{
			IDUsuario: sql.NullInt32{Int32: testUser.IDUsuario, Valid: true},
			IDRespuestaAMensaje: sql.NullInt32{Int32: testMessage.IDMensaje, Valid: true},
			Texto: "Este es un mensaje de prueba",
		}

	createdMessage, err := queries.CreateMensaje(ctx, mensaje)
	if err != nil {
		t.Errorf("Error al crear mensaje: %s\n", err)
	}
	if createdMessage.Texto != mensaje.Texto || createdMessage.IDUsuario != mensaje.IDUsuario || createdMessage.IDRespuestaAMensaje != mensaje.IDRespuestaAMensaje {
		t.Errorf("El mensaje creado no coincide con los parametros pasados")
	}


	//5. recuperamos el mensaje por su id y volvemos a verificar que coinciden los datos
	messageGet, err := queries.GetMensajeById(ctx, createdMessage.IDMensaje)
	if err != nil {
		t.Errorf("Error al obtener mensaje: %s\n", err)
	}
	if messageGet.Texto != createdMessage.Texto || messageGet.IDUsuario != createdMessage.IDUsuario || messageGet.IDRespuestaAMensaje != createdMessage.IDRespuestaAMensaje {
		t.Errorf("El mensaje recuperado no coincide con los parametros pasados")
	}


	//6. actualizamos el mensaje y verificamos que se actualizo correctamente
		//parametros de actualizacion
		updateParams := sqlc.UpdateMensajeParams{
			IDMensaje: createdMessage.IDMensaje,
			Texto:    "Esta es la modificacion del mensaje de prueba",
			MeGusta: sql.NullInt32{Int32: 5, Valid: true},
		}
	err = queries.UpdateMensaje(ctx, updateParams)
	if err != nil {
		t.Errorf("Error al actualizar mensaje: %s\n", err)
	}
	UpdatedMessage, err := queries.GetMensajeById(ctx, createdMessage.IDMensaje)
	if err != nil {
		t.Errorf("Error al obtener mensaje: %s\n", err)
	}
	if UpdatedMessage.Texto != updateParams.Texto || UpdatedMessage.MeGusta != updateParams.MeGusta || UpdatedMessage.IDMensaje != updateParams.IDMensaje {
		t.Errorf("El mensaje actualizado no coincide con los parametros pasados")
	}


	//7. listamos mensajes y verificamos que el mensaje actualizado este
	messagesList, err := queries.ListMensajes(ctx)
	if err != nil {
		t.Errorf("Error al listar mensajes: %s\n", err)
	}
	found := false
	for _, i := range messagesList {
		
		if i.IDMensaje == updateParams.IDMensaje || i.Texto == updateParams.Texto || i.MeGusta == updateParams.MeGusta {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("El mensaje actualizado no se encuentra en la lista de mensajes")
	}


	//8. eliminamos el mensaje y verificamos que ya no existe
	err = queries.DeleteMensaje(ctx, createdMessage.IDMensaje)
	if err != nil {
		t.Errorf("Error al eliminar mensaje: %s\n", err)
	}
	_, err = queries.GetUserById(ctx, createdMessage.IDMensaje)
	if err == nil {
		t.Errorf("El mensaje aun sigue estando en la DB.")
	}

	
	//9 eliminamos usuario y mensaje de prueba
		//7. eliminamos el usuario y verificamos que ya no existe
	err = queries.DeleteUser(ctx, testUser.IDUsuario)
	err = queries.DeleteMensaje(ctx, testMessage.IDMensaje)
}