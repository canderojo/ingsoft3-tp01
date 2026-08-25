package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Profesional representa una fila de la tabla "profesionales".
//
// Las etiquetas `db:"..."` le dicen a sqlx a qué columna de la base
// corresponde cada campo cuando arma un Profesional a partir de una
// fila devuelta por una query (esto es lo que reemplaza al "mapeo
// automático" que haría un ORM, pero acá es explícito y vos lo ves).
//
// Las etiquetas `json:"..."` controlan el nombre de cada campo cuando
// este struct se convierte a JSON para la respuesta HTTP.
type Profesional struct {
	ID                   int        `db:"id" json:"id"`
	Nombre               string     `db:"nombre" json:"nombre"`
	Especialidad         string     `db:"especialidad" json:"especialidad"`
	HoraInicioAtencion   HoraDelDia `db:"hora_inicio_atencion" json:"hora_inicio_atencion"`
	HoraFinAtencion      HoraDelDia `db:"hora_fin_atencion" json:"hora_fin_atencion"`
	DuracionTurnoMinutos int        `db:"duracion_turno_minutos" json:"duracion_turno_minutos"`
	PrecioConsulta       float64    `db:"precio_consulta" json:"precio_consulta"`
}

// HoraDelDia envuelve una columna TIME de Postgres. El driver pgx (en
// el modo database/sql que usa sqlx) devuelve las columnas TIME como
// texto plano ("09:00:00"), no como time.Time, así que hace falta este
// tipo intermedio con su propio Scan para poder mapearlas y después
// operar con métodos como Hour()/Minute().
type HoraDelDia struct {
	time.Time
}

const formatoHoraDelDia = "15:04:05"

// Scan implementa sql.Scanner: lo llama sqlx/database/sql al leer la
// columna desde la base.
func (h *HoraDelDia) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("HoraDelDia.Scan: tipo inesperado %T", value)
	}
	t, err := time.Parse(formatoHoraDelDia, s)
	if err != nil {
		return err
	}
	h.Time = t
	return nil
}

// Value implementa driver.Valuer: lo llama sqlx/database/sql al
// mandar este campo como parámetro de una query (no se usa hoy, pero
// hace falta para que el tipo sea simétrico).
func (h HoraDelDia) Value() (driver.Value, error) {
	return h.Time.Format(formatoHoraDelDia), nil
}

// MarshalJSON hace que HoraDelDia se serialice como "09:00:00" en vez
// del formato RFC3339 completo que usaría time.Time por defecto.
func (h HoraDelDia) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.Time.Format(formatoHoraDelDia) + `"`), nil
}

// Especialidades válidas. Usamos constantes en vez de "strings mágicos"
// sueltos por el código, así el compilador nos avisa si escribimos mal
// alguna en otro archivo.
const (
	EspecialidadDermatologia   = "dermatologia"
	EspecialidadNutricion      = "nutricion"
	EspecialidadEcografia      = "ecografia"
	EspecialidadEndocrinologia = "endocrinologia"
)
