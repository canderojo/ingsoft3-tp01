import { describe, it, expect, vi, afterEach } from "vitest";
import { formatHora, formatPrecio, formatHoraSlot, formatFechaHora, claveSlot, hoyISO } from "./format";

// Estructura de vitest: describe agrupa los tests de una función, it es un test,
// y expect(...).toBe(...) es el assert. Es el mismo patrón AAA que los tests del
// backend (preparar, ejecutar, comprobar), con otros nombres.

// Intl separa algunas partes con espacios especiales (no separables):
// los pasamos a espacios comunes para poder comparar el texto.
const normalizar = (texto) => texto.replace(/\s/g, " ");

describe("formatHora", () => {
  // Sin doble: formatHora es una función pura. Es un test parametrizado:
  // it.each corre el mismo test con cada fila de la tabla (equivale a la
  // tabla con t.Run del backend). Los casos undefined y null comprueban que,
  // si la hora no viene, devuelva vacío en vez de romperse.
  it.each([
    ["09:30:00", "09:30"],
    ["20:00:00", "20:00"],
    [undefined, ""],
    [null, ""],
  ])("formatHora(%s) devuelve %j", (entrada, esperado) => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado. El Arrange son los datos de la tabla.
    expect(formatHora(entrada)).toBe(esperado);
  });
});

describe("formatPrecio", () => {
  // Sin doble: función pura. normalizar cambia los espacios especiales que
  // pone Intl por espacios comunes, para poder comparar el texto.
  it("muestra el precio en pesos, con punto de miles y sin decimales", () => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(normalizar(formatPrecio(5800))).toBe("$ 5.800");
  });
});

describe("horarios que vienen del backend", () => {
  // El backend manda la hora de Argentina con sufijo Z: la app tiene que
  // mostrar esos dígitos tal cual, sin restarle 3 horas.

  // Sin doble: funciones puras. Protegen la regla de la zona horaria del lado
  // del front: si alguien sacara el timeZone: "UTC" del código, la hora se
  // mostraría 3 horas antes y estos tests se pondrían en rojo.
  it("formatHoraSlot muestra la hora tal cual viene, sin correrla 3 horas", () => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(normalizar(formatHoraSlot("2026-10-16T19:30:00Z"))).toBe("07:30 p. m.");
  });

  // Sin doble: funciones puras. Protegen la regla de la zona horaria del lado
  // del front: si alguien sacara el timeZone: "UTC" del código, la hora se
  // mostraría 3 horas antes y estos tests se pondrían en rojo.
  it("formatFechaHora muestra el día y la hora tal cual vienen", () => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(normalizar(formatFechaHora("2026-10-16T10:00:00Z"))).toBe("viernes, 16 de octubre, 10:00 a. m.");
  });
});

describe("claveSlot", () => {
  // Sin doble: función pura. El front usa esta clave para saber qué horarios
  // tachar.
  it("el mismo horario escrito de dos formas da la misma clave", () => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(claveSlot("2026-10-16T10:00:00Z")).toBe(claveSlot("2026-10-16T10:00:00.000Z"));
  });

  // Sin doble: función pura. El front usa esta clave para saber qué horarios
  // tachar.
  it("dos horarios distintos dan claves distintas", () => {
    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(claveSlot("2026-10-16T10:00:00Z")).not.toBe(claveSlot("2026-10-16T10:30:00Z"));
  });
});

describe("hoyISO", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  // Stub del reloj: vi.useFakeTimers y vi.setSystemTime reemplazan la hora
  // real por una fija (como relojFijo en el backend). afterEach vuelve al
  // reloj real para no afectar a los otros tests.
  it("devuelve la fecha de hoy en formato AAAA-MM-DD", () => {
    // Arrange: se fija el reloj en una fecha y hora conocidas.
    // Reloj fijo al mediodía: da el mismo día en Argentina y en UTC (el
    // pipeline corre en UTC), así el test no depende de dónde se corre.
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-10-16T15:00:00Z"));

    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    expect(hoyISO()).toBe("2026-10-16");
  });
});
