import { describe, it, expect } from "vitest";
import { textoResumenTurnos } from "./resumen";

describe("textoResumenTurnos", () => {
  // Sin doble: textoResumenTurnos es una función pura. Es un test
  // parametrizado: it.each corre el mismo test con cada fila de la tabla, y
  // hay una fila por cada camino de la función (sin lista, sin activos, uno
  // solo, varios con y sin confirmar).
  it.each([
    ["sin turnos (lista vacía)", [], "Todavía no tenés turnos reservados."],
    ["sin turnos (no vino la lista)", null, "Todavía no tenés turnos reservados."],
    ["sólo turnos cancelados o completados", [{ estado: "cancelado" }, { estado: "completado" }], "No tenés turnos activos."],
    ["un turno confirmado", [{ estado: "confirmado" }], "Tenés 1 turno activo."],
    ["dos activos, uno sin confirmar", [{ estado: "confirmado" }, { estado: "pendiente" }], "Tenés 2 turnos activos (1 sin confirmar)."],
    ["tres pendientes", [{ estado: "pendiente" }, { estado: "pendiente" }, { estado: "pendiente" }], "Tenés 3 turnos activos (3 sin confirmar)."],
  ])("%s", (_caso, turnos, esperado) => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado. El Arrange son los datos de la tabla.
    expect(textoResumenTurnos(turnos)).toBe(esperado);
  });
});
