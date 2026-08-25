import { get } from "./client";
import { mockListarProfesionales, mockObtenerProfesional, mockHorariosDisponibles } from "./mockData";

const MOCK = import.meta.env.VITE_USE_MOCK_DATA === "true";

export const ESPECIALIDADES = [
  { value: "dermatologia", label: "Dermatología" },
  { value: "nutricion", label: "Nutrición" },
  { value: "ecografia", label: "Ecografía" },
  { value: "endocrinologia", label: "Endocrinología" },
];

export function listarProfesionales(especialidad) {
  if (MOCK) return mockListarProfesionales(especialidad);
  const query = especialidad ? `?especialidad=${especialidad}` : "";
  return get(`/profesionales${query}`);
}

export function obtenerProfesional(id) {
  if (MOCK) return mockObtenerProfesional(id);
  return get(`/profesionales/${id}`);
}

export function horariosDisponibles(id, fecha) {
  if (MOCK) return mockHorariosDisponibles(id, fecha);
  return get(`/profesionales/${id}/horarios-disponibles?fecha=${fecha}`);
}
