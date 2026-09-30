package service

import (
	"errors"
	"testing"
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ---- Regla: horario de atención (función pura, sin doble) ----

func TestDentroDeHorarioAtencion(t *testing.T) {
	profesional := *profesionalDeManana() // atiende de 9:00 a 13:00

	casos := []struct {
		nombre   string
		inicio   time.Time
		fin      time.Time
		esperado bool
	}{
		{"turno en el medio de la mañana entra", diaSiguienteALas(10, 0), diaSiguienteALas(10, 30), true},
		{"turno que empieza justo a la apertura entra", diaSiguienteALas(9, 0), diaSiguienteALas(9, 30), true},
		{"turno que termina justo al cierre entra", diaSiguienteALas(12, 30), diaSiguienteALas(13, 0), true},
		{"turno que empieza antes de la apertura no entra", diaSiguienteALas(8, 45), diaSiguienteALas(9, 15), false},
		{"turno que termina después del cierre no entra", diaSiguienteALas(12, 45), diaSiguienteALas(13, 15), false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			resultado := dentroDeHorarioAtencion(c.inicio, c.fin, profesional)

			if resultado != c.esperado {
				t.Errorf("de %s a %s: se esperaba %v y dio %v",
					c.inicio.Format("15:04"), c.fin.Format("15:04"), c.esperado, resultado)
			}
		})
	}
}

// ---- Regla: no se reserva en el pasado ----

func TestCrearTurno_EnElPasado_EsRechazadoYNoSeGuarda(t *testing.T) {
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)
	ayer := time.Date(2026, 10, 14, 10, 0, 0, 0, time.UTC)

	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: ayer})

	if !errors.Is(err, ErrFechaEnElPasado) {
		t.Errorf("se esperaba ErrFechaEnElPasado y dio: %v", err)
	}
	if len(repo.turnosCreados) != 0 {
		t.Errorf("no se tenía que guardar nada, y se guardaron %d turnos", len(repo.turnosCreados))
	}
}

// ---- Regla: horario de atención, aplicada en la reserva ----

func TestCrearTurno_FueraDeHorario_EsRechazado(t *testing.T) {
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)

	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(15, 0)})

	if !errors.Is(err, ErrFueraDeHorario) {
		t.Errorf("se esperaba ErrFueraDeHorario y dio: %v", err)
	}
}

// ---- Reglas: superposición del profesional y del paciente ----

func TestCrearTurno_ProfesionalOcupado_EsRechazado(t *testing.T) {
	repo := &repoDoble{profesional: profesionalDeManana(), ocupadoProfesional: true}
	servicio := NuevoTurnos(repo, relojFijo)

	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	if !errors.Is(err, ErrSuperposicionProfesional) {
		t.Errorf("se esperaba ErrSuperposicionProfesional y dio: %v", err)
	}
}

func TestCrearTurno_PacienteOcupado_EsRechazado(t *testing.T) {
	repo := &repoDoble{profesional: profesionalDeManana(), ocupadoPaciente: true}
	servicio := NuevoTurnos(repo, relojFijo)

	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	if !errors.Is(err, ErrSuperposicionPaciente) {
		t.Errorf("se esperaba ErrSuperposicionPaciente y dio: %v", err)
	}
}

// ---- Regla: el turno nuevo queda pendiente, con fin calculado y precio copiado ----

func TestCrearTurno_Valido_QuedaPendienteConElPrecioDelProfesional(t *testing.T) {
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)

	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	if err != nil {
		t.Fatalf("no se esperaba error y dio: %v", err)
	}
	if len(repo.turnosCreados) != 1 {
		t.Fatalf("se esperaba guardar 1 turno y se guardaron %d", len(repo.turnosCreados))
	}
	guardado := repo.turnosCreados[0]
	if guardado.Estado != models.EstadoPendiente {
		t.Errorf("el estado tenía que ser pendiente y fue %q", guardado.Estado)
	}
	if !guardado.FechaHoraFin.Equal(diaSiguienteALas(10, 30)) {
		t.Errorf("el fin tenía que ser 10:30 (inicio + 30 min) y fue %s", guardado.FechaHoraFin.Format("15:04"))
	}
	if guardado.Precio != 8000 {
		t.Errorf("el precio tenía que copiarse del profesional (8000) y fue %v", guardado.Precio)
	}
}

// ---- Regla: máquina de estados ----

func TestCambiarEstadoTurno(t *testing.T) {
	casos := []struct {
		nombre    string
		desde     string
		hacia     string
		permitido bool
	}{
		{"pendiente a confirmado se permite", models.EstadoPendiente, models.EstadoConfirmado, true},
		{"pendiente a cancelado se permite", models.EstadoPendiente, models.EstadoCancelado, true},
		{"confirmado a cancelado se permite", models.EstadoConfirmado, models.EstadoCancelado, true},
		{"pendiente a completado no se permite", models.EstadoPendiente, models.EstadoCompletado, false},
		{"confirmado a completado no se permite a mano", models.EstadoConfirmado, models.EstadoCompletado, false},
		{"cancelado no vuelve a pendiente", models.EstadoCancelado, models.EstadoPendiente, false},
		{"completado no se puede cancelar", models.EstadoCompletado, models.EstadoCancelado, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := &repoDoble{turno: &models.Turno{
				ID:              7,
				Estado:          c.desde,
				FechaHoraInicio: diaSiguienteALas(10, 0), // un turno futuro
			}}
			servicio := NuevoTurnos(repo, relojFijo)

			_, err := servicio.CambiarEstadoTurno(7, c.hacia)

			if c.permitido && err != nil {
				t.Errorf("se esperaba que se permita y dio error: %v", err)
			}
			if !c.permitido && !errors.Is(err, ErrTransicionInvalida) {
				t.Errorf("se esperaba ErrTransicionInvalida y dio: %v", err)
			}
			if !c.permitido && len(repo.cambiosDeEstado) != 0 {
				t.Errorf("una transición inválida no tiene que tocar la base, y se pidieron %v", repo.cambiosDeEstado)
			}
		})
	}
}

// ---- Regla: auto-completado de turnos confirmados que ya pasaron (mock) ----

func TestObtenerTurno_ConfirmadoQueYaPaso_SeMarcaCompletadoUnaSolaVez(t *testing.T) {
	repo := &repoDoble{turno: &models.Turno{
		ID:           7,
		Estado:       models.EstadoConfirmado,
		FechaHoraFin: ahoraFijo.Add(-time.Hour), // terminó hace una hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	servicio.ObtenerTurno(7)

	if len(repo.cambiosDeEstado) != 1 || repo.cambiosDeEstado[0] != models.EstadoCompletado {
		t.Errorf("se esperaba un solo cambio a completado, y hubo: %v", repo.cambiosDeEstado)
	}
}

func TestObtenerTurno_ConfirmadoQueTodaviaNoPaso_NoSeModifica(t *testing.T) {
	repo := &repoDoble{turno: &models.Turno{
		ID:           7,
		Estado:       models.EstadoConfirmado,
		FechaHoraFin: ahoraFijo.Add(time.Hour), // termina dentro de una hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	servicio.ObtenerTurno(7)

	if len(repo.cambiosDeEstado) != 0 {
		t.Errorf("un turno que no pasó no tiene que cambiar de estado, y se pidió: %v", repo.cambiosDeEstado)
	}
}

// ---- Regla: anticipación mínima de 10 minutos para reservar ----

func TestCrearTurno_AnticipacionMinima(t *testing.T) {
	ahora := diaSiguienteALas(9, 50)

	casos := []struct {
		nombre      string
		inicio      time.Time
		errEsperado error
	}{
		{"con 10 minutos justos de anticipación se acepta", diaSiguienteALas(10, 0), nil},
		{"con 9 minutos de anticipación se rechaza", diaSiguienteALas(9, 59), ErrAnticipacionInsuficiente},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := &repoDoble{profesional: profesionalDeManana()}
			servicio := NuevoTurnos(repo, relojEn(ahora))

			_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: c.inicio})

			if !errors.Is(err, c.errEsperado) {
				t.Errorf("se esperaba %v y dio: %v", c.errEsperado, err)
			}
		})
	}
}

// ---- Regla: no se puede cancelar un turno que ya empezó ----

func TestCambiarEstadoTurno_CancelarUnTurnoQueYaEmpezo_EsRechazado(t *testing.T) {
	repo := &repoDoble{turno: &models.Turno{
		ID:              7,
		Estado:          models.EstadoConfirmado,
		FechaHoraInicio: ahoraFijo.Add(-30 * time.Minute), // empezó hace media hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	_, err := servicio.CambiarEstadoTurno(7, models.EstadoCancelado)

	if !errors.Is(err, ErrCancelarTurnoPasado) {
		t.Errorf("se esperaba ErrCancelarTurnoPasado y dio: %v", err)
	}
	if len(repo.cambiosDeEstado) != 0 {
		t.Errorf("no se tenía que tocar la base, y se pidieron %v", repo.cambiosDeEstado)
	}
}
