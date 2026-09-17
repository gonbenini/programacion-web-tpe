package logic

import (
	"errors"
)

type NuevoUsuario struct {
	Nombre      string `json:"nombre"`
	Mail        string `json:"mail"`
	Contrasenia string `json:"contrasenia"`
}

func ValidarUsuario(u NuevoUsuario) error {
	if u.Nombre == "" {
		return errors.New("el nombre es obligatorio y no puede estar vacío")
	}
	if u.Mail == "" {
		return errors.New("el mail es obligatorio y no puede estar vacío")
	}
	if u.Contrasenia == "" {
		return errors.New("la contraseña es obligatoria y no puede estar vacía")
	}
	return nil
}
