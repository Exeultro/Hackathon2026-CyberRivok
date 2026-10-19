FROM node:20-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
COPY content/ /content/
RUN npm run build

FROM golang:1.26.4-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H app
WORKDIR /app
COPY --from=backend /out/api /app/api
COPY --from=backend /src/internal/store/migrations /app/internal/store/migrations
COPY --from=backend /src/frontend/dist /app/frontend/dist
COPY --from=backend /src/content /app/content
USER app
EXPOSE 8000
ENTRYPOINT ["/app/api"]