export type Operation = 'add' | 'subtract' | 'multiply' | 'divide'

const CALCULATE_URL = '/api/v1/calculate'
const TIMEOUT_MS = 5000

type ResponseBody = { result?: unknown; error?: unknown } | null

// Sends one calculation to the backend and returns the result.
// Every failure is thrown as an Error whose message is safe to show to the user.
export async function calculate(operation: Operation, a: number, b: number): Promise<number> {
  let response: Response
  try {
    response = await fetch(CALCULATE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ operation, a, b }),
      signal: AbortSignal.timeout(TIMEOUT_MS),
    })
  } catch {
    throw new Error('Could not reach the server. Please try again.')
  }

  // Error pages from proxies may not be JSON, so a parse failure is not fatal here.
  const body = (await response.json().catch(() => null)) as ResponseBody

  if (!response.ok) {
    throw new Error(typeof body?.error === 'string' ? body.error : `Request failed (${response.status})`)
  }
  if (typeof body?.result !== 'number') {
    throw new Error('Unexpected response from the server.')
  }
  return body.result
}
