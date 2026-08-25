// Datos de prueba SOLO para desarrollo del frontend, mientras el backend/DB
// real no está levantado. Se activan con VITE_USE_MOCK_DATA=true en .env.
// Para sacarlos: borrar este archivo y las referencias a MOCK en
// api/profesionales.js y api/turnos.js.

import { ApiError } from "./client";
import { claveSlot } from "../utils/format";

function delay(data) {
  return new Promise((resolve) => setTimeout(() => resolve(data), 300));
}

const PROFESIONALES = [
  {
    id: 1,
    nombre: "Dra. Milagros Ferreyra",
    especialidad: "dermatologia",
    hora_inicio_atencion: "09:00:00",
    hora_fin_atencion: "14:00:00",
    duracion_turno_minutos: 30,
    precio_consulta: 18000,
  },
  {
    id: 2,
    nombre: "Lic. Bianca Torres",
    especialidad: "nutricion",
    hora_inicio_atencion: "10:00:00",
    hora_fin_atencion: "18:00:00",
    duracion_turno_minutos: 40,
    precio_consulta: 15500,
  },
  {
    id: 3,
    nombre: "Dra. Yamila Correa",
    especialidad: "ecografia",
    hora_inicio_atencion: "08:00:00",
    hora_fin_atencion: "13:00:00",
    duracion_turno_minutos: 30,
    precio_consulta: 22000,
  },
  {
    id: 4,
    nombre: "Dr. Ezequiel Paz",
    especialidad: "endocrinologia",
    hora_inicio_atencion: "14:00:00",
    hora_fin_atencion: "19:00:00",
    duracion_turno_minutos: 30,
    precio_consulta: 19800,
  },
];

// Turnos ya reservados, para que "Mis turnos" tenga algo que encontrar
// sin necesidad de reservar uno primero, y para que se vean horarios
// ocupados en la grilla de un profesional. DNI "12345678" / email "paciente@demo.com".
let turnos = [
  {
    id: 1001,
    profesional_id: 2,
    paciente_id: 1,
    fecha_hora_inicio: proximaFecha(3, 10, 0),
    fecha_hora_fin: proximaFecha(3, 10, 40),
    estado: "pendiente",
    precio: 15500,
  },
  {
    id: 1002,
    profesional_id: 1,
    paciente_id: 1,
    fecha_hora_inicio: proximaFecha(0, 9, 30),
    fecha_hora_fin: proximaFecha(0, 10, 0),
    estado: "confirmado",
    precio: 18000,
  },
];
let siguienteId = 1003;

function proximaFecha(diasDesdeHoy, hora, minuto) {
  const fecha = new Date();
  fecha.setDate(fecha.getDate() + diasDesdeHoy);
  fecha.setHours(hora, minuto, 0, 0);
  return fecha.toISOString();
}

function slotsOcupados(profesionalId) {
  return new Set(
    turnos
      .filter((t) => t.profesional_id === profesionalId && t.estado !== "cancelado")
      .map((t) => claveSlot(t.fecha_hora_inicio))
  );
}

export function mockListarProfesionales(especialidad) {
  const resultado = especialidad
    ? PROFESIONALES.filter((p) => p.especialidad === especialidad)
    : PROFESIONALES;
  return delay(resultado);
}

export function mockObtenerProfesional(id) {
  const profesional = PROFESIONALES.find((p) => p.id === Number(id));
  return profesional
    ? delay(profesional)
    : Promise.reject(new ApiError("Profesional no encontrado", 404));
}

export function mockHorariosDisponibles(id, fecha) {
  const profesional = PROFESIONALES.find((p) => p.id === Number(id));
  if (!profesional) return Promise.reject(new ApiError("Profesional no encontrado", 404));

  const [horaInicio, minInicio] = profesional.hora_inicio_atencion.split(":").map(Number);
  const [horaFin, minFin] = profesional.hora_fin_atencion.split(":").map(Number);
  const duracion = profesional.duracion_turno_minutos;

  const inicio = new Date(`${fecha}T00:00:00`);
  inicio.setHours(horaInicio, minInicio, 0, 0);
  const fin = new Date(`${fecha}T00:00:00`);
  fin.setHours(horaFin, minFin, 0, 0);

  const ocupados = slotsOcupados(profesional.id);
  const slots = [];
  for (let t = new Date(inicio); t < fin; t.setMinutes(t.getMinutes() + duracion)) {
    const iso = new Date(t).toISOString();
    if (!ocupados.has(claveSlot(iso))) slots.push(iso);
  }
  return delay(slots);
}

export function mockCrearTurno(datos) {
  const profesional = PROFESIONALES.find((p) => p.id === Number(datos.profesional_id));
  if (!profesional) return Promise.reject(new ApiError("Profesional no encontrado", 404));

  if (slotsOcupados(profesional.id).has(claveSlot(datos.fecha_hora_inicio))) {
    return Promise.reject(
      new ApiError("Ese horario ya no está disponible. Elegí otro.", 409)
    );
  }

  const inicio = new Date(datos.fecha_hora_inicio);
  const fin = new Date(inicio.getTime() + profesional.duracion_turno_minutos * 60000);

  const turno = {
    id: siguienteId++,
    profesional_id: profesional.id,
    paciente_id: 1,
    fecha_hora_inicio: inicio.toISOString(),
    fecha_hora_fin: fin.toISOString(),
    estado: "pendiente",
    precio: profesional.precio_consulta,
  };
  turnos = [...turnos, turno];
  return delay(turno);
}

export function mockObtenerTurno(id) {
  const turno = turnos.find((t) => t.id === Number(id));
  return turno ? delay(turno) : Promise.reject(new ApiError("Turno no encontrado", 404));
}

export function mockListarTurnosDePaciente() {
  // en el mock, cualquier búsqueda devuelve los turnos del paciente de prueba
  return delay(turnos);
}

export function mockCambiarEstadoTurno(id, estado) {
  const turno = turnos.find((t) => t.id === Number(id));
  if (!turno) return Promise.reject(new ApiError("Turno no encontrado", 404));
  turno.estado = estado;
  return delay({ ...turno });
}
