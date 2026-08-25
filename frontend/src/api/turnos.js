import { get, post, patch } from "./client";
import {
  mockCrearTurno,
  mockObtenerTurno,
  mockListarTurnosDePaciente,
  mockCambiarEstadoTurno,
} from "./mockData";

const MOCK = import.meta.env.VITE_USE_MOCK_DATA === "true";

export const ESTADOS = {
  pendiente: "Pendiente",
  confirmado: "Confirmado",
  cancelado: "Cancelado",
  completado: "Completado",
};

// Transiciones permitidas, replicadas del backend (models.Turno.TransicionesPermitidas)
// para poder decidir qué acciones mostrar sin ida y vuelta al servidor.
export const TRANSICIONES_PERMITIDAS = {
  pendiente: ["confirmado", "cancelado"],
  confirmado: ["completado", "cancelado"],
  cancelado: [],
  completado: [],
};

export function crearTurno(datos) {
  if (MOCK) return mockCrearTurno(datos);
  return post("/turnos", datos);
}

export function obtenerTurno(id) {
  if (MOCK) return mockObtenerTurno(id);
  return get(`/turnos/${id}`);
}

export function listarTurnosDePaciente({ dni, email }) {
  if (MOCK) return mockListarTurnosDePaciente();
  const param = dni ? `dni=${encodeURIComponent(dni)}` : `email=${encodeURIComponent(email)}`;
  return get(`/turnos?${param}`);
}

export function cambiarEstadoTurno(id, estado) {
  if (MOCK) return mockCambiarEstadoTurno(id, estado);
  return patch(`/turnos/${id}/estado`, { estado });
}
