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

type EditarMensaje struct {
	IdUsuario           *int32  `json:"id_usuario"`
	IdRespuestaAMensaje *int32  `json:"id_respuesta_a_mensaje"`
	Texto               *string `json:"texto"`
	MeGusta             *int32  `json:"me_gusta"`
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

func ValidarEdicionMensaje(req EditarMensaje) error {
	if req.Texto != nil && *req.Texto == "" {
		return errors.New("el texto no puede estar vacio si se intenta editar")
	}
	if req.MeGusta != nil && *req.MeGusta < 0 {
		return errors.New("la cantidad de me gusta no puede ser negativa")
	}
	if req.IdUsuario != nil && *req.IdUsuario <= 0 {
		return errors.New("el id_usuario debe ser mayor a 0")
	}
	if req.IdRespuestaAMensaje != nil && *req.IdRespuestaAMensaje <= 0 {
		return errors.New("el id_respuesta_a_mensaje debe ser mayor a 0")
	}

	return nil
}
