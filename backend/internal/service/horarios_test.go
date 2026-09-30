package service

import (
	"testing"
	"time"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ---- Regla: horarios disponibles = grilla del día menos los ocupados ----

func TestHorariosDisponibles_SacaLosTurnosOcupados(t *testing.T) {
	repo := &repoDoble{
		profesional: profesionalDeManana(), // de 9 a 13, cada 30 min: 8 huecos
		turnosDelDia: []models.Turno{
			{FechaHoraInicio: diaSiguienteALas(10, 0), FechaHoraFin: diaSiguienteALas(10, 30)},
		},
	}
	servicio := NuevoTurnos(repo, relojFijo)
	dia := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)

	disponibles, err := servicio.HorariosDisponibles(1, dia)

	if err != nil {
		t.Fatalf("no se esperaba error y dio: %v", err)
	}
	if len(disponibles) != 7 {
		t.Errorf("se esperaban 7 horarios (8 menos el ocupado) y hubo %d", len(disponibles))
	}
	for _, h := range disponibles {
		if h.Equal(diaSiguienteALas(10, 0)) {
			t.Errorf("las 10:00 están ocupadas y aparecieron como disponibles")
		}
	}
}

// ---- Regla: un horario deja de ofrecerse 10 minutos antes de empezar ----

func TestHorariosDisponibles_AnticipacionMinima(t *testing.T) {
	dia := time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC)

	casos := []struct {
		nombre          string
		ahora           time.Time
		primeroEsperado time.Time
	}{
		{"a las 10:50 todavía se ofrece el de las 11:00", diaSiguienteALas(10, 50), diaSiguienteALas(11, 0)},
		{"a las 10:51 el de las 11:00 ya no se ofrece", diaSiguienteALas(10, 51), diaSiguienteALas(11, 30)},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := &repoDoble{profesional: profesionalDeManana()}
			servicio := NuevoTurnos(repo, relojEn(c.ahora))

			disponibles, err := servicio.HorariosDisponibles(1, dia)

			if err != nil {
				t.Fatalf("no se esperaba error y dio: %v", err)
			}
			if len(disponibles) == 0 || !disponibles[0].Equal(c.primeroEsperado) {
				t.Errorf("el primer horario tenía que ser %s y la lista fue %v",
					c.primeroEsperado.Format("15:04"), disponibles)
			}
		})
	}
}
