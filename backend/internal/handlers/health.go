// Package handlers contiene las funciones que atienden requests HTTP.
// Cada handler recibe (w http.ResponseWriter, r *http.Request): "w" es
// donde escribimos la respuesta, "r" es el request que llegó.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// HealthHandler agrupa las dependencias que necesita el endpoint /health.
// En Go, un "puntero" (*sqlx.DB) es una referencia a un valor en memoria
// en vez de una copia — así todos los handlers comparten la misma
// conexión a la base en vez de abrir una nueva cada vez.
type HealthHandler struct {
	DB *sqlx.DB
}

// Health responde si la API está viva y si puede hablar con Postgres.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"

	if h.DB == nil {
		dbStatus = "sin configurar"
	} else if err := h.DB.Ping(); err != nil {
		dbStatus = "error"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     dbStatus,
	})
}
