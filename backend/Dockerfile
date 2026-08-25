# ---- Etapa 1: build (imagen con el compilador de Go) ----
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ---- Etapa 2: runtime (imagen mínima, sin compilador) ----
FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]