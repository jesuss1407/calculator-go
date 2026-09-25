import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Forward API calls to the Go backend so the browser sees one origin and no CORS is needed.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  test: {
    // Only unit and component tests; the Playwright specs in e2e/ run with `npm run e2e`.
    include: ['src/**/*.test.{ts,tsx}'],
    environment: 'jsdom',
    setupFiles: './src/setupTests.ts',
    coverage: {
      provider: 'v8',
      // skipFull: false lists fully covered files in the terminal table too.
      reporter: [['text', { skipFull: false }], 'html'],
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/main.tsx', 'src/setupTests.ts', 'src/**/*.test.*'],
    },
  },
})
