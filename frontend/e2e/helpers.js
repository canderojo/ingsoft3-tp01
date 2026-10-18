import { expect } from '@playwright/test'

// Por defecto las rutas son relativas al baseURL del front: nginx (o el
// proxy de Vite) reenvía /profesionales y /turnos a la API. Con `base` se
// le pega directo a la API (por ejemplo 'http://localhost:8080').

let contador = 0

// Datos de paciente que no se repiten entre corridas. El DNI es la clave
// con la que el backend reutiliza pacientes, y el email es UNIQUE en la
// base: un email repetido con un DNI nuevo hace fallar el alta con 500.
export function datosUnicos() {
  const ahora = Date.now()
  contador += 1
  // 10.000.000 a 99.999.999: siempre 8 dígitos y nunca empieza con 0.
  // Cambia cada milisegundo y el contador separa dos llamadas en el mismo ms.
  const dni = String(10_000_000 + ((ahora + contador) % 90_000_000))
  return {
    nombre: `E2E Paciente ${ahora}`,
    dni,
    email: `e2e.${ahora}.${contador}@example.com`,
  }
}

// Fecha YYYY-MM-DD de dentro de `dias` días según el reloj de Argentina
// (UTC-3), que es el que usa el backend para decidir qué horario ya pasó.
function fechaEnDias(dias) {
  const hoyArgentina = new Date(Date.now() - 3 * 60 * 60 * 1000)
  hoyArgentina.setUTCDate(hoyArgentina.getUTCDate() + dias)
  return hoyArgentina.toISOString().slice(0, 10)
}

// Primer horario libre del primer profesional, buscando desde mañana hasta
// dentro de 14 días. Si no hay ninguno, el test falla: no se saltea.
export async function buscarHorarioLibre(request, base = '') {
  const resProfesionales = await request.get(`${base}/profesionales`)
  await expect(resProfesionales).toBeOK()
  const profesionales = await resProfesionales.json()
  if (profesionales.length === 0) {
    throw new Error('GET /profesionales devolvió una lista vacía: no hay con quién reservar')
  }
  const profesional = profesionales[0]

  for (let dias = 1; dias <= 14; dias++) {
    const fecha = fechaEnDias(dias)
    const res = await request.get(`${base}/profesionales/${profesional.id}/horarios-disponibles?fecha=${fecha}`)
    await expect(res).toBeOK()
    const horarios = await res.json()
    if (horarios.length > 0) {
      return { profesional, fecha, horario: horarios[0] }
    }
  }

  throw new Error(
    `${profesional.nombre} (id ${profesional.id}) no tiene ningún horario libre entre mañana y dentro de 14 días`
  )
}

// Red de seguridad: cancela los turnos de ese DNI que todavía se pueden
// cancelar. Solo pendiente y confirmado admiten pasar a cancelado; un
// completado devolvería 422 "transición de estado no permitida".
export async function cancelarTurnosDe(request, dni, base = '') {
  if (!dni) return
  const res = await request.get(`${base}/turnos?dni=${dni}`)
  await expect(res).toBeOK()
  const turnos = await res.json()
  for (const turno of turnos) {
    if (turno.estado !== 'pendiente' && turno.estado !== 'confirmado') continue
    const resCancelar = await request.patch(`${base}/turnos/${turno.id}/estado`, {
      data: { estado: 'cancelado' },
    })
    await expect(resCancelar).toBeOK()
  }
}
