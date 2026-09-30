package service

import (
	"testing"
	"time"
)

// ---- Reloj de la app: hora de reloj de Argentina, etiquetada como UTC ----

func TestHoraDeRelojArgentina(t *testing.T) {
	casos := []struct {
		nombre   string
		instante time.Time
		esperado time.Time
	}{
		{"las 21:30 UTC son las 18:30 en Argentina",
			time.Date(2026, 9, 30, 21, 30, 0, 0, time.UTC),
			time.Date(2026, 9, 30, 18, 30, 0, 0, time.UTC)},
		{"la 1:00 UTC del 1/10 son las 22:00 del 30/9 en Argentina",
			time.Date(2026, 10, 1, 1, 0, 0, 0, time.UTC),
			time.Date(2026, 9, 30, 22, 0, 0, 0, time.UTC)},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			resultado := horaDeRelojArgentina(c.instante)

			if !resultado.Equal(c.esperado) {
				t.Errorf("se esperaba %s y dio %s", c.esperado, resultado)
			}
		})
	}
}
