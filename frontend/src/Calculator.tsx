import { useState, type FormEvent } from 'react'
import { calculate, type Operation } from './api'
import { InfoPopover } from './InfoPopover'
import { formatNumber, parseNumber } from './number'

const OPERATIONS: { value: Operation; label: string; symbol: string }[] = [
  { value: 'add', label: 'Add', symbol: '+' },
  { value: 'subtract', label: 'Subtract', symbol: '−' },
  { value: 'multiply', label: 'Multiply', symbol: '×' },
  { value: 'divide', label: 'Divide', symbol: '÷' },
]

type Operand = 'a' | 'b'

// A single status value makes it impossible to show a result and an error at the same time.
type Status =
  | { kind: 'idle' }
  | { kind: 'loading' }
  | { kind: 'success'; result: number }
  | { kind: 'error'; message: string }

export function Calculator() {
  const [operands, setOperands] = useState<Record<Operand, string>>({ a: '', b: '' })
  const [operation, setOperation] = useState<Operation>('add')
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<Operand, string>>>({})
  const [status, setStatus] = useState<Status>({ kind: 'idle' })

  // Any edit resets the status, so a previous result is never shown next to new inputs.
  function handleOperandChange(operand: Operand, value: string) {
    setOperands((current) => ({ ...current, [operand]: value }))
    setFieldErrors((current) => ({ ...current, [operand]: undefined }))
    setStatus({ kind: 'idle' })
  }

  function handleOperationChange(value: Operation) {
    setOperation(value)
    setStatus({ kind: 'idle' })
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const a = parseNumber(operands.a)
    const b = parseNumber(operands.b)
    if (!a.ok || !b.ok) {
      setFieldErrors({ a: a.ok ? undefined : a.error, b: b.ok ? undefined : b.error })
      return
    }

    setStatus({ kind: 'loading' })
    try {
      const result = await calculate(operation, a.value, b.value)
      setStatus({ kind: 'success', result })
    } catch (error) {
      setStatus({ kind: 'error', message: error instanceof Error ? error.message : 'Something went wrong.' })
    }
  }

  const isLoading = status.kind === 'loading'

  return (
    <main className="calculator">
      <header className="header">
        <h1>Calculator</h1>
        <InfoPopover />
      </header>

      <form onSubmit={handleSubmit}>
        {/* Disabling the fieldset locks every control while a request is in flight. */}
        <fieldset className="controls" disabled={isLoading}>
          <NumberField
            id="a"
            label="First number"
            value={operands.a}
            error={fieldErrors.a}
            onChange={(value) => handleOperandChange('a', value)}
          />

          <fieldset>
            <legend>Operation</legend>
            <div className="operations">
              {OPERATIONS.map(({ value, label, symbol }) => (
                <label key={value} className="operation">
                  <input
                    type="radio"
                    name="operation"
                    value={value}
                    checked={operation === value}
                    onChange={() => handleOperationChange(value)}
                  />
                  <span className="operation-symbol" aria-hidden="true">
                    {symbol}
                  </span>
                  <span className="visually-hidden">{label}</span>
                </label>
              ))}
            </div>
          </fieldset>

          <NumberField
            id="b"
            label="Second number"
            value={operands.b}
            error={fieldErrors.b}
            onChange={(value) => handleOperandChange('b', value)}
          />

          <button type="submit" className="submit">
            {isLoading ? 'Calculating…' : 'Calculate'}
          </button>
        </fieldset>

        <div className="result">
          <span>Result</span>
          <output aria-live="polite">{status.kind === 'success' ? formatNumber(status.result) : '—'}</output>
        </div>

        {status.kind === 'error' && (
          <p className="request-error" role="alert">
            {status.message}
          </p>
        )}
      </form>
    </main>
  )
}

type NumberFieldProps = {
  id: string
  label: string
  value: string
  error?: string
  onChange: (value: string) => void
}

// type="text" rather than type="number": number inputs report "" for invalid text,
// which would hide what the user actually typed from validation.
function NumberField({ id, label, value, error, onChange }: NumberFieldProps) {
  const errorId = `${id}-error`

  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        type="text"
        inputMode="decimal"
        autoComplete="off"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        aria-invalid={error !== undefined}
        aria-describedby={error ? errorId : undefined}
      />
      {error && (
        <p id={errorId} className="field-error">
          {error}
        </p>
      )}
    </div>
  )
}
