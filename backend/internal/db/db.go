// Package db se encarga de abrir la conexión a Postgres.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// El "_" antes del import significa "importalo solo por su efecto
	// secundario". pgx se registra a sí mismo como driver de SQL al
	// importarse, pero nunca lo llamamos por su nombre directamente:
	// sqlx lo usa internamente cuando le decimos sqlx.Connect("pgx", ...).
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect abre una conexión a Postgres. sqlx es una capa fina sobre el
// paquete estándar database/sql: agrega métodos como Get/Select que
// mapean filas de la base directo a structs de Go, pero el SQL que
// escribimos sigue siendo SQL plano (no hay ORM generando queries).
func Connect(databaseURL string) (*sqlx.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL no está configurada")
	}

	conn, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
