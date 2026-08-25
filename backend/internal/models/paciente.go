package models

// Paciente representa una fila de la tabla "pacientes".
//
// Telefono es *string (puntero a string) en vez de string porque la
// columna es opcional (puede ser NULL en la base). Un puntero nil
// significa "no hay teléfono cargado"; un string común no puede
// representar "ausencia de valor" de forma distinta a "string vacío".
type Paciente struct {
	ID       int     `db:"id" json:"id"`
	Nombre   string  `db:"nombre" json:"nombre"`
	DNI      string  `db:"dni" json:"dni"`
	Email    string  `db:"email" json:"email"`
	Telefono *string `db:"telefono" json:"telefono,omitempty"`
}
