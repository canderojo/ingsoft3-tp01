package service

import (
	"testing"
	"time"
)

// ---- Reloj de la app: hora de reloj de Argentina, etiquetada como UTC ----

// Sin doble: horaDeRelojArgentina es una función pura. Es un test
// parametrizado. El segundo caso prueba el cambio de día: la 1:00 UTC del
// 1/10 todavía es el 30/9 en Argentina.
func TestHoraDeRelojArgentina(t *testing.T) {
	// Arrange: la tabla de instantes UTC y su hora de reloj en Argentina.
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
			// Act: se pasa el instante a hora de reloj de Argentina.
			resultado := horaDeRelojArgentina(c.instante)

			// Assert: da la hora (y el día) esperados.
			if !resultado.Equal(c.esperado) {
				t.Errorf("se esperaba %s y dio %s", c.esperado, resultado)
			}
		})
	}
}
