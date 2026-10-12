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
// El código HTTP refleja el estado de la base: 200 si el Ping anda, y
// 503 (Service Unavailable) si la base dio error o no está configurada.
// Así un monitor o el orquestador se entera del problema mirando solo
// el código, sin tener que leer el JSON.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"

	if h.DB == nil {
		dbStatus = "sin configurar"
	} else if err := h.DB.Ping(); err != nil {
		dbStatus = "error"
	}

	code := http.StatusOK
	status := "ok"
	if dbStatus != "ok" {
		code = http.StatusServiceUnavailable
		status = "degradado"
	}

	// El orden importa: primero los headers, después WriteHeader con el
	// código y recién al final el body. Una vez que se escribe el body,
	// Go ya mandó los headers y no se pueden cambiar.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
		"db":     dbStatus,
	})
}
