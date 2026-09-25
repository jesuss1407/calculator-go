import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'

function mockFetch(implementation: () => Promise<Response>) {
  const fetchMock = vi.fn(implementation)
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('calculate', () => {
  it('posts the operation and operands as JSON and returns the result', async () => {
    const fetchMock = mockFetch(async () => Response.json({ result: 5 }))

    await expect(calculate('add', 2, 3)).resolves.toBe(5)

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/calculate',
      expect.objectContaining({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ operation: 'add', a: 2, b: 3 }),
      }),
    )
  })

  it('sends no b for square root', async () => {
    const fetchMock = mockFetch(async () => Response.json({ result: 3 }))

    await expect(calculate('sqrt', 9)).resolves.toBe(3)

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/calculate',
      expect.objectContaining({ body: '{"operation":"sqrt","a":9}' }),
    )
  })

  it('throws the error message returned by the server', async () => {
    mockFetch(async () => Response.json({ error: 'cannot divide by zero' }, { status: 422 }))

    await expect(calculate('divide', 1, 0)).rejects.toThrow('cannot divide by zero')
  })

  it('throws a fallback message when an error response is not JSON', async () => {
    mockFetch(async () => new Response('<html>Bad Gateway</html>', { status: 502 }))

    await expect(calculate('add', 1, 2)).rejects.toThrow('Request failed (502)')
  })

  it('throws a network message when the server cannot be reached', async () => {
    mockFetch(async () => {
      throw new TypeError('Failed to fetch')
    })

    await expect(calculate('add', 1, 2)).rejects.toThrow('Could not reach the server. Please try again.')
  })

  it('throws when a successful response has no numeric result', async () => {
    mockFetch(async () => Response.json({ result: 'five' }))

    await expect(calculate('add', 2, 3)).rejects.toThrow('Unexpected response from the server.')
  })
})
