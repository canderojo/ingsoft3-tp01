// Validaciones del formulario de reserva, para avisarle a la paciente antes
// de mandar el pedido al backend. Cada función devuelve el mensaje de error,
// o "" si el dato está bien.

export function validarDNI(dni) {
  const limpio = (dni ?? "").trim();
  if (limpio === "") {
    return "Ingresá tu DNI.";
  }
  if (!/^\d+$/.test(limpio)) {
    return "El DNI tiene que tener sólo números, sin puntos.";
  }
  if (limpio.length < 7 || limpio.length > 8) {
    return "El DNI tiene que tener 7 u 8 números.";
  }
  return "";
}

export function validarEmail(email) {
  const limpio = (email ?? "").trim();
  if (limpio === "") {
    return "Ingresá tu email.";
  }
  if (!limpio.includes("@") || !limpio.includes(".")) {
    return "El email no parece válido.";
  }
  return "";
}
