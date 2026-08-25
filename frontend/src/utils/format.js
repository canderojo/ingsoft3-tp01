export function formatHora(horaHHMMSS) {
  return horaHHMMSS?.slice(0, 5) ?? "";
}

export function formatPrecio(precio) {
  return new Intl.NumberFormat("es-AR", {
    style: "currency",
    currency: "ARS",
    maximumFractionDigits: 0,
  }).format(precio);
}

// El backend manda las horas de atención en reloj de Argentina pero
// serializadas con sufijo "Z" (UTC) sin conversión real (no maneja huso
// horario). Si dejáramos que el navegador las interprete como UTC de
// verdad, se corren -3hs al mostrarlas. Por eso en todo este archivo
// forzamos timeZone: "UTC" (para mostrar) o usamos los getters UTC* (para
// comparar): así se toman los dígitos tal cual vienen, sin conversión.

export function formatFechaHora(iso) {
  const fecha = new Date(iso);
  return fecha.toLocaleString("es-AR", {
    weekday: "long",
    day: "numeric",
    month: "long",
    hour: "2-digit",
    minute: "2-digit",
    timeZone: "UTC",
  });
}

export function formatHoraSlot(iso) {
  return new Date(iso).toLocaleTimeString("es-AR", {
    hour: "2-digit",
    minute: "2-digit",
    timeZone: "UTC",
  });
}

// Clave de comparación por hora de reloj (ignora cómo cada Date serializa
// el huso horario), para saber si dos timestamps son "el mismo horario".
export function claveSlot(iso) {
  const d = new Date(iso);
  return `${d.getUTCFullYear()}-${d.getUTCMonth()}-${d.getUTCDate()}-${d.getUTCHours()}-${d.getUTCMinutes()}`;
}

export function hoyISO() {
  const hoy = new Date();
  const offset = hoy.getTimezoneOffset();
  const local = new Date(hoy.getTime() - offset * 60 * 1000);
  return local.toISOString().slice(0, 10);
}
