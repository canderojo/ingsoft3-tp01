// Package db se encarga de abrir la conexión a Postgres.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// El "_" antes del import significa "importalo solo por su efecto
	// secundario". pgx se registra a sí mismo como driver de SQL al
	// importarse, pero nunca lo llamamos por su nombre directamente:
	// sqlx lo usa internamente cuando le decimos sqlx.Open("pgx", ...).
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect abre una conexión a Postgres. sqlx es una capa fina sobre el
// paquete estándar database/sql: agrega métodos como Get/Select que
// mapean filas de la base directo a structs de Go, pero el SQL que
// escribimos sigue siendo SQL plano (no hay ORM generando queries).
//
// Usamos sqlx.Open en vez de sqlx.Connect: Open prepara el pool de
// conexiones pero no se conecta al arrancar. Cada consulta abre la
// conexión cuando la necesita, así que si la base tarda en estar
// disponible (por ejemplo Neon despertando) el backend no se cae al
// iniciar y se recupera solo cuando la base vuelve. Para saber el
// estado real de la base está /health, que hace el Ping.
func Connect(databaseURL string) (*sqlx.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL no está configurada")
	}

	conn, err := sqlx.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
