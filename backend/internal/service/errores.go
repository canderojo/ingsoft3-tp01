// Package service contiene las reglas de negocio del turnero: nada
// de SQL acá (eso vive en repository), nada de HTTP acá (eso vive en
// handlers). Esta separación es la que permite testear las reglas
// (horarios, superposición, máquina de estados) sin levantar Postgres
// ni un servidor HTTP.
package service

import "errors"

// Errores de negocio conocidos. Los handlers los traducen a códigos
// HTTP específicos (ver responderErrorDeNegocio); cualquier otro error
// (por ejemplo, uno de conexión a la base) se trata como error interno.
var (
	ErrProfesionalNoExiste      = errors.New("el profesional no existe")
	ErrFueraDeHorario           = errors.New("el turno está fuera del horario de atención del profesional")
	ErrFechaEnElPasado          = errors.New("no se puede reservar un turno en el pasado")
	ErrSuperposicionProfesional = errors.New("el profesional ya tiene un turno en ese horario")
	ErrSuperposicionPaciente    = errors.New("el paciente ya tiene un turno en ese horario")
	ErrTurnoNoExiste            = errors.New("el turno no existe")
	ErrTransicionInvalida       = errors.New("transición de estado no permitida")
)
