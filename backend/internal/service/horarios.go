package service

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/repository"
)

// HorariosDisponibles calcula los huecos horarios libres de un
// profesional en una fecha puntual: arma la grilla completa de
// turnos posibles (horario de atención dividido en bloques de
// duración fija) y le saca los que ya están ocupados.
func HorariosDisponibles(db *sqlx.DB, profesionalID int, fecha time.Time) ([]time.Time, error) {
	profesional, err := repository.ObtenerProfesional(db, profesionalID)
	if err == sql.ErrNoRows {
		return nil, ErrProfesionalNoExiste
	}
	if err != nil {
		return nil, err
	}

	ocupados, err := repository.ListarTurnosDeProfesionalEnFecha(db, profesionalID, fecha)
	if err != nil {
		return nil, err
	}

	duracion := time.Duration(profesional.DuracionTurnoMinutos) * time.Minute
	inicioAtencion := combinarFechaYHora(fecha, profesional.HoraInicioAtencion.Time)
	finAtencion := combinarFechaYHora(fecha, profesional.HoraFinAtencion.Time)

	disponibles := []time.Time{}
	ahora := time.Now()

	for inicio := inicioAtencion; !inicio.Add(duracion).After(finAtencion); inicio = inicio.Add(duracion) {
		fin := inicio.Add(duracion)

		if inicio.Before(ahora) {
			continue
		}
		if seSuperponeConAlguno(inicio, fin, ocupados) {
			continue
		}

		disponibles = append(disponibles, inicio)
	}

	return disponibles, nil
}

// combinarFechaYHora arma un time.Time con la fecha del primer
// argumento y la hora del segundo (se usa para "aterrizar" un
// HoraDelDia, que no tiene fecha real, sobre el día que se está
// consultando).
func combinarFechaYHora(fecha, hora time.Time) time.Time {
	return time.Date(fecha.Year(), fecha.Month(), fecha.Day(), hora.Hour(), hora.Minute(), hora.Second(), 0, fecha.Location())
}

func seSuperponeConAlguno(inicio, fin time.Time, turnos []models.Turno) bool {
	for _, t := range turnos {
		if inicio.Before(t.FechaHoraFin) && fin.After(t.FechaHoraInicio) {
			return true
		}
	}
	return false
}
