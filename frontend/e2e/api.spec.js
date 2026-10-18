import { test, expect } from '@playwright/test'
import { datosUnicos, buscarHorarioLibre, cancelarTurnosDe } from './helpers.js'

// Pruebas de integración contra la API, sin navegador: se le pega directo
// al backend, no al front. Sin barra final, porque las rutas empiezan con '/'.
const API = (process.env.API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '')

// DNIs que usó el test en curso, para que afterEach limpie sus turnos.
let dnisUsados = []

test.afterEach(async ({ request }) => {
  for (const dni of dnisUsados) {
    await cancelarTurnosDe(request, dni, API)
  }
  dnisUsados = []
})

// Body de POST /turnos para un paciente y un horario.
function altaTurno(profesional, horario, { nombre, dni, email }) {
  return {
    profesional_id: profesional.id,
    fecha_hora_inicio: horario,
    nombre,
    dni,
    email,
    telefono: '',
  }
}

// true si `horario` está en la lista de horarios disponibles del profesional
// para esa fecha. Compara instantes, no strings: "…Z" y "….000Z" son lo mismo.
async function horarioDisponible(request, profesional, fecha, horario) {
  const res = await request.get(`${API}/profesionales/${profesional.id}/horarios-disponibles?fecha=${fecha}`)
  await expect(res).toBeOK()
  const buscado = new Date(horario).getTime()
  return (await res.json()).some((h) => new Date(h).getTime() === buscado)
}

// Verifica contra la API (API_BASE_URL) que reservar guarda el turno y ocupa
// el horario, y que cancelarlo lo libera. Si falla, la reserva o la
// cancelación dejaron de persistir, o la grilla de disponibles miente.
test('reservar guarda el turno y cancelar libera el horario', async ({ request }) => {
  const { profesional, fecha, horario } = await buscarHorarioLibre(request, API)
  const datos = datosUnicos()
  dnisUsados.push(datos.dni)

  const resCrear = await request.post(`${API}/turnos`, { data: altaTurno(profesional, horario, datos) })
  expect(resCrear.status()).toBe(201)
  const turno = await resCrear.json()
  expect(turno.estado).toBe('pendiente')

  const resLista = await request.get(`${API}/turnos?dni=${datos.dni}`)
  await expect(resLista).toBeOK()
  expect((await resLista.json()).map((t) => t.id)).toContain(turno.id)

  const resObtener = await request.get(`${API}/turnos/${turno.id}`)
  expect(resObtener.status()).toBe(200)

  expect(await horarioDisponible(request, profesional, fecha, horario)).toBe(false)

  const resCancelar = await request.patch(`${API}/turnos/${turno.id}/estado`, {
    data: { estado: 'cancelado' },
  })
  expect(resCancelar.status()).toBe(200)

  const resDespues = await request.get(`${API}/turnos/${turno.id}`)
  await expect(resDespues).toBeOK()
  expect((await resDespues.json()).estado).toBe('cancelado')

  expect(await horarioDisponible(request, profesional, fecha, horario)).toBe(true)
})

// Verifica contra la API (API_BASE_URL) que un alta sin nombre se rechaza
// con 400 y no deja nada guardado. Si falla, el backend acepta altas
// incompletas o crea al paciente aunque rechace el turno.
test('un alta sin nombre da 400 y no crea nada', async ({ request }) => {
  const { profesional, horario } = await buscarHorarioLibre(request, API)
  const datos = datosUnicos()
  dnisUsados.push(datos.dni)

  const res = await request.post(`${API}/turnos`, {
    data: altaTurno(profesional, horario, { ...datos, nombre: '' }),
  })
  expect(res.status()).toBe(400)
  expect((await res.json()).error).toContain('faltan campos obligatorios')

  const resLista = await request.get(`${API}/turnos?dni=${datos.dni}`)
  expect(resLista.status()).toBe(200)
  expect(await resLista.json()).toEqual([])
})

// Verifica contra la API (API_BASE_URL) que un horario tomado no se puede
// reservar de nuevo (409) y que al cancelarlo otra paciente sí puede. Si
// falla, hay turnos superpuestos o un turno cancelado sigue bloqueando.
test('un horario ocupado da 409 y se libera al cancelar', async ({ request }) => {
  const { profesional, horario } = await buscarHorarioLibre(request, API)
  const pacienteA = datosUnicos()
  const pacienteB = datosUnicos()
  dnisUsados.push(pacienteA.dni, pacienteB.dni)

  const resA = await request.post(`${API}/turnos`, { data: altaTurno(profesional, horario, pacienteA) })
  expect(resA.status()).toBe(201)
  const turnoA = await resA.json()

  const resB = await request.post(`${API}/turnos`, { data: altaTurno(profesional, horario, pacienteB) })
  expect(resB.status()).toBe(409)
  expect((await resB.json()).error).toContain('ya tiene un turno en ese horario')

  const resCancelarA = await request.patch(`${API}/turnos/${turnoA.id}/estado`, {
    data: { estado: 'cancelado' },
  })
  expect(resCancelarA.status()).toBe(200)

  const resB2 = await request.post(`${API}/turnos`, { data: altaTurno(profesional, horario, pacienteB) })
  expect(resB2.status()).toBe(201)
  const turnoB = await resB2.json()

  const resCancelarB = await request.patch(`${API}/turnos/${turnoB.id}/estado`, {
    data: { estado: 'cancelado' },
  })
  expect(resCancelarB.status()).toBe(200)
})
