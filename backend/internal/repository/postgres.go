package repository

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// Postgres envuelve la conexión a la base y expone como métodos las
// funciones de este paquete que usa la capa service. No agrega lógica:
// cada método delega en la función existente pasándole p.DB. Existe
// para que service dependa de una interfaz (service.Repositorio) y no
// de *sqlx.DB, y así poder reemplazar la base por un doble en los tests.
type Postgres struct {
	DB *sqlx.DB
}

func (p Postgres) ObtenerProfesional(id int) (*models.Profesional, error) {
	return ObtenerProfesional(p.DB, id)
}

func (p Postgres) ListarTurnosDeProfesionalEnFecha(profesionalID int, fecha time.Time) ([]models.Turno, error) {
	return ListarTurnosDeProfesionalEnFecha(p.DB, profesionalID, fecha)
}

func (p Postgres) ExisteSuperposicionProfesional(profesionalID int, inicio, fin time.Time) (bool, error) {
	return ExisteSuperposicionProfesional(p.DB, profesionalID, inicio, fin)
}

func (p Postgres) ExisteSuperposicionPaciente(pacienteID int, inicio, fin time.Time) (bool, error) {
	return ExisteSuperposicionPaciente(p.DB, pacienteID, inicio, fin)
}

func (p Postgres) BuscarOCrearPaciente(nombre, dni, email, telefono string) (*models.Paciente, error) {
	return BuscarOCrearPaciente(p.DB, nombre, dni, email, telefono)
}

func (p Postgres) CrearTurno(t models.Turno) (*models.Turno, error) {
	return CrearTurno(p.DB, t)
}

func (p Postgres) ObtenerTurno(id int) (*models.Turno, error) {
	return ObtenerTurno(p.DB, id)
}

func (p Postgres) ActualizarEstadoTurno(id int, nuevoEstado string) (*models.Turno, error) {
	return ActualizarEstadoTurno(p.DB, id, nuevoEstado)
}

func (p Postgres) ListarTurnosDePaciente(pacienteID int) ([]models.Turno, error) {
	return ListarTurnosDePaciente(p.DB, pacienteID)
}
