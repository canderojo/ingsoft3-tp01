package service

import (
	"database/sql"
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// repoDoble reemplaza a repository.Postgres en los tests: cumple la
// interfaz Repositorio sin tocar ninguna base.
//
// Hace dos trabajos:
//   - responde lo que el test le cargó antes (stub);
//   - anota lo que el service le pidió (turnosCreados, cambiosDeEstado),
//     para que el test pueda verificar la interacción (mock).
type repoDoble struct {
	// Lo que responde
	profesional        *models.Profesional
	turno              *models.Turno
	turnosDelDia       []models.Turno
	ocupadoProfesional bool
	ocupadoPaciente    bool

	// Lo que registra
	turnosCreados   []models.Turno
	cambiosDeEstado []string
}

func (r *repoDoble) ObtenerProfesional(id int) (*models.Profesional, error) {
	if r.profesional == nil {
		return nil, sql.ErrNoRows
	}
	return r.profesional, nil
}

func (r *repoDoble) ListarTurnosDeProfesionalEnFecha(profesionalID int, fecha time.Time) ([]models.Turno, error) {
	return r.turnosDelDia, nil
}

func (r *repoDoble) ExisteSuperposicionProfesional(profesionalID int, inicio, fin time.Time) (bool, error) {
	return r.ocupadoProfesional, nil
}

func (r *repoDoble) ExisteSuperposicionPaciente(pacienteID int, inicio, fin time.Time) (bool, error) {
	return r.ocupadoPaciente, nil
}

func (r *repoDoble) BuscarOCrearPaciente(nombre, dni, email, telefono string) (*models.Paciente, error) {
	return &models.Paciente{ID: 42, Nombre: nombre, DNI: dni, Email: email}, nil
}

func (r *repoDoble) CrearTurno(t models.Turno) (*models.Turno, error) {
	r.turnosCreados = append(r.turnosCreados, t)
	t.ID = 1
	return &t, nil
}

func (r *repoDoble) ObtenerTurno(id int) (*models.Turno, error) {
	if r.turno == nil {
		return nil, sql.ErrNoRows
	}
	copia := *r.turno
	return &copia, nil
}

func (r *repoDoble) ActualizarEstadoTurno(id int, nuevoEstado string) (*models.Turno, error) {
	r.cambiosDeEstado = append(r.cambiosDeEstado, nuevoEstado)
	actualizado := *r.turno
	actualizado.Estado = nuevoEstado
	return &actualizado, nil
}

func (r *repoDoble) ListarTurnosDePaciente(pacienteID int) ([]models.Turno, error) {
	return nil, nil
}

// ---- Datos de prueba compartidos ----

// ahoraFijo es el "ahora" de todos los tests: el 15/10/2026 a las 8:00.
// Con un reloj fijo, los tests dan lo mismo hoy que dentro de un año.
var ahoraFijo = time.Date(2026, 10, 15, 8, 0, 0, 0, time.UTC)

func relojFijo() time.Time { return ahoraFijo }

// relojEn arma un reloj que siempre devuelve la hora pedida. Sirve para
// los tests que necesitan un "ahora" distinto de ahoraFijo.
func relojEn(ahora time.Time) func() time.Time {
	return func() time.Time { return ahora }
}

// hora arma un HoraDelDia (una hora sin fecha, como la columna TIME).
func hora(h, m int) models.HoraDelDia {
	return models.HoraDelDia{Time: time.Date(0, 1, 1, h, m, 0, 0, time.UTC)}
}

// profesionalDeManana atiende de 9 a 13, con turnos de 30 minutos.
func profesionalDeManana() *models.Profesional {
	return &models.Profesional{
		ID:                   1,
		HoraInicioAtencion:   hora(9, 0),
		HoraFinAtencion:      hora(13, 0),
		DuracionTurnoMinutos: 30,
		PrecioConsulta:       8000,
	}
}

// diaSiguienteALas arma una fecha del día siguiente al "ahora" fijo
// (el 16/10/2026), a la hora pedida. Sirve para reservar "mañana".
func diaSiguienteALas(h, m int) time.Time {
	return time.Date(2026, 10, 16, h, m, 0, 0, time.UTC)
}
