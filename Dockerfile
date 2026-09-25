# One image that serves the frontend and the API from a single Go process.
#
#   docker build -t calculator .
#   docker run --rm -p 8080:8080 calculator      # then open http://localhost:8080

# 1) Frontend: compile the React app into static files (dist/).
FROM node:24-alpine AS frontend
WORKDIR /src/frontend
# Copy the manifests first so the npm install layer stays cached until dependencies change.
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# 2) Backend: a static Go binary (no cgo), so it runs on a minimal base image.
FROM golang:1.27-alpine AS backend
WORKDIR /src/backend
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

# 3) Runtime: only the binary and the built frontend. Distroless has no shell or
#    package manager, and the :nonroot variant runs as an unprivileged user.
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=backend /out/server /app/server
COPY --from=frontend /src/frontend/dist /app/public
ENV PORT=8080 \
    STATIC_DIR=/app/public
EXPOSE 8080
ENTRYPOINT ["/app/server"]
