// Arma el texto de resumen que se puede mostrar arriba de "Mis turnos":
// cuántos turnos activos (pendientes o confirmados) tiene la paciente y
// cuántos de ellos todavía no confirmó.
export function textoResumenTurnos(turnos) {
  if (!turnos || turnos.length === 0) {
    return "Todavía no tenés turnos reservados.";
  }

  const activos = turnos.filter((t) => t.estado === "pendiente" || t.estado === "confirmado").length;
  const sinConfirmar = turnos.filter((t) => t.estado === "pendiente").length;

  if (activos === 0) {
    return "No tenés turnos activos.";
  }

  let texto = activos === 1 ? "Tenés 1 turno activo" : `Tenés ${activos} turnos activos`;
  if (sinConfirmar > 0) {
    texto += sinConfirmar === 1 ? " (1 sin confirmar)" : ` (${sinConfirmar} sin confirmar)`;
  }
  return texto + ".";
}
