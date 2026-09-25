# Calculator

A full-stack calculator: a **Go REST API** built only on the standard library, and a **React + TypeScript** frontend built with Vite.

It supports addition, subtraction, multiplication, division, powers, square roots and percentages. Both sides validate input, errors are clear, the form works from phone to desktop, and both apps have unit tests with coverage reports.

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

```
calculator-go/
├── .github/workflows/ci.yml         GitHub Actions: lint, test, coverage, build, and a Docker smoke test
├── Dockerfile                       One image that serves the frontend and the API together
├── .dockerignore
├── backend/                         Go module, no third-party dependencies
│   ├── main.go                      Wiring: reads PORT and STATIC_DIR, starts the HTTP server
│   ├── main_test.go                 Routing when the built frontend is served too
│   └── internal/
│       ├── calculator/              Pure arithmetic, no HTTP or JSON
│       │   ├── calculator.go
│       │   └── calculator_test.go
│       └── api/                     HTTP layer: JSON decoding, validation, status codes
│           ├── handler.go
│           └── handler_test.go
└── frontend/                        React 19 + TypeScript + Vite
    ├── vite.config.ts               Dev proxy to the backend, test and coverage config
    └── src/
        ├── main.tsx                 Entry point
        ├── Calculator.tsx           The form: state, validation messages, result
        ├── InfoPopover.tsx          "How to use" help popup
        ├── api.ts                   The only code that talks to the backend
        ├── number.ts                Parsing and display formatting
        ├── index.css                All styles, mobile-first
        └── *.test.ts(x)             Unit tests next to the code they cover
```

---

## Prerequisites

| Tool | Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.27 or later | `go.mod` declares `go 1.27.0`. With Go 1.21 or later and the default `GOTOOLCHAIN=auto`, the required toolchain downloads automatically. |
| [Node.js](https://nodejs.org/) | 20.19+ or 22.12+ | Required by Vite 8. npm comes with Node. |
| [Docker](https://docs.docker.com/get-docker/) | Any recent version | Optional. Only needed to [run everything in one container](#option-b-docker-one-container); then Go and Node aren't needed at all. |

Developed and tested with Go 1.27.0, Node 25.8.1 and Docker 29.8 on Windows 11.

---

## Setup

```bash
git clone <repository-url>
cd calculator-go

# Frontend dependencies (the backend has none to install)
cd frontend
npm install
```

---

## Running the app

There are two ways to run it:

- **Option A: development servers.** Two terminals, with hot reload. Use this while working on the code.
- **Option B: Docker.** One command, one container. Use this to try the app without installing Go or Node.

### Option A: development servers

Run the backend and the frontend in two terminals.

**Terminal 1: backend** (http://localhost:8080)

```bash
cd backend
go run .
```

To use another port, set `PORT`, for example `PORT=9090 go run .`. The frontend's dev proxy targets port 8080, so change `server.proxy` in `frontend/vite.config.ts` as well if you do.

**Terminal 2: frontend** (http://localhost:5173)

```bash
cd frontend
npm run dev
```

Open **http://localhost:5173**. The Vite dev server forwards every `/api` request to the backend, so the browser only ever talks to one origin and no CORS setup is needed.

### Option B: Docker (one container)

From the repository root:

```bash
docker build -t calculator .
docker run --rm -p 8080:8080 calculator
```

Open **http://localhost:8080**. This one address serves both the page and the API. The curl examples below work against it unchanged. Press Ctrl+C, or run `docker stop`, to shut it down.

To publish it on a different host port, change only the left-hand number: `-p 3000:8080` serves the app at http://localhost:3000. Do this if port 8080 is already taken, for example by `go run .` from Option A.

**How the image is built.** The [Dockerfile](Dockerfile) has three stages:

| Stage | Base image | What it does |
|---|---|---|
| `frontend` | `node:24-alpine` | `npm ci`, then `npm run build`, which type-checks and bundles into `dist/` |
| `backend` | `golang:1.27-alpine` | Compiles a static Go binary (`CGO_ENABLED=0`) |
| runtime | `gcr.io/distroless/static-debian13:nonroot` | Holds only the binary and the built `dist/` files. It runs as a non-root user and has no shell or package manager. |

The final image is about **9 MB**. The binary serves the frontend because the image sets `STATIC_DIR=/app/public`. When `STATIC_DIR` is unset, as in Option A, it serves only the API.

### Using the app

Enter two numbers, pick an operation and press **Calculate** or Enter. Square root (√) needs only one number, so the second field is hidden while it's selected. The **i** button next to the title lists the number formats that are accepted.

---

## Testing and coverage

### Backend

```bash
cd backend
go vet ./...
go test ./...

# Coverage
go test ./... -coverprofile coverage.out
go tool cover -func coverage.out                       # per-function summary
go tool cover -html coverage.out -o coverage.html      # optional browsable report
```

> The flags use the space-separated form (`-coverprofile coverage.out`) on purpose.

### Frontend

```bash
cd frontend
npm test              # Vitest, React Testing Library, jsdom
npm run coverage      # tests + coverage (text summary, HTML in frontend/coverage/)
npm run lint          # Oxlint
npm run build         # type-check (tsc) + production bundle
```

### Coverage results

These numbers were measured on 2026-09-25. Coverage is reported, not targeted: the tests concentrate on business rules and on every error a user or API client can trigger.

**Backend** (86 test cases: 51 for `calculator`, 28 for `api`, 7 for `main`)

| Package | Statements | Not covered |
|---|---|---|
| `internal/calculator` | 100.0% | — |
| `internal/api` | 92.3% | The 500 fallback (no current error can reach it) and the log line for a failed response write |
| `main` | 38.1% | `newHandler`, the routing, is at 100%. `main()` itself (reading env vars, starting the server) is wiring and deliberately not unit-tested. |
| **Total** | **80.7%** | |

**Frontend** (46 tests in 4 files)

| File | Statements | Branches |
|---|---|---|
| `Calculator.tsx` | 100% | 88.46% |
| `InfoPopover.tsx` | 100% | 100% |
| `api.ts` | 100% | 100% |
| `number.ts` | 100% | 100% |
| **All files** | **100%** | **92.1%** |

The three untested branches are all in the submit handler: some combinations of which fields are invalid, and the fallback message for a thrown value that isn't an `Error`. That fallback can't happen today, because `api.ts` only throws `Error`s.

### Continuous integration

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) runs on every push to `main`, on every pull request, and on demand from the Actions tab. The backend and frontend jobs run in parallel on Ubuntu. The Docker job starts only after both pass, so no time is spent building an image from code that fails its tests. Each job uses the same commands as above:

| Job | Steps |
|---|---|
| **Backend (Go)** | `gofmt` check → `go vet` → `go test` with coverage → coverage summary |
| **Frontend (React)** | `npm ci` → `npm run lint` → `npm run coverage` → `npm run build` |
| **Docker image** | `docker build` → start the container → smoke test: `GET /health`, the page is served, and `2 + 3` returns `{"result":5}` |

- **Toolchains:** the Go version comes from `backend/go.mod`, so CI and local development can't drift apart. The frontend job uses Node 24 LTS with the npm cache enabled.
- **Coverage:** the backend and frontend jobs write their coverage tables to the run's summary page. Coverage is reported but never used as a pass/fail gate.
- **Why a smoke test and not just a build:** an image can build but still not work, for example if the server can't find the frontend files. The smoke test runs the real container and fails on any unexpected response. The container's logs are printed either way, so a failure is easy to diagnose.
- **Permissions:** the workflow only has read access to the repository.

> **Windows note:** with Git's default `core.autocrlf=true`, a local `gofmt -l` can list files only because the working copy has CRLF line endings. CI checks out LF, so those files still pass there.

---

## API reference

Base URL: `http://localhost:8080`. Every response from the handlers has `Content-Type: application/json`.

### `POST /api/v1/calculate`

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | One of the operations below (lowercase) |
| `a` | number | yes | First operand |
| `b` | number | for every operation except `sqrt` | Second operand. For `sqrt` it must be left out. |

| `operation` | Result | Operands |
|---|---|---|
| `add` | a + b | `a`, `b` |
| `subtract` | a − b | `a`, `b` |
| `multiply` | a × b | `a`, `b` |
| `divide` | a ÷ b | `a`, `b` |
| `power` | a raised to the power b | `a`, `b` |
| `sqrt` | √a | `a` only |
| `percentage` | a% of b, computed as `(a / 100) * b` | `a`, `b` |

```json
{ "operation": "divide", "a": 10, "b": 4 }
{ "operation": "sqrt", "a": 16 }
```

**Success: `200 OK`**

```json
{ "result": 2.5 }
```

**Errors** all use one shape:

```json
{ "error": "cannot divide by zero" }
```

| Status | `error` message | When |
|---|---|---|
| `400` | `request body must be a JSON object with operation, a and b` | Empty or malformed body, wrong types (`"a": "2"`), unknown fields, a number outside float64 range (`1e999`) |
| `400` | `request body must contain a single JSON object` | Extra data after the JSON object, such as `{…}{}` |
| `400` | `a and b are required` | `a` or `b` missing or `null`, for any operation except `sqrt` |
| `400` | `sqrt takes only a` | `sqrt` without `a`, or with a `b` |
| `400` | `unknown operation` | `operation` missing or not one of the supported values |
| `422` | `cannot divide by zero` | `b` is 0 in a division (including `0 / 0`), or `0` raised to a negative power |
| `422` | `result is not a real number` | The square root of a negative number, or a negative base with a fractional exponent, such as `(-4)^0.5` |
| `422` | `result is out of range` | The result overflows float64, such as `1e308 * 10` or `10^400` |
| `500` | `internal server error` | Unexpected failure (details are logged on the server) |
| `405` | `Method Not Allowed` (plain text) | Wrong HTTP method, returned by Go's router |

### `GET /health`

```json
{ "status": "ok" }
```

### Example calls

With curl (macOS, Linux, Git Bash):

```bash
curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":2,"b":3}'
# {"result":5}
```

In PowerShell 7+, call `curl.exe` explicitly and put the JSON in single quotes:

```powershell
curl.exe -s -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d '{"operation":"add","a":2,"b":3}'
```

More requests and the responses they return. Each was run against the running server:

| Request body | Status | Response |
|---|---|---|
| `{"operation":"add","a":2,"b":3}` | 200 | `{"result":5}` |
| `{"operation":"divide","a":10,"b":4}` | 200 | `{"result":2.5}` |
| `{"operation":"power","a":2,"b":10}` | 200 | `{"result":1024}` |
| `{"operation":"sqrt","a":2}` | 200 | `{"result":1.4142135623730951}` |
| `{"operation":"percentage","a":10,"b":200}` | 200 | `{"result":20}` |
| `{"operation":"percentage","a":7,"b":300}` | 200 | `{"result":21.000000000000004}` (the UI shows `21`) |
| `{"operation":"divide","a":1,"b":0}` | 422 | `{"error":"cannot divide by zero"}` |
| `{"operation":"power","a":0,"b":-1}` | 422 | `{"error":"cannot divide by zero"}` |
| `{"operation":"sqrt","a":-4}` | 422 | `{"error":"result is not a real number"}` |
| `{"operation":"power","a":-4,"b":0.5}` | 422 | `{"error":"result is not a real number"}` |
| `{"operation":"multiply","a":1e308,"b":10}` | 422 | `{"error":"result is out of range"}` |
| `{"operation":"add","a":1}` | 400 | `{"error":"a and b are required"}` |
| `{"operation":"sqrt","a":16,"b":2}` | 400 | `{"error":"sqrt takes only a"}` |
| `{"operation":"modulo","a":1,"b":2}` | 400 | `{"error":"unknown operation"}` |
| `{"operation":` | 400 | `{"error":"request body must be a JSON object with operation, a and b"}` |
| `{"operation":"add","a":1,"b":2}{}` | 400 | `{"error":"request body must contain a single JSON object"}` |

```bash
curl -s http://localhost:8080/health                 # {"status":"ok"}
curl -s -i http://localhost:8080/api/v1/calculate    # GET → 405 Method Not Allowed
```

---

## Design decisions

The guiding constraint was **the smallest code that still separates concerns, is testable and is easy to explain**.

### Architecture

```
Browser                                        Go server
┌─────────────────────────────┐               ┌──────────────────────────────────────┐
│ Calculator.tsx  UI + state  │               │ main.go             wiring           │
│ number.ts       parse/format│   POST        │ internal/api        HTTP + JSON      │
│ api.ts          fetch ──────┼──────────────▶│        │ calls                       │
└─────────────────────────────┘ /api/v1/      │ internal/calculator pure arithmetic  │
                                 calculate    └──────────────────────────────────────┘
```

- **The math doesn't know about HTTP.** `internal/calculator` imports only `errors` and `math`. It exposes `Calculate(op, a, b) (float64, error)`, four sentinel errors and `Operation.IsUnary()`. The HTTP layer maps those errors to status codes with `errors.Is`. The domain decides *what* went wrong; the handler decides *how* to report it.
- **No interfaces or mocks on the backend.** The calculator is pure and fast, so the handler tests call the real one. An interface that existed only for mocking would add indirection with no benefit.
- **Standard library only.** Go 1.22+ routing patterns (`"POST /api/v1/calculate"`) provide method matching and automatic 405 responses, so no router framework is needed.
- **The UI never calls `fetch`.** `api.ts` owns the URL, the request, a 5-second timeout and the translation of every failure into a message that's safe to show. That makes it the single seam the component tests mock.
- **No abstractions ahead of need.** Seven operations still need no operation registry, no custom hook and no folder hierarchy. When square root added the first one-operand operation, one `IsUnary()` method was all the handler needed. Each abstraction can be added when a real need appears.

### API

- **One endpoint with the operation in the body,** not `/add`, `/divide` and so on. A calculation is an action, not a resource. One endpoint means one handler, one validation path and one client function.
- **Named `a` and `b` fields.** They read better than an array and give precise errors. They're decoded into `*float64` so a missing or `null` operand is rejected; with a plain `float64`, Go's JSON decoder would silently turn it into `0`.
- **400 vs 422.** 400 means the request is malformed. 422 means it's well-formed, but the math has no answer. A client can tell "fix your request" apart from "this calculation is undefined".
- **Strict decoding.** Unknown fields are rejected, which catches typos. So is anything after the first JSON object: the handler decodes again and requires end of input.
- **Results are always finite.** JSON can't represent Infinity or NaN; without a check, an overflow would produce a broken response. `Calculate` turns any non-finite result into `result is out of range`, whatever the inputs.
- **A flat `{"error": "..."}` body.** Nothing needs machine-readable codes yet, and the status code already gives the category.
- **A versioned path (`/api/v1`) and `GET /health`.** Both are cheap now and awkward to add later.
- **Unary operations reject `b` rather than ignore it.** The domain says which operations use only `a` (`IsUnary`), and the handler enforces the matching JSON shape. A `b` sent with `sqrt` is almost certainly a client mistake, so it gets a 400, consistent with rejecting unknown fields. Adding `sqrt` changed none of the existing messages.
- **Explicit domain errors instead of relying on NaN or Infinity.** `math.Pow` and `math.Sqrt` return NaN or ±Inf for undefined results. `Calculate` checks those cases first, so clients get a specific message: `0^-1` is a division by zero, and `√-4` or `(-4)^0.5` is "not a real number". A general "out of range" is left only for real overflow.

### Validation: who checks what

- **Frontend: format.** `parseNumber` uses a strict pattern plus `Number.isFinite`. `Number()` alone is too lenient: `Number('')` is `0` and `Number('0x10')` is `16`. Invalid input never reaches the network.
- **Backend: meaning.** Division by zero, results that aren't real numbers, and overflow are only checked on the server, so those rules live in one place. The backend also validates the request shape itself and never assumes the client already did.

### Frontend

- **`<input type="text" inputMode="decimal">` rather than `type="number"`.** A number input reports `""` for partial or invalid text, which hides what the user typed from validation. `inputMode` still brings up a numeric keyboard on phones.
- **One status value** (`idle | loading | success | error`) rather than separate flags, so a result and an error can never appear together.
- **Stale results are cleared.** Editing either number or changing the operation resets the status to idle, so a result never sits next to inputs that didn't produce it.
- **The form is locked while a request is in flight.** A disabled `<fieldset>` prevents double submits and edits that a late response could overwrite.
- **Display rounding to 15 significant digits**, the precision spreadsheets use. `0.1 + 0.2` shows as `0.3`, and integers up to 15 digits display exactly. The API still returns the full float64 value.
- **A help popup on the native Popover API.** The browser handles opening, closing with Esc or an outside click, and focus order, so the component has no state or effects. A test checks that every example in the popup agrees with `parseNumber`, so the help text can't drift from the real rules.
- **Accessibility.** Labels are linked to inputs, fields get `aria-invalid` and `aria-describedby`, the result is announced through a live region and request errors use `role="alert"`. Operations are real radio buttons, so arrow keys work. Touch targets are at least 44 px.
- **Square root hides the second field.** It isn't validated or sent either, and the first field's label changes from "First number" to "Number". Whatever was typed in the second field is kept for when another operation is chosen.
- **Responsive layout.** Mobile-first, a single column, a centered card on wide screens, and long results wrap rather than scroll sideways at 320 px. The seven operation buttons sit in four columns: `+ − × ÷` on the first row and `xʸ √ %` on the second. That keeps each button at least 44 px wide even on a 320 px screen, with no media query.

### Docker

- **One image, one process, one port.** Instead of separate frontend and backend containers (for example nginx and Go wired together with docker-compose), the Go binary also serves the built frontend. The page and the API share an origin, so the relative `/api/v1/calculate` URL works as it does in development, with no proxy and no CORS.
- **Opt-in static serving.** `main.go` serves files only when `STATIC_DIR` is set. That keeps the `api` package unchanged and local development identical. The routing lives in `newHandler`, and `main_test.go` checks that static files never shadow `/api/...` or `/health`.
- **A multi-stage build.** Node and the Go toolchain exist only in the build stages. The runtime image is distroless and runs as a non-root user, about 9 MB in total.
- **Build-cache friendly.** `package.json` and `package-lock.json` are copied before the source, so `npm ci` only reruns when dependencies change. `.dockerignore` keeps host `node_modules` out of the build, because those contain native binaries for the host platform, not Linux.

### Tooling

- **Vitest + React Testing Library.** Tests mock `api.ts` at the module boundary rather than using a network-mocking library, which saves a dependency and leaves the seam easy to see.
- **Oxlint,** the default linter in Vite's current React template: fast, with no configuration to maintain.
- **Table-driven Go tests** with `httptest` for the handler. Each case asserts the status code, the `Content-Type` and the exact response body.

---

## Assumptions and known limitations

- **One operation per request, on real numbers only.** Each request is one operation on one or two operands. There is no expression parsing (`2 + 3 × 4`) and no complex numbers.
- **Percentage means "a% of b"** and is computed as `(a / 100) * b`. It's an ordinary two-operand operation, so it needed no change to the request shape. The formula keeps plain floating-point behavior: `7% of 300` returns `21.000000000000004` from the API, while the UI's display rounding shows `21`. That matches how `0.1 + 0.2` behaves. Dividing first also means `a * b` never overflows when the result itself fits: `200% of 1e307` returns `2e+307`.
- **Power follows `math.Pow` conventions.** `0^0` is `1`. A negative base with a fractional exponent is rejected even when a real answer exists. For example, `(-8)^(1/3)` returns "not a real number" rather than `-2`, because `1/3` can't be represented exactly and the principal value is complex.
- **The frontend knows which operation is unary.** `Calculator.tsx` hides the second field when `sqrt` is selected. That duplicates the backend's `IsUnary()`, which is acceptable for one operation; an operations endpoint would be the fix if the list grew.
- **IEEE 754 float64 precision.** The API returns `0.1 + 0.2` as `0.30000000000000004`, and integers above 2^53 lose precision. The UI rounds only for display.
- **Negative zero.** A result such as `0 × -1` comes back as `{"result":-0}`. That's valid JSON, and the UI displays it as `0`.
- **Operation names are case-sensitive** (`"ADD"` is rejected). JSON *field names* follow Go's `encoding/json` default of case-insensitive matching, so `"A"` is accepted as `a`.
- **Not production-hardened.** There is no authentication, persistence, rate limiting or CORS configuration. The frontend reaches the API through the Vite dev proxy (Option A) or from the same origin (Docker). The container has no Docker `HEALTHCHECK`, because the distroless image has no shell or curl to run one. `GET /health` exists for an orchestrator to probe instead.
- **Modern browsers.** The help popup needs the Popover API (Chrome 114+, Firefox 125+, Safari 17+). The calculator itself works without it.
- **iPhone keypad.** The keypad iOS shows for `inputMode="decimal"` has only digits and a decimal point, with no minus sign and no `e`. On an iPhone, negative numbers and scientific notation can only be pasted. Desktop browsers and Android keyboards are not affected.

---

## Possible next steps

These were deliberately left out to keep the first version small:

- **Negative-zero normalization,** so that results such as `0 × -1` come back as `0` instead of `-0`.
- **Server hardening:** a request body size limit, `http.Server` timeouts and graceful shutdown. `main.go` currently uses plain `http.ListenAndServe`. That's fine for local use and the demo container, but not for internet-facing production.
- **Browser end-to-end tests:** a Playwright test that fills in the form in a real browser. It could run in the Docker CI job against the container that's already started there.
