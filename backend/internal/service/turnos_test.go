package service

import (
	"errors"
	"testing"
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ---- Regla: horario de atención (función pura, sin doble) ----

// Sin doble: dentroDeHorarioAtencion es una función pura (recibe datos y
// devuelve un resultado, no habla con la base). Es un test parametrizado:
// cada fila de la tabla corre como un subtest con su propio nombre.
func TestDentroDeHorarioAtencion(t *testing.T) {
	// Arrange: un profesional de mañana y la tabla de turnos a probar.
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
			// Act: se pregunta si el turno entra en el horario de atención.
			resultado := dentroDeHorarioAtencion(c.inicio, c.fin, profesional)

			// Assert: la respuesta coincide con la esperada para este caso.
			if resultado != c.esperado {
				t.Errorf("de %s a %s: se esperaba %v y dio %v",
					c.inicio.Format("15:04"), c.fin.Format("15:04"), c.esperado, resultado)
			}
		})
	}
}

// ---- Regla: no se reserva en el pasado ----

// Stub y mock a la vez. Stub: el doble le da el profesional al servicio.
// Mock: el segundo if revisa qué se le pidió a la base (turnosCreados
// tiene que quedar vacía, o sea, no se intentó guardar nada).
func TestCrearTurno_EnElPasado_EsRechazadoYNoSeGuarda(t *testing.T) {
	// Arrange: el doble con el profesional, el servicio y una fecha de ayer.
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)
	ayer := time.Date(2026, 10, 14, 10, 0, 0, 0, time.UTC)

	// Act: se intenta reservar un turno para ayer.
	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: ayer})

	// Assert: se rechaza con ErrFechaEnElPasado y no se guarda nada.
	if !errors.Is(err, ErrFechaEnElPasado) {
		t.Errorf("se esperaba ErrFechaEnElPasado y dio: %v", err)
	}
	if len(repo.turnosCreados) != 0 {
		t.Errorf("no se tenía que guardar nada, y se guardaron %d turnos", len(repo.turnosCreados))
	}
}

// ---- Regla: horario de atención, aplicada en la reserva ----

// Stub: el doble sólo le da el profesional al servicio; el test revisa el
// error que devolvió, no qué se le pidió a la base.
//
// Diferencia con TestDentroDeHorarioAtencion: aquél prueba la función del
// horario sola, con varios casos de borde. Éste prueba que CrearTurno USE
// esa función y rechace la reserva con ErrFueraDeHorario. Si alguien
// borrara el chequeo de horario de CrearTurno, aquél seguiría en verde y
// éste se pondría rojo.
func TestCrearTurno_FueraDeHorario_EsRechazado(t *testing.T) {
	// Arrange: el doble con el profesional de mañana y el servicio.
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se intenta reservar a las 15:00, cuando ya no atiende.
	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(15, 0)})

	// Assert: se rechaza con ErrFueraDeHorario.
	if !errors.Is(err, ErrFueraDeHorario) {
		t.Errorf("se esperaba ErrFueraDeHorario y dio: %v", err)
	}
}

// ---- Reglas: superposición del profesional y del paciente ----

// Stub: el doble contesta que el profesional ya está ocupado. La
// superposición en sí la calcula Postgres con SQL; acá se prueba que el
// servicio reaccione bien a esa respuesta.
func TestCrearTurno_ProfesionalOcupado_EsRechazado(t *testing.T) {
	// Arrange: el doble dice que el profesional ya tiene un turno a esa hora.
	repo := &repoDoble{profesional: profesionalDeManana(), ocupadoProfesional: true}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se intenta reservar a las 10:00.
	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	// Assert: se rechaza con ErrSuperposicionProfesional.
	if !errors.Is(err, ErrSuperposicionProfesional) {
		t.Errorf("se esperaba ErrSuperposicionProfesional y dio: %v", err)
	}
}

// Stub: el doble contesta que la paciente ya está ocupada. Igual que el
// anterior, se prueba la reacción del servicio, no el SQL.
func TestCrearTurno_PacienteOcupado_EsRechazado(t *testing.T) {
	// Arrange: el doble dice que la paciente ya tiene un turno a esa hora.
	repo := &repoDoble{profesional: profesionalDeManana(), ocupadoPaciente: true}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se intenta reservar a las 10:00.
	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	// Assert: se rechaza con ErrSuperposicionPaciente.
	if !errors.Is(err, ErrSuperposicionPaciente) {
		t.Errorf("se esperaba ErrSuperposicionPaciente y dio: %v", err)
	}
}

// ---- Regla: el turno nuevo queda pendiente, con fin calculado y precio copiado ----

// Stub y mock. Stub: el doble le da el profesional. Mock: no se mira lo
// que devolvió CrearTurno, sino el turno que el servicio le mandó a
// guardar a la base (repo.turnosCreados).
func TestCrearTurno_Valido_QuedaPendienteConElPrecioDelProfesional(t *testing.T) {
	// Arrange: el doble con el profesional (precio 8000) y el servicio.
	repo := &repoDoble{profesional: profesionalDeManana()}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se reserva un turno válido para mañana a las 10:00.
	_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: diaSiguienteALas(10, 0)})

	// Assert: se guardó un turno pendiente, de 10:00 a 10:30, a 8000.
	// t.Fatalf corta el test acá mismo si falla: si hubo error o no se guardó
	// ningún turno, no tiene sentido seguir revisando el estado o el precio
	// (repo.turnosCreados[0] ni siquiera existiría).
	// t.Errorf, en cambio, marca el test en rojo pero sigue: así, si fallan el
	// estado y el precio a la vez, se ven los dos errores en una sola corrida.
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

// Stub: el doble le da el turno en el estado de origen. En los casos
// prohibidos también actúa como mock: se revisa que no se le haya pedido
// ningún cambio a la base (cambiosDeEstado vacía). Es un test
// parametrizado.
// "confirmado a completado" está prohibido a mano: un turno sólo llega a
// completado por el auto-completado (ver los tests de ObtenerTurno).
func TestCambiarEstadoTurno(t *testing.T) {
	// Arrange: la tabla de transiciones, permitidas y prohibidas.
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
			// Arrange (de este caso): el doble con un turno en el estado de
			// origen, y el servicio.
			repo := &repoDoble{turno: &models.Turno{
				ID:     7,
				Estado: c.desde,
				// El turno es futuro a propósito: si ya hubiera empezado, la regla de no
				// cancelar turnos pasados rechazaría los casos de cancelación.
				FechaHoraInicio: diaSiguienteALas(10, 0),
			}}
			servicio := NuevoTurnos(repo, relojFijo)

			// Act: se pide pasar el turno al estado de destino.
			_, err := servicio.CambiarEstadoTurno(7, c.hacia)

			// Assert: las permitidas no dan error; las prohibidas dan
			// ErrTransicionInvalida y no tocan la base.
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

// Mock: el assert no mira lo que devolvió ObtenerTurno, mira qué le pidió
// el servicio a la base: exactamente un cambio de estado, a completado. Es
// el equivalente al Verify(..., Times.Once) de Moq en la guía de la
// cátedra.
func TestObtenerTurno_ConfirmadoQueYaPaso_SeMarcaCompletadoUnaSolaVez(t *testing.T) {
	// Arrange: un turno confirmado que terminó hace una hora.
	repo := &repoDoble{turno: &models.Turno{
		ID:           7,
		Estado:       models.EstadoConfirmado,
		FechaHoraFin: ahoraFijo.Add(-time.Hour), // terminó hace una hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se lee el turno.
	servicio.ObtenerTurno(7)

	// Assert: se pidió un solo cambio de estado, a completado.
	if len(repo.cambiosDeEstado) != 1 || repo.cambiosDeEstado[0] != models.EstadoCompletado {
		t.Errorf("se esperaba un solo cambio a completado, y hubo: %v", repo.cambiosDeEstado)
	}
}

// Mock: revisa que el servicio NO le haya pedido ningún cambio a la base.
// Junto con el test anterior cubre los dos caminos del if del
// auto-completado (ya pasó / todavía no pasó), porque Go no mide
// cobertura de ramas.
func TestObtenerTurno_ConfirmadoQueTodaviaNoPaso_NoSeModifica(t *testing.T) {
	// Arrange: un turno confirmado que termina dentro de una hora.
	repo := &repoDoble{turno: &models.Turno{
		ID:           7,
		Estado:       models.EstadoConfirmado,
		FechaHoraFin: ahoraFijo.Add(time.Hour), // termina dentro de una hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se lee el turno.
	servicio.ObtenerTurno(7)

	// Assert: no se pidió ningún cambio de estado.
	if len(repo.cambiosDeEstado) != 0 {
		t.Errorf("un turno que no pasó no tiene que cambiar de estado, y se pidió: %v", repo.cambiosDeEstado)
	}
}

// ---- Regla: anticipación mínima de 10 minutos para reservar ----

// Stub: el doble le da el profesional. Es un test parametrizado que prueba
// el borde de la regla de 10 minutos: con 10 minutos justos se acepta y
// con 9 se rechaza. Usa relojEn para que "ahora" sean las 9:50 del día
// siguiente.
func TestCrearTurno_AnticipacionMinima(t *testing.T) {
	// Arrange: "ahora" son las 9:50 y la tabla con los dos bordes.
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
			// Arrange (de este caso): el doble con el profesional y el
			// servicio con el reloj en las 9:50.
			repo := &repoDoble{profesional: profesionalDeManana()}
			servicio := NuevoTurnos(repo, relojEn(ahora))

			// Act: se intenta reservar a la hora del caso.
			_, err := servicio.CrearTurno(CrearTurnoInput{ProfesionalID: 1, FechaHoraInicio: c.inicio})

			// Assert: el error es el esperado (nil si se acepta).
			if !errors.Is(err, c.errEsperado) {
				t.Errorf("se esperaba %v y dio: %v", c.errEsperado, err)
			}
		})
	}
}

// ---- Regla: no se puede cancelar un turno que ya empezó ----

// Stub y mock. Stub: el doble le da un turno que empezó hace media hora.
// Mock: el segundo if revisa que no se le haya pedido ningún cambio a la
// base.
func TestCambiarEstadoTurno_CancelarUnTurnoQueYaEmpezo_EsRechazado(t *testing.T) {
	// Arrange: un turno confirmado que empezó hace media hora.
	repo := &repoDoble{turno: &models.Turno{
		ID:              7,
		Estado:          models.EstadoConfirmado,
		FechaHoraInicio: ahoraFijo.Add(-30 * time.Minute), // empezó hace media hora
	}}
	servicio := NuevoTurnos(repo, relojFijo)

	// Act: se intenta cancelarlo.
	_, err := servicio.CambiarEstadoTurno(7, models.EstadoCancelado)

	// Assert: se rechaza con ErrCancelarTurnoPasado y no se toca la base.
	if !errors.Is(err, ErrCancelarTurnoPasado) {
		t.Errorf("se esperaba ErrCancelarTurnoPasado y dio: %v", err)
	}
	if len(repo.cambiosDeEstado) != 0 {
		t.Errorf("no se tenía que tocar la base, y se pidieron %v", repo.cambiosDeEstado)
	}
}
