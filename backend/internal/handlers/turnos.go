package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/repository"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/service"
)

// TurnosHandler agrupa los endpoints de reserva y consulta de turnos.
type TurnosHandler struct {
	DB *sqlx.DB
}

// crearTurnoRequest es el body esperado de POST /turnos. El paciente
// no necesita existir de antemano: se identifica (o se crea) por DNI.
type crearTurnoRequest struct {
	ProfesionalID   int    `json:"profesional_id"`
	FechaHoraInicio string `json:"fecha_hora_inicio"` // formato RFC3339, ej: "2026-08-20T10:00:00Z"
	Nombre          string `json:"nombre"`
	DNI             string `json:"dni"`
	Email           string `json:"email"`
	Telefono        string `json:"telefono"`
}

// Crear responde POST /turnos: reserva un turno nuevo.
func (h *TurnosHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var body crearTurnoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "body inválido")
		return
	}

	if body.ProfesionalID == 0 || body.Nombre == "" || body.DNI == "" || body.Email == "" {
		responderError(w, http.StatusBadRequest, "faltan campos obligatorios (profesional_id, nombre, dni, email)")
		return
	}

	inicio, err := time.Parse(time.RFC3339, body.FechaHoraInicio)
	if err != nil {
		responderError(w, http.StatusBadRequest, "fecha_hora_inicio inválida, usar formato RFC3339 (ej: 2026-08-20T10:00:00Z)")
		return
	}

	turno, err := service.CrearTurno(h.DB, service.CrearTurnoInput{
		ProfesionalID:   body.ProfesionalID,
		FechaHoraInicio: inicio,
		Nombre:          body.Nombre,
		DNI:             body.DNI,
		Email:           body.Email,
		Telefono:        body.Telefono,
	})
	if err != nil {
		responderErrorDeNegocio(w, err)
		return
	}

	responderJSON(w, http.StatusCreated, turno)
}

// Obtener responde GET /turnos/{id}.
func (h *TurnosHandler) Obtener(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	turno, err := service.ObtenerTurno(h.DB, id)
	if err != nil {
		responderErrorDeNegocio(w, err)
		return
	}

	responderJSON(w, http.StatusOK, turno)
}

// ListarDePaciente responde GET /turnos?dni=...  o  ?email=...
// ("mis turnos", sin necesidad de login).
func (h *TurnosHandler) ListarDePaciente(w http.ResponseWriter, r *http.Request) {
	dni := r.URL.Query().Get("dni")
	email := r.URL.Query().Get("email")

	if dni == "" && email == "" {
		responderError(w, http.StatusBadRequest, "hay que pasar ?dni= o ?email=")
		return
	}

	paciente, err := repository.BuscarPacientePorDNIoEmail(h.DB, dni, email)
	if err == sql.ErrNoRows {
		responderJSON(w, http.StatusOK, []interface{}{})
		return
	}
	if err != nil {
		responderError(w, http.StatusInternalServerError, "no se pudo buscar el paciente")
		return
	}

	turnos, err := service.ListarTurnosDePaciente(h.DB, paciente.ID)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "no se pudieron listar los turnos")
		return
	}

	responderJSON(w, http.StatusOK, turnos)
}

// cambiarEstadoRequest es el body esperado de PATCH /turnos/{id}/estado.
type cambiarEstadoRequest struct {
	Estado string `json:"estado"`
}

// CambiarEstado responde PATCH /turnos/{id}/estado: confirma, cancela
// o completa un turno, respetando la máquina de estados definida en
// models.TransicionesPermitidas.
func (h *TurnosHandler) CambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var body cambiarEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "body inválido")
		return
	}

	turno, err := service.CambiarEstadoTurno(h.DB, id, body.Estado)
	if err != nil {
		responderErrorDeNegocio(w, err)
		return
	}

	responderJSON(w, http.StatusOK, turno)
}

// responderErrorDeNegocio traduce errores conocidos de la capa de
// negocio (service) al código HTTP correspondiente. Cualquier otro
// error (por ejemplo, uno de conexión a la base) se trata como 500.
func responderErrorDeNegocio(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrProfesionalNoExiste), errors.Is(err, service.ErrTurnoNoExiste):
		responderError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrFueraDeHorario),
		errors.Is(err, service.ErrFechaEnElPasado),
		errors.Is(err, service.ErrTransicionInvalida):
		responderError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, service.ErrSuperposicionProfesional),
		errors.Is(err, service.ErrSuperposicionPaciente):
		responderError(w, http.StatusConflict, err.Error())
	default:
		responderError(w, http.StatusInternalServerError, "error interno")
	}
}
