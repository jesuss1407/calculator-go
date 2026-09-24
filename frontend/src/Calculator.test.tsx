import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'
import { Calculator } from './Calculator'

vi.mock('./api')
const calculateMock = vi.mocked(calculate)

type User = ReturnType<typeof userEvent.setup>

beforeEach(() => {
  calculateMock.mockReset()
})

async function renderAndFill(a: string, operation: string, b: string) {
  const user = userEvent.setup()
  render(<Calculator />)
  await user.type(screen.getByLabelText('First number'), a)
  await user.click(screen.getByRole('radio', { name: operation }))
  await user.type(screen.getByLabelText('Second number'), b)
  return user
}

function submitButton() {
  return screen.getByRole('button', { name: /^calculat/i })
}

describe('Calculator', () => {
  it('shows field errors and does not call the API when input is invalid', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.type(screen.getByLabelText('First number'), 'abc')
    await user.click(submitButton())

    expect(screen.getByText('Enter a valid number')).toBeInTheDocument()
    expect(screen.getByText('Required')).toBeInTheDocument()
    expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'true')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('sends the parsed operands and shows the formatted result', async () => {
    calculateMock.mockResolvedValue(0.30000000000000004)
    const user = await renderAndFill('0.1', 'Add', '0.2')

    await user.click(submitButton())

    expect(await screen.findByText('0.3')).toBeInTheDocument()
    expect(calculateMock).toHaveBeenCalledWith('add', 0.1, 0.2)
  })

  it('shows the error message from the API', async () => {
    calculateMock.mockRejectedValue(new Error('cannot divide by zero'))
    const user = await renderAndFill('1', 'Divide', '0')

    await user.click(submitButton())

    expect(await screen.findByRole('alert')).toHaveTextContent('cannot divide by zero')
  })

  it('disables the form while the request is pending', async () => {
    calculateMock.mockReturnValue(new Promise(() => {}))
    const user = await renderAndFill('2', 'Add', '3')

    await user.click(submitButton())

    expect(submitButton()).toHaveTextContent('Calculating…')
    expect(submitButton()).toBeDisabled()
    expect(screen.getByLabelText('First number')).toBeDisabled()
    expect(screen.getByRole('radio', { name: 'Subtract' })).toBeDisabled()
  })

  it.each([
    { change: 'an operand', makeChange: (user: User) => user.type(screen.getByLabelText('Second number'), '0') },
    { change: 'the operation', makeChange: (user: User) => user.click(screen.getByRole('radio', { name: 'Multiply' })) },
  ])('clears the result when $change changes', async ({ makeChange }) => {
    calculateMock.mockResolvedValue(5)
    const user = await renderAndFill('2', 'Add', '3')
    await user.click(submitButton())
    expect(await screen.findByText('5')).toBeInTheDocument()

    await makeChange(user)

    expect(screen.queryByText('5')).not.toBeInTheDocument()
  })
})
