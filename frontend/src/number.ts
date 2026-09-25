// Optional sign, digits with an optional decimal point, optional exponent.
// Accepts "3", "-3", "+3", ".5", "5.", "1e5". Rejects "1.2.3", "0x10", "Infinity", "NaN".
const NUMBER_PATTERN = /^[+-]?(\d+\.?\d*|\.\d+)(e[+-]?\d+)?$/i

export type ParseResult = { ok: true; value: number } | { ok: false; error: string }

// Number() alone is too lenient: Number('') is 0 and Number('0x10') is 16.
export function parseNumber(input: string): ParseResult {
  const trimmed = input.trim()
  if (trimmed === '') {
    return { ok: false, error: 'Required' }
  }

  const value = Number(trimmed)
  if (!NUMBER_PATTERN.test(trimmed) || !Number.isFinite(value)) {
    return { ok: false, error: 'Enter a valid number' }
  }
  return { ok: true, value }
}

// Rounds to 15 significant digits (as spreadsheets do) to hide floating-point
// noise such as 0.1 + 0.2 = 0.30000000000000004.
export function formatNumber(value: number): string {
  const rounded = Number(value.toPrecision(15))
  // Rounding a value next to the float64 maximum can overflow to Infinity; show it as-is then.
  return String(Number.isFinite(rounded) ? rounded : value)
}
