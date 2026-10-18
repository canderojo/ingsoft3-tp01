import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { configDefaults } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/profesionales': 'http://localhost:8080',
      '/turnos': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
    },
  },
  test: {
    exclude: [...configDefaults.exclude, 'e2e/**'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'lcov', 'json-summary'],
      include: ['src/utils/**', 'src/api/client.js'],
      exclude: ['**/*.test.js'],
      thresholds: { lines: 80, branches: 75 },
    },
  },
})
