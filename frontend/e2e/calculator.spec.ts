import { expect, test } from '@playwright/test'

// A deliberately small suite. Each test covers something the unit and component tests
// can't: the real browser talking to the real Go API, native browser behavior (form
// submission, the Popover API) and real layout.

test.beforeEach(async ({ page }) => {
  await page.goto('/')
})

test('calculates through the real API and rounds the result for display', async ({ page }) => {
  await page.getByLabel('First number').fill('0.1')
  await page.getByRole('radio', { name: 'Add' }).check()
  await page.getByLabel('Second number').fill('0.2')
  await page.getByRole('button', { name: 'Calculate' }).click()

  // The API returns 0.30000000000000004; the UI shows it rounded.
  await expect(page.getByRole('status')).toHaveText('0.3')
})

test('square root sends only one number and submits with Enter', async ({ page }) => {
  await page.getByRole('radio', { name: 'Square root' }).check()
  await expect(page.getByLabel('Second number')).toHaveCount(0)

  // The API rejects a b for sqrt, so a correct result also proves none was sent.
  const number = page.getByLabel('Number', { exact: true })
  await number.fill('16')
  await number.press('Enter')

  await expect(page.getByRole('status')).toHaveText('4')
})

test('shows the error returned by the API', async ({ page }) => {
  await page.getByLabel('First number').fill('1')
  await page.getByRole('radio', { name: 'Divide' }).check()
  await page.getByLabel('Second number').fill('0')
  await page.getByRole('button', { name: 'Calculate' }).click()

  await expect(page.getByRole('alert')).toHaveText('cannot divide by zero')
  await expect(page.getByRole('status')).toHaveText('—')
})

test('opens and closes the help popup', async ({ page }) => {
  const helpButton = page.getByRole('button', { name: 'How to use' })
  const helpContent = page.getByRole('heading', { name: 'Valid numbers' })
  await expect(helpContent).toBeHidden()

  await helpButton.click()
  await expect(helpContent).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(helpContent).toBeHidden()

  await helpButton.click()
  await page.getByRole('button', { name: 'Close' }).click()
  await expect(helpContent).toBeHidden()
})

test('fits a 320px-wide screen with full-size touch targets', async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 640 })
  // Waits for React to render all seven buttons, so the checks below can't pass on an empty page.
  await expect(page.getByRole('radio')).toHaveCount(7)

  const pageWidth = await page.evaluate(() => document.documentElement.scrollWidth)
  expect(pageWidth, 'page scrolls horizontally').toBeLessThanOrEqual(320)

  for (const operation of await page.getByRole('radio').all()) {
    const box = await operation.boundingBox()
    expect(box?.width, 'operation button width').toBeGreaterThanOrEqual(44)
    expect(box?.height, 'operation button height').toBeGreaterThanOrEqual(44)
  }
})
