package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// BuscarOCrearPaciente implementa la regla de "identificación por DNI
// sin login": si ya existe un paciente con ese DNI, lo reutiliza; si
// no existe, lo crea con los datos recibidos.
func BuscarOCrearPaciente(db *sqlx.DB, nombre, dni, email, telefono string) (*models.Paciente, error) {
	var paciente models.Paciente

	err := db.Get(&paciente, `SELECT * FROM pacientes WHERE dni = $1`, dni)
	if err == nil {
		return &paciente, nil
	}
	// sql.ErrNoRows es el error especial que devuelve database/sql
	// cuando una query que esperaba una fila no encontró ninguna.
	// Cualquier otro error (por ejemplo, no poder conectar) lo
	// propagamos en vez de asumir "no existe".
	if err != sql.ErrNoRows {
		return nil, err
	}

	var tel *string
	if telefono != "" {
		tel = &telefono
	}

	err = db.Get(&paciente, `
		INSERT INTO pacientes (nombre, dni, email, telefono)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`, nombre, dni, email, tel)
	if err != nil {
		return nil, err
	}
	return &paciente, nil
}

// BuscarPacientePorDNIoEmail se usa en "GET /turnos?dni= o ?email=".
func BuscarPacientePorDNIoEmail(db *sqlx.DB, dni, email string) (*models.Paciente, error) {
	var paciente models.Paciente
	var err error

	if dni != "" {
		err = db.Get(&paciente, `SELECT * FROM pacientes WHERE dni = $1`, dni)
	} else {
		err = db.Get(&paciente, `SELECT * FROM pacientes WHERE email = $1`, email)
	}

	if err != nil {
		return nil, err
	}
	return &paciente, nil
}
