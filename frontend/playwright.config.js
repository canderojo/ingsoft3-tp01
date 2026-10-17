import { defineConfig } from '@playwright/test'

export default defineConfig({
  // Los specs de Playwright viven aparte de los tests de vitest (src/**).
  testDir: './e2e',

  // 60 s por test: en el free tier de Render el primer pedido despierta al
  // servicio dormido (cold start) y puede tardar bastante en contestar.
  timeout: 60_000,

  // Mismo motivo: cada expect espera hasta 15 s a que aparezca lo que busca.
  expect: { timeout: 15_000 },

  // Un solo reintento: absorbe un cold start lento sin tapar un fallo real.
  retries: 1,

  // El reporte HTML se genera siempre (sin abrirse solo) y 'list' muestra
  // el avance en la consola.
  reporter: [['html', { open: 'never' }], ['list']],

  use: {
    // Por defecto, el front del compose local; en CI se apunta a QA
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',

    // La traza solo se graba en el reintento: pesa poco y explica el fallo.
    trace: 'on-first-retry',

    // Captura de pantalla solo cuando un test falla.
    screenshot: 'only-on-failure',
  },
})
