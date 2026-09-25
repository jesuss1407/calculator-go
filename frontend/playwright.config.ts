import { defineConfig, devices } from '@playwright/test'

// The E2E tests run against an app that is already running. By default that's the Docker
// container on port 8080 (`docker run --rm -p 8080:8080 calculator`), so they test what ships.
// Set BASE_URL to target something else, such as the Vite dev server on port 5173.
export default defineConfig({
  testDir: './e2e',
  forbidOnly: !!process.env.CI,
  reporter: process.env.CI ? [['github'], ['html', { open: 'never' }]] : 'list',
  use: {
    baseURL: process.env.BASE_URL ?? 'http://localhost:8080',
    trace: 'retain-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
})
