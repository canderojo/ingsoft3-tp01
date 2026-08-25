// Package repository es la única capa del backend que sabe escribir
// SQL. Los handlers HTTP y las reglas de negocio nunca tocan la base
// directo: le piden datos a este paquete. Esto es lo que nos permite
// testear las reglas de negocio más adelante sin necesitar Postgres
// corriendo (le pasamos datos de prueba en memoria en vez de llamar
// a estas funciones).
package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/models"
)

// ListarProfesionales trae profesionales, filtrando por especialidad
// si se pasa un valor (string vacío = sin filtro, trae todos).
func ListarProfesionales(db *sqlx.DB, especialidad string) ([]models.Profesional, error) {
	// Inicializamos el slice vacío (en vez de nil) para que, si no hay
	// resultados, la respuesta JSON sea "[]" y no "null".
	profesionales := []models.Profesional{}

	query := `SELECT * FROM profesionales`
	args := []interface{}{}
	if especialidad != "" {
		query += ` WHERE especialidad = $1`
		args = append(args, especialidad)
	}
	query += ` ORDER BY nombre`

	// db.Select ejecuta la query y mapea cada fila del resultado a un
	// elemento del slice que le pasamos por puntero (&profesionales).
	err := db.Select(&profesionales, query, args...)
	return profesionales, err
}

// ObtenerProfesional busca un profesional por id. Devuelve un puntero
// nil junto con el error si no existe (el error será sql.ErrNoRows).
func ObtenerProfesional(db *sqlx.DB, id int) (*models.Profesional, error) {
	var profesional models.Profesional
	// db.Get es como db.Select pero espera (y devuelve) una sola fila.
	err := db.Get(&profesional, `SELECT * FROM profesionales WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &profesional, nil
}
