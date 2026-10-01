package service

import (
	"testing"
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ---- Resumen de turnos de una paciente (próximo turno y activos) ----

// Stub: el doble sólo devuelve la lista de turnos de la paciente (acá,
// vacía). No se revisa qué se le pidió a la base.
func TestResumenDeTurnos_SinTurnos_NoHayProximo(t *testing.T) {
	// Arrange: un doble sin turnos cargados y el servicio con el reloj fijo.
	repo := &repoDoble{}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se pide el resumen de la paciente.
	resumen, err := servicio.ResumenDeTurnos(42)

	// Assert: no hay error, ni próximo turno, ni turnos activos.
	if err != nil {
		t.Fatalf("no se esperaba error y dio: %v", err)
	}
	if resumen.Proximo != nil || resumen.Activos != 0 {
		t.Errorf("sin turnos no tiene que haber próximo ni activos, y dio %+v", resumen)
	}
}

// Stub: el doble sólo devuelve la lista de turnos de la paciente. Hay un
// turno por cada camino de la función: cancelado, completado, uno que ya
// empezó y dos activos (uno mañana y otro hoy, que es el más cercano).
// El confirmado lleva FechaHoraFin para que el auto-completado no lo tome
// como terminado: sin fin, su horario "ya pasó" y el servicio intentaría
// marcarlo como completado.
func TestResumenDeTurnos_CuentaSoloLosActivosQueTodaviaNoEmpezaron(t *testing.T) {
	// Arrange: el doble con un turno por cada caso y el servicio con el reloj fijo.
	hoyALas := func(h int) time.Time { return time.Date(2026, 10, 15, h, 0, 0, 0, time.UTC) }
	repo := &repoDoble{turnosDelPaciente: []models.Turno{
		{ID: 1, Estado: models.EstadoCancelado, FechaHoraInicio: diaSiguienteALas(9, 0)},                                           // cancelado: no cuenta
		{ID: 2, Estado: models.EstadoCompletado, FechaHoraInicio: hoyALas(6)},                                                      // completado: no cuenta
		{ID: 3, Estado: models.EstadoPendiente, FechaHoraInicio: hoyALas(7)},                                                       // ya empezó: no cuenta
		{ID: 4, Estado: models.EstadoConfirmado, FechaHoraInicio: diaSiguienteALas(10, 0), FechaHoraFin: diaSiguienteALas(10, 30)}, // activo
		{ID: 5, Estado: models.EstadoPendiente, FechaHoraInicio: hoyALas(15)},                                                      // activo, hoy, el más cercano
	}}
	servicio := NuevoTurnos(repo, relojFijo) // ahora: 15/10 a las 8:00

	// Act: se pide el resumen de la paciente.
	resumen, err := servicio.ResumenDeTurnos(42)

	// Assert: cuenta sólo los dos activos, el próximo es el de hoy a las 15
	// y avisa que tiene un turno hoy.
	if err != nil {
		t.Fatalf("no se esperaba error y dio: %v", err)
	}
	if resumen.Activos != 2 || resumen.SinConfirmar != 1 {
		t.Errorf("se esperaban 2 activos y 1 sin confirmar, y dio %d y %d", resumen.Activos, resumen.SinConfirmar)
	}
	if resumen.Proximo == nil || resumen.Proximo.ID != 5 {
		t.Errorf("el próximo tenía que ser el turno 5 (hoy a las 15), y fue %+v", resumen.Proximo)
	}
	if !resumen.TieneTurnosHoy {
		t.Errorf("tiene un turno hoy a las 15 y TieneTurnosHoy dio false")
	}
}

// Stub: el doble sólo devuelve la lista de turnos de la paciente (un
// único turno, mañana).
func TestResumenDeTurnos_SinTurnosHoy(t *testing.T) {
	// Arrange: el doble con un solo turno confirmado para mañana.
	repo := &repoDoble{turnosDelPaciente: []models.Turno{
		{ID: 4, Estado: models.EstadoConfirmado, FechaHoraInicio: diaSiguienteALas(10, 0), FechaHoraFin: diaSiguienteALas(10, 30)},
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se pide el resumen de la paciente.
	resumen, _ := servicio.ResumenDeTurnos(42)

	// Assert: como el turno es mañana, no tiene turnos hoy.
	if resumen.TieneTurnosHoy {
		t.Errorf("el único turno es mañana y TieneTurnosHoy dio true")
	}
}
