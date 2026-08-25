package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/repository"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/service"
)

// ProfesionalesHandler agrupa los endpoints de solo lectura sobre
// profesionales (no hay alta/baja/modificación por API: se cargan a
// mano en la base, ver db/init.sql).
type ProfesionalesHandler struct {
	DB *sqlx.DB
}

// Listar responde GET /profesionales (opcionalmente ?especialidad=...).
func (h *ProfesionalesHandler) Listar(w http.ResponseWriter, r *http.Request) {
	especialidad := r.URL.Query().Get("especialidad")

	profesionales, err := repository.ListarProfesionales(h.DB, especialidad)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "no se pudo listar profesionales")
		return
	}

	responderJSON(w, http.StatusOK, profesionales)
}

// Obtener responde GET /profesionales/{id}.
func (h *ProfesionalesHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	profesional, err := repository.ObtenerProfesional(h.DB, id)
	if err == sql.ErrNoRows {
		responderError(w, http.StatusNotFound, "profesional no encontrado")
		return
	}
	if err != nil {
		responderError(w, http.StatusInternalServerError, "no se pudo obtener el profesional")
		return
	}

	responderJSON(w, http.StatusOK, profesional)
}

// formatoFecha es el formato esperado en ?fecha= (sin hora, ej: "2026-08-20").
const formatoFecha = "2006-01-02"

// HorariosDisponibles responde GET /profesionales/{id}/horarios-disponibles?fecha=YYYY-MM-DD:
// la lista de horarios en los que el paciente puede reservar un turno
// con ese profesional ese día.
func (h *ProfesionalesHandler) HorariosDisponibles(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	fechaStr := r.URL.Query().Get("fecha")
	if fechaStr == "" {
		responderError(w, http.StatusBadRequest, "hay que pasar ?fecha=YYYY-MM-DD")
		return
	}

	fecha, err := time.Parse(formatoFecha, fechaStr)
	if err != nil {
		responderError(w, http.StatusBadRequest, "fecha inválida, usar formato YYYY-MM-DD")
		return
	}

	horarios, err := service.HorariosDisponibles(h.DB, id, fecha)
	if err != nil {
		responderErrorDeNegocio(w, err)
		return
	}

	responderJSON(w, http.StatusOK, horarios)
}
