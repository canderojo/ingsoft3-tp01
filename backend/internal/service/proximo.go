package service

import (
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ResumenDePaciente junta lo que la pantalla "Mis turnos" necesita mostrar
// arriba de la lista: cuál es el próximo turno de la paciente y cuántos
// turnos activos tiene.
type ResumenDePaciente struct {
	Proximo        *models.Turno
	Activos        int
	SinConfirmar   int
	TieneTurnosHoy bool
}

// ResumenDeTurnos arma el resumen de una paciente. El próximo turno es el
// activo (pendiente o confirmado) que empieza antes, entre los que todavía
// no empezaron; si no tiene ninguno, Proximo queda en nil.
func (s *Turnos) ResumenDeTurnos(pacienteID int) (*ResumenDePaciente, error) {
	turnos, err := s.ListarTurnosDePaciente(pacienteID)
	if err != nil {
		return nil, err
	}

	ahora := s.ahora()
	resumen := &ResumenDePaciente{}
	for i := range turnos {
		t := &turnos[i]
		if t.Estado != models.EstadoPendiente && t.Estado != models.EstadoConfirmado {
			continue
		}
		if !t.FechaHoraInicio.After(ahora) {
			continue
		}

		resumen.Activos++
		if t.Estado == models.EstadoPendiente {
			resumen.SinConfirmar++
		}
		if mismoDia(t.FechaHoraInicio, ahora) {
			resumen.TieneTurnosHoy = true
		}
		if resumen.Proximo == nil || t.FechaHoraInicio.Before(resumen.Proximo.FechaHoraInicio) {
			resumen.Proximo = t
		}
	}
	return resumen, nil
}

// mismoDia dice si dos horarios caen en la misma fecha del calendario.
func mismoDia(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}
