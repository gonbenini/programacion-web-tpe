package logic

import (
	"errors"
)

// NuevoMensaje representa el payload (JSON) esperado para crear un mensaje
type NuevoMensaje struct {
	IdUsuario           int32  `json:"id_usuario"`
	IdRespuestaAMensaje *int32 `json:"id_respuesta_a_mensaje"`
	Texto               string `json:"texto"`
}

// ValidarMensaje verifica que los campos obligatorios del mensaje estén presentes y sean correctos
func ValidarMensaje(req NuevoMensaje) error {
	if req.Texto == "" {
		return errors.New("el texto del mensaje es obligatorio")
	}

	if req.IdUsuario <= 0 {
		return errors.New("el id_usuario es invalido o requerido")
	}

	return nil
}
