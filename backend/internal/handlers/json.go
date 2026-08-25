package handlers

import (
	"encoding/json"
	"net/http"
)

// responderJSON escribe cualquier valor como respuesta JSON con el
// código de estado indicado. Está centralizado acá para que todos los
// handlers respondan de forma consistente.
func responderJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

// responderError es un atajo para responderJSON con el formato
// {"error": "..."} que usan todos los handlers ante un fallo.
func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
}
