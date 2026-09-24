import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { InfoPopover } from './InfoPopover'
import { parseNumber } from './number'

// jsdom applies the browser rule that hides a closed popover, but it cannot open one.
// So these tests check the wiring, and read the closed popover's content with { hidden: true }.
describe('InfoPopover', () => {
  it('wires both buttons to the popover, which starts closed', () => {
    render(<InfoPopover />)

    const popover = document.getElementById('info-popover')
    expect(popover).toHaveAttribute('popover', 'auto')
    expect(popover).not.toBeVisible()

    expect(screen.getByRole('button', { name: 'How to use' })).toHaveAttribute('popovertarget', 'info-popover')

    const close = screen.getByRole('button', { name: 'Close', hidden: true })
    expect(close).toHaveAttribute('popovertarget', 'info-popover')
    expect(close).toHaveAttribute('popovertargetaction', 'hide')
  })

  it('only shows examples that agree with the validation rules', () => {
    render(<InfoPopover />)

    const valid = examplesIn('Valid numbers')
    const invalid = examplesIn('Not accepted')

    expect(valid.length).toBeGreaterThan(0)
    expect(invalid.length).toBeGreaterThan(0)
    for (const example of valid) {
      expect(parseNumber(example).ok, `"${example}" should be valid`).toBe(true)
    }
    for (const example of invalid) {
      expect(parseNumber(example).ok, `"${example}" should be rejected`).toBe(false)
    }
  })
})

function examplesIn(listName: string): string[] {
  const list = screen.getByRole('list', { name: listName, hidden: true })
  return Array.from(list.querySelectorAll('code'), (code) => code.textContent ?? '')
}
