package service

import (
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// Repositorio describe QUÉ necesita el negocio de la base, sin decir
// CÓMO se guarda. En la app real entra repository.Postgres; en los
// tests entra un doble en memoria, así las reglas se prueban sin
// levantar Postgres.
type Repositorio interface {
	ObtenerProfesional(id int) (*models.Profesional, error)
	ListarTurnosDeProfesionalEnFecha(profesionalID int, fecha time.Time) ([]models.Turno, error)
	ExisteSuperposicionProfesional(profesionalID int, inicio, fin time.Time) (bool, error)
	ExisteSuperposicionPaciente(pacienteID int, inicio, fin time.Time) (bool, error)
	BuscarOCrearPaciente(nombre, dni, email, telefono string) (*models.Paciente, error)
	CrearTurno(t models.Turno) (*models.Turno, error)
	ObtenerTurno(id int) (*models.Turno, error)
	ActualizarEstadoTurno(id int, nuevoEstado string) (*models.Turno, error)
	ListarTurnosDePaciente(pacienteID int) ([]models.Turno, error)
}

// Turnos agrupa las reglas de negocio del turnero junto con sus dos
// dependencias externas: la base (repo) y el reloj (ahora). Recibir el
// reloj como función permite fijar "la hora actual" en un test.
type Turnos struct {
	repo  Repositorio
	ahora func() time.Time
}

// NuevoTurnos arma el servicio. En main.go se usa con
// repository.Postgres y time.Now.
func NuevoTurnos(repo Repositorio, ahora func() time.Time) *Turnos {
	return &Turnos{repo: repo, ahora: ahora}
}
