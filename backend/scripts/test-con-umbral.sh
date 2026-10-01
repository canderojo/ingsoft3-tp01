#!/bin/sh
# Corre los tests del service, mide la cobertura y FALLA si queda por debajo
# del umbral. Go no trae una opción de umbral (como el Threshold de coverlet o
# los thresholds de vitest), así que la comparación se hace en este script.
set -e

UMBRAL="${UMBRAL:-70}"   # % mínimo de sentencias (justificado en decisiones.md)
SALIDA="${SALIDA:-/out}" # carpeta donde queda coverage.out

mkdir -p "$SALIDA"

# Sólo se mide internal/service: es donde están las reglas de negocio.
go test -v ./internal/service/... -coverprofile="$SALIDA/coverage.out"

TOTAL=$(go tool cover -func="$SALIDA/coverage.out" | tail -n 1 | awk '{print $NF}' | tr -d '%')
echo "Cobertura de internal/service: ${TOTAL}% de sentencias (umbral: ${UMBRAL}%)"

if awk -v t="$TOTAL" -v u="$UMBRAL" 'BEGIN { exit !(t < u) }'; then
  echo "ERROR: la cobertura (${TOTAL}%) no llega al umbral (${UMBRAL}%)"
  exit 1
fi
