import { test, expect } from '@playwright/test'
import { datosUnicos, buscarHorarioLibre, cancelarTurnosDe } from './helpers.js'

// Texto de un botón de horario. Es tolerante porque el formato de la hora
// depende del idioma del navegador: "17:00", "05:00 p. m." o "5:00pm".
const HORARIO = /^\d{1,2}:\d{2}\s*([ap]\.?\s?m\.?)?$/i

// DNI que usó el test en curso, para que afterEach limpie sus turnos.
let dniUsado = null

test.afterEach(async ({ request }) => {
  await cancelarTurnosDe(request, dniUsado)
  dniUsado = null
})

// Llega al formulario de reserva solo con clics, desde la home: nunca se
// hace goto a /profesionales/:id, que nginx manda a la API.
async function abrirFormularioDeReserva(page, profesional, fecha) {
  await page.goto('/')
  await page.getByRole('link', { name: profesional.nombre }).click()

  await page.getByLabel('Elegí una fecha').fill(fecha)

  // El front descarta la respuesta de horarios de una fecha anterior, así que
  // la grilla de hoy ya no pisa a la nueva. Mientras llega la de la fecha
  // elegida no hay ningún horario habilitado: el primero que aparece es suyo.
  await page.getByRole('button', { name: HORARIO, disabled: false }).first().click()
  await expect(page.getByRole('heading', { name: 'Completá tus datos para reservar' })).toBeVisible()
}

async function completarDatos(page, { nombre, dni, email }) {
  await page.getByLabel('Nombre y apellido').fill(nombre)
  await page.getByLabel('DNI', { exact: true }).fill(dni)
  await page.getByLabel('Email', { exact: true }).fill(email)
}

// Verifica el camino feliz completo por pantalla: reservar un turno y
// cancelarlo, contra el front y la API reales del entorno de E2E_BASE_URL.
// Si falla, una paciente no puede reservar o cancelar desde la app.
test('reservar un turno por pantalla y cancelarlo', async ({ page, request }) => {
  const { profesional, fecha } = await buscarHorarioLibre(request)
  const datos = datosUnicos()
  dniUsado = datos.dni

  await abrirFormularioDeReserva(page, profesional, fecha)
  await completarDatos(page, datos)
  await page.getByRole('button', { name: 'Confirmar reserva' }).click()

  await expect(page.getByText('Tu turno se reservó con éxito.')).toBeVisible()
  const titulo = page.getByRole('heading', { name: /^Turno #\d+$/ })
  await expect(titulo).toBeVisible()
  const id = Number((await titulo.textContent()).replace('Turno #', ''))

  await page.getByRole('button', { name: 'Cancelar turno' }).click()
  await expect(page.getByText('Este turno ya está cancelado y no admite cambios.')).toBeVisible()

  const res = await request.get(`/turnos?dni=${datos.dni}`)
  await expect(res).toBeOK()
  const turnos = await res.json()
  expect(turnos).toHaveLength(1)
  expect(turnos[0]).toMatchObject({ id, estado: 'cancelado' })
})

// Verifica que el front frena un DNI de 7 dígitos con un aviso accesible
// (role="alert") y no llega a crear nada, contra el entorno de E2E_BASE_URL.
// Si falla, la validación del DNI dejó de andar o se creó un turno igual.
test('un DNI de 7 dígitos muestra el error y no crea el turno', async ({ page, request }) => {
  const { profesional, fecha } = await buscarHorarioLibre(request)
  const datos = datosUnicos()
  const dniCorto = datos.dni.slice(0, 7)
  dniUsado = dniCorto

  await abrirFormularioDeReserva(page, profesional, fecha)
  await completarDatos(page, { ...datos, dni: dniCorto })
  await page.getByRole('button', { name: 'Confirmar reserva' }).click()

  const alerta = page.getByRole('alert')
  await expect(alerta).toBeVisible()
  await expect(alerta).toHaveText('El DNI debe tener 8 dígitos, sin puntos ni espacios.')

  const res = await request.get(`/turnos?dni=${dniCorto}`)
  expect(res.status()).toBe(200)
  expect(await res.json()).toEqual([])
})

// Verifica que "Mis turnos" encuentra por DNI un turno creado por API y que
// se puede cancelar desde su detalle, contra el entorno de E2E_BASE_URL.
// Si falla, una paciente no puede ver ni cancelar los turnos que ya sacó.
test('Mis turnos encuentra el turno por DNI y lo cancela', async ({ page, request }) => {
  const { profesional, horario } = await buscarHorarioLibre(request)
  const datos = datosUnicos()
  dniUsado = datos.dni

  // Preparación: el turno se crea directo por la API.
  const resCrear = await request.post('/turnos', {
    data: {
      profesional_id: profesional.id,
      fecha_hora_inicio: horario,
      nombre: datos.nombre,
      dni: datos.dni,
      email: datos.email,
      telefono: '',
    },
  })
  expect(resCrear.status()).toBe(201)
  const turno = await resCrear.json()

  await page.goto('/mis-turnos')
  await page.getByRole('radio', { name: 'DNI' }).check()
  await page.getByRole('textbox', { name: 'Dato para buscar tus turnos' }).fill(datos.dni)
  await page.getByRole('button', { name: 'Buscar' }).click()

  await page.getByRole('link', { name: new RegExp(`Turno #${turno.id}\\b`) }).click()
  await expect(page.getByRole('heading', { name: `Turno #${turno.id}`, exact: true })).toBeVisible()

  await page.getByRole('button', { name: 'Cancelar turno' }).click()
  await expect(page.getByText('Este turno ya está cancelado y no admite cambios.')).toBeVisible()
})
