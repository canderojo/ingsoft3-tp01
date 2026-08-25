package models

import "time"

// Turno representa una fila de la tabla "turnos".
type Turno struct {
	ID              int       `db:"id" json:"id"`
	ProfesionalID   int       `db:"profesional_id" json:"profesional_id"`
	PacienteID      int       `db:"paciente_id" json:"paciente_id"`
	FechaHoraInicio time.Time `db:"fecha_hora_inicio" json:"fecha_hora_inicio"`
	FechaHoraFin    time.Time `db:"fecha_hora_fin" json:"fecha_hora_fin"`
	Estado          string    `db:"estado" json:"estado"`
	Precio          float64   `db:"precio" json:"precio"`
}

// Estados válidos de un turno y las transiciones permitidas entre
// ellos (regla de negocio 5: desde "pendiente" solo se puede pasar a
// "confirmado" o "cancelado", nunca directo a "completado").
const (
	EstadoPendiente  = "pendiente"
	EstadoConfirmado = "confirmado"
	EstadoCancelado  = "cancelado"
	EstadoCompletado = "completado"
)

// TransicionesPermitidas mapea cada estado a la lista de estados a
// los que puede pasar. Es un "mapa" (map[string][]string): la
// estructura de datos de Go para diccionarios clave→valor.
var TransicionesPermitidas = map[string][]string{
	EstadoPendiente:  {EstadoConfirmado, EstadoCancelado},
	EstadoConfirmado: {EstadoCompletado, EstadoCancelado},
}
