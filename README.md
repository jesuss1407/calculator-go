# Calculator

A full-stack calculator with a **Go REST API** built only on the standard library and a **React + TypeScript** frontend built with Vite.

It supports addition, subtraction, multiplication, division, powers, square roots, and percentages. Both layers validate input, errors are handled explicitly, the UI is responsive, and the project includes unit, integration, and end-to-end tests with coverage reports.

---

## Contents

- [Project structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Running the app](#running-the-app)
- [Testing and coverage](#testing-and-coverage)
- [API reference](#api-reference)
- [Design decisions](#design-decisions)
- [Assumptions and known limitations](#assumptions-and-known-limitations)
- [Possible next steps](#possible-next-steps)

---

## Project structure

```text
calculator-go/
├── .github/workflows/ci.yml         GitHub Actions: test, coverage, build, Docker smoke test, E2E
├── Dockerfile                       Multi-stage build for frontend + backend
├── .dockerignore
├── backend/
│   ├── main.go                      Application wiring and optional static-file serving
│   ├── main_test.go                 Routing tests when serving the built frontend
│   └── internal/
│       ├── calculator/              Pure arithmetic, no HTTP or JSON
│       │   ├── calculator.go
│       │   └── calculator_test.go
│       └── api/                     HTTP layer: JSON decoding, validation, status codes
│           ├── handler.go
│           └── handler_test.go
└── frontend/
    ├── vite.config.ts               Dev proxy, Vitest, and coverage config
    ├── playwright.config.ts         E2E config: Chromium and BASE_URL
    ├── e2e/
    │   └── calculator.spec.ts       Playwright end-to-end tests
    └── src/
        ├── main.tsx                 Entry point
        ├── Calculator.tsx           UI and state
        ├── InfoPopover.tsx          Help popup
        ├── api.ts                   Backend communication
        ├── number.ts                Parsing and display formatting
        ├── index.css                Responsive styles
        └── *.test.ts(x)             Unit and component tests
```

---

## Prerequisites

| Tool | Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.27 or later | `go.mod` declares `go 1.27.0` |
| [Node.js](https://nodejs.org/) | 22.13+ (24 LTS recommended) | The minimum required by the test tooling (jsdom, jest-dom); npm is included |
| [Docker](https://docs.docker.com/get-docker/) | Any recent version | Optional; only required for the single-container setup |

---

## Setup

```bash
git clone <repository-url>
cd calculator-go

cd frontend
npm install
```

The Go backend uses only the standard library, so it has no third-party dependencies to install.

---

## Running the app

There are two supported ways to run the project.

### Option A: development servers

Run the backend and frontend in separate terminals.

**Terminal 1: backend**

```bash
cd backend
go run .
```

The API runs at:

```text
http://localhost:8080
```

**Terminal 2: frontend**

```bash
cd frontend
npm run dev
```

Open:

```text
http://localhost:5173
```

Vite forwards `/api` requests to the Go backend, so no CORS configuration is needed during local development.

### Option B: Docker

From the repository root:

```bash
docker build -t calculator .
docker run --rm -p 8080:8080 calculator
```

Open:

```text
http://localhost:8080
```

The same address serves both the React application and the Go API.

The Docker image uses three stages:

| Stage | Base image | Purpose |
|---|---|---|
| `frontend` | `node:24-alpine` | Installs dependencies and builds the React app |
| `backend` | `golang:1.27-alpine` | Compiles a static Go binary |
| runtime | `gcr.io/distroless/static-debian13:nonroot` | Runs only the server binary and built frontend |

The final image is about **9 MB** and runs as a non-root user.

### Using the app

Enter one or two numbers, choose an operation, and press **Calculate** or Enter.

Square root (`√`) is unary, so the second input is hidden while it is selected. The help button next to the title shows accepted number formats.

---

## Testing and coverage

### Backend

```bash
cd backend

go vet ./...
go test ./...

go test ./... -coverprofile coverage.out
go tool cover -func coverage.out
go tool cover -html coverage.out -o coverage.html
```

### Frontend

```bash
cd frontend

npm test
npm run coverage
npm run lint
npm run build
```

### End-to-end (Playwright)

The Playwright suite runs against a real browser and a running full-stack application.

Install Chromium once:

```bash
cd frontend
npx playwright install chromium
```

Run against Docker:

```bash
# terminal 1, repository root
docker build -t calculator .
docker run --rm -p 8080:8080 calculator

# terminal 2
cd frontend
npm run e2e
```

To run against the development servers instead:

```bash
BASE_URL=http://localhost:5173 npm run e2e
```

PowerShell:

```powershell
$env:BASE_URL="http://localhost:5173"
npm run e2e
```

The E2E suite covers:

- a real frontend-to-backend calculation
- square-root behavior and Enter submission
- API error propagation
- the native help popover
- responsive behavior at 320 px

### Coverage results

Coverage is reported rather than used as a pass/fail threshold.

| Area | Coverage |
|---|---:|
| `internal/calculator` | 100.0% |
| `internal/api` | 92.9% |
| Backend total | 82.0% |
| Frontend statements | 100% |
| Frontend branches | 92.5% |

`main()` is primarily application wiring and is intentionally not unit-tested; the routing in `newHandler` is covered.

### Continuous integration

`.github/workflows/ci.yml` runs on pushes to `main`, pull requests, and manual dispatches.

| Job | Checks |
|---|---|
| **Backend (Go)** | `gofmt` → `go vet` → tests with coverage |
| **Frontend (React)** | `npm ci` → lint → tests with coverage → production build |
| **Docker image and E2E** | Docker build → container smoke test → Playwright against the same container |

The backend and frontend jobs run in parallel. The Docker/E2E job starts only after both succeed. Failed Playwright runs upload the report and traces as CI artifacts.

---

## API reference

Base URL:

```text
http://localhost:8080
```

### `POST /api/v1/calculate`

#### Request

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | Supported operation |
| `a` | number | yes | First operand |
| `b` | number | except for `sqrt` | Second operand |

Supported operations:

| `operation` | Result | Operands |
|---|---|---|
| `add` | a + b | `a`, `b` |
| `subtract` | a − b | `a`, `b` |
| `multiply` | a × b | `a`, `b` |
| `divide` | a ÷ b | `a`, `b` |
| `power` | a raised to b | `a`, `b` |
| `sqrt` | √a | `a` |
| `percentage` | a% of b, calculated as `(a / 100) * b` | `a`, `b` |

Examples:

```json
{ "operation": "divide", "a": 10, "b": 4 }
```

```json
{ "operation": "sqrt", "a": 16 }
```

#### Success

`200 OK`

```json
{ "result": 2.5 }
```

#### Errors

Errors returned by the calculate handler use one JSON shape:

```json
{ "error": "cannot divide by zero" }
```

| Status | Example / meaning |
|---|---|
| `400` | malformed JSON, wrong field types, missing operands, unknown operation, invalid request shape |
| `422` | division by zero, non-real result, or result outside float64 range |
| `500` | unexpected server failure |
| `405` | wrong HTTP method, answered by Go's router as plain text (`Method Not Allowed`) |
| `404` | unknown path, answered by Go's router as plain text (`404 page not found`) |

Representative error cases:

| Request | Status | Response |
|---|---:|---|
| `{"operation":"divide","a":1,"b":0}` | 422 | `{"error":"cannot divide by zero"}` |
| `{"operation":"sqrt","a":-4}` | 422 | `{"error":"result is not a real number"}` |
| `{"operation":"sqrt","a":16,"b":2}` | 400 | `{"error":"sqrt takes only a"}` |
| `{"operation":"modulo","a":1,"b":2}` | 400 | `{"error":"unknown operation"}` |

### `GET /health`

```json
{ "status": "ok" }
```

### Example call

```bash
curl -s -X POST http://localhost:8080/api/v1/calculate   -H 'Content-Type: application/json'   -d '{"operation":"percentage","a":10,"b":200}'
```

Response:

```json
{ "result": 20 }
```

---

## Design decisions

The guiding principle was to keep the project **small, testable, and easy to explain**.

- **Clear separation of concerns.** `internal/calculator` contains pure arithmetic, `internal/api` handles HTTP and JSON, and React components use `api.ts` instead of calling `fetch` directly.
- **Go standard library only.** `net/http` is sufficient for the API and keeps the dependency footprint minimal.
- **One calculation endpoint.** Every operation uses `POST /api/v1/calculate`, which keeps routing and validation centralized.
- **Backend owns mathematical rules.** The frontend validates input format; the backend handles division by zero, non-real results, and overflow.
- **Minimal abstractions.** No backend interfaces, global state library, operation registry, or framework layers were added without a concrete need.
- **Unary operations are explicit.** `sqrt` is identified as unary and the API rejects an unnecessary `b` instead of silently ignoring it.
- **Finite API results.** Non-finite results are converted into domain errors because JSON cannot represent `NaN` or `Infinity`.
- **Single production container.** The Go process serves the built React assets and API from one origin, avoiding a second web-server container and production CORS configuration.
- **Layered testing.** Domain tests, HTTP handler tests, React component tests, Docker smoke tests, and a small Playwright suite cover different responsibilities without unnecessary duplication.

---

## Assumptions and known limitations

- **One operation per request.** There is no expression parsing such as `2 + 3 × 4`, and complex numbers are not supported.
- **IEEE-754 float64 arithmetic.** Values such as `0.1 + 0.2` may contain floating-point artifacts in the API response; the UI rounds results for display.
- **Percentage semantics.** `percentage` means **a% of b** and is calculated as `(a / 100) * b`.
- **Power follows Go's `math.Pow`.** `0^0` returns `1`; a negative base with a fractional exponent is rejected as a non-real result.
- **Negative zero is not normalized in the API.** For example, `0 × -1` may return `-0`, while the UI displays `0`.
- **Not production-hardened.** Authentication, persistence, rate limiting, request-size limits, server timeouts, graceful shutdown, and configurable CORS are outside the scope of this take-home.

---

## Possible next steps

The following were deliberately left out to keep the implementation focused:

- normalize negative-zero results
- add request-size limits, server timeouts, and graceful shutdown
- add Firefox and WebKit projects if cross-browser E2E coverage becomes important
