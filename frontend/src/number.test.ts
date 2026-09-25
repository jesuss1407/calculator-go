import { describe, expect, it } from 'vitest'
import { formatNumber, parseNumber } from './number'

describe('parseNumber', () => {
  it.each([
    ['3', 3],
    ['-3', -3],
    ['+3', 3],
    ['0', 0],
    ['2.5', 2.5],
    ['.5', 0.5],
    ['5.', 5],
    ['1e5', 100000],
    ['1.5E-3', 0.0015],
    ['  42  ', 42],
  ])('accepts %j', (input, expected) => {
    expect(parseNumber(input)).toEqual({ ok: true, value: expected })
  })

  it.each(['', '   '])('reports %j as required', (input) => {
    expect(parseNumber(input)).toEqual({ ok: false, error: 'Required' })
  })

  it.each(['abc', '1.2.3', '-', '.', '1e', '1,5', '0x10', 'Infinity', 'NaN', '1e400'])(
    'rejects %j',
    (input) => {
      expect(parseNumber(input)).toEqual({ ok: false, error: 'Enter a valid number' })
    },
  )
})

describe('formatNumber', () => {
  it.each([
    [0.1 + 0.2, '0.3'],
    [5, '5'],
    [-2.5, '-2.5'],
    [-0, '0'],
    [123456789012345, '123456789012345'],
    [1e21, '1e+21'],
    [1e-7, '1e-7'],
    [Number.MAX_VALUE, '1.7976931348623157e+308'],
    [-Number.MAX_VALUE, '-1.7976931348623157e+308'],
  ])('formats %d as %j', (value, expected) => {
    expect(formatNumber(value)).toBe(expected)
  })
})
