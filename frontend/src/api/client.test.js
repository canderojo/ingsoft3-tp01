import { describe, it, expect, vi, afterEach } from "vitest";
import { get, post, ApiError } from "./client";

// fetch es la frontera con la red: en estos tests lo reemplazamos por un
// doble (vi.fn) que devuelve la respuesta que cada test necesita.

// respuesta arma lo que devuelve el fetch de mentira: si salió bien (ok), el
// código HTTP (status) y el body. Con jsonFalla simula un body que no se puede
// leer.
function respuesta({ ok, status, body, jsonFalla = false }) {
  return {
    ok,
    status,
    json: jsonFalla ? () => Promise.reject(new Error("no es JSON")) : () => Promise.resolve(body),
  };
}

// Después de cada test se saca el fetch de mentira, para que no pase al
// siguiente.
afterEach(() => {
  vi.unstubAllGlobals();
});

describe("client", () => {
  // Stub y mock. Stub: el fetch de mentira contesta un 422 con un mensaje.
  // Mock: las dos últimas líneas no miran el resultado, revisan cómo se usó
  // fetch (una sola vez, con POST, a una URL que termina en /turnos). Es el
  // equivalente al Verify de la guía. Es también el caso de error: el mensaje
  // del backend tiene que llegar en el ApiError, porque es lo que ve la
  // paciente en pantalla.
  it("un error del backend se convierte en ApiError con su mensaje y su status", async () => {
    // Arrange: fetch contesta un 422 con el mensaje del backend.
    const fetchDoble = vi.fn().mockResolvedValue(
      respuesta({ ok: false, status: 422, body: { error: "el turno está fuera del horario de atención" } })
    );
    vi.stubGlobal("fetch", fetchDoble);

    // Act: se pide crear un turno.
    const promesa = post("/turnos", { profesional_id: 1 });

    // Assert: el error que llega y cómo se usó fetch.
    // rejects: se espera que la promesa falle. Se revisa con qué error falla.
    await expect(promesa).rejects.toBeInstanceOf(ApiError);
    await expect(promesa).rejects.toMatchObject({
      message: "el turno está fuera del horario de atención",
      status: 422,
    });
    expect(fetchDoble).toHaveBeenCalledTimes(1);
    // La URL termina en /turnos, sin importar qué VITE_API_URL tenga cada entorno
    // (en local puede haber un .env que la defina; en el pipeline no).
    expect(fetchDoble).toHaveBeenCalledWith(
      expect.stringMatching(/\/turnos$/),
      expect.objectContaining({ method: "POST" })
    );
  });

  // Stub: el fetch de mentira contesta un 500 con un body que no se puede
  // leer. Caso de error: se usa el mensaje genérico con el código.
  it("si el error no trae un body legible, usa un mensaje genérico con el status", async () => {
    // Arrange: fetch contesta un 500 con un body ilegible.
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(respuesta({ ok: false, status: 500, jsonFalla: true })));

    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    await expect(get("/turnos/1")).rejects.toMatchObject({ message: "Error 500", status: 500 });
  });

  // Stub: el fetch de mentira contesta 204 (sin contenido). resolves: se
  // espera que la promesa salga bien, y se revisa con qué valor.
  it("una respuesta 204 (sin contenido) devuelve null", async () => {
    // Arrange: fetch contesta 204 sin body.
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(respuesta({ ok: true, status: 204 })));

    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    await expect(get("/algo")).resolves.toBeNull();
  });

  // Stub: el fetch de mentira contesta 200 con un body. Es el camino feliz.
  it("una respuesta OK devuelve el JSON del body", async () => {
    // Arrange: fetch contesta 200 con un body.
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(respuesta({ ok: true, status: 200, body: { id: 7 } })));

    // Act y Assert en la misma línea: se llama a la función y se compara el
    // resultado.
    await expect(get("/turnos/7")).resolves.toEqual({ id: 7 });
  });
});
