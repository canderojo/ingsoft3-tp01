package repository

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ListarTurnosDeProfesionalEnFecha trae los turnos ya ocupados (no
// cancelados) de un profesional en un día puntual. Se usa para
// calcular los huecos horarios disponibles.
func ListarTurnosDeProfesionalEnFecha(db *sqlx.DB, profesionalID int, fecha time.Time) ([]models.Turno, error) {
	turnos := []models.Turno{}
	// "fecha_hora_inicio::date" convierte el timestamp a solo fecha
	// (sin hora) para compararlo contra el día que nos pasan.
	err := db.Select(&turnos, `
		SELECT * FROM turnos
		WHERE profesional_id = $1
		  AND fecha_hora_inicio::date = $2::date
		  AND estado != 'cancelado'
		ORDER BY fecha_hora_inicio
	`, profesionalID, fecha)
	return turnos, err
}

// ExisteSuperposicionProfesional implementa la parte de base de datos
// de la regla de negocio 3 (un profesional no puede tener dos turnos
// que se solapen en el tiempo). Dos rangos [inicioA,finA) y
// [inicioB,finB) se superponen si inicioA < finB Y finA > inicioB.
func ExisteSuperposicionProfesional(db *sqlx.DB, profesionalID int, inicio, fin time.Time) (bool, error) {
	var count int
	err := db.Get(&count, `
		SELECT COUNT(*) FROM turnos
		WHERE profesional_id = $1
		  AND estado != 'cancelado'
		  AND fecha_hora_inicio < $3
		  AND fecha_hora_fin > $2
	`, profesionalID, inicio, fin)
	return count > 0, err
}

// ExisteSuperposicionPaciente es lo mismo que la anterior, pero para
// la regla 4 (un paciente no puede tener dos turnos que se solapen,
// aunque sean con profesionales distintos).
func ExisteSuperposicionPaciente(db *sqlx.DB, pacienteID int, inicio, fin time.Time) (bool, error) {
	var count int
	err := db.Get(&count, `
		SELECT COUNT(*) FROM turnos
		WHERE paciente_id = $1
		  AND estado != 'cancelado'
		  AND fecha_hora_inicio < $3
		  AND fecha_hora_fin > $2
	`, pacienteID, inicio, fin)
	return count > 0, err
}

// CrearTurno inserta un turno nuevo. El precio ya viene calculado
// (snapshot del precio_consulta del profesional) desde la capa de
// negocio, no se calcula acá.
func CrearTurno(db *sqlx.DB, t models.Turno) (*models.Turno, error) {
	var turno models.Turno
	err := db.Get(&turno, `
		INSERT INTO turnos (profesional_id, paciente_id, fecha_hora_inicio, fecha_hora_fin, estado, precio)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`, t.ProfesionalID, t.PacienteID, t.FechaHoraInicio, t.FechaHoraFin, t.Estado, t.Precio)
	if err != nil {
		return nil, err
	}
	return &turno, nil
}

func ObtenerTurno(db *sqlx.DB, id int) (*models.Turno, error) {
	var turno models.Turno
	err := db.Get(&turno, `SELECT * FROM turnos WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &turno, nil
}

func ActualizarEstadoTurno(db *sqlx.DB, id int, nuevoEstado string) (*models.Turno, error) {
	var turno models.Turno
	err := db.Get(&turno, `
		UPDATE turnos SET estado = $1 WHERE id = $2 RETURNING *
	`, nuevoEstado, id)
	if err != nil {
		return nil, err
	}
	return &turno, nil
}

// ListarTurnosDePaciente se usa en "GET /turnos?dni= o ?email=" ("mis turnos").
func ListarTurnosDePaciente(db *sqlx.DB, pacienteID int) ([]models.Turno, error) {
	turnos := []models.Turno{}
	err := db.Select(&turnos, `
		SELECT * FROM turnos WHERE paciente_id = $1 ORDER BY fecha_hora_inicio DESC
	`, pacienteID)
	return turnos, err
}
