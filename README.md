[![CI](https://github.com/canderojo/ingsoft3-tp01/actions/workflows/ci.yml/badge.svg)](https://github.com/canderojo/ingsoft3-tp01/actions/workflows/ci.yml)

# Turnero — Centro de Salud de la Mujer

Sistema de reserva de turnos para un centro de salud con profesionales de 4 especialidades: **dermatología, nutrición, ecografía y endocrinología**. Los pacientes ven la disponibilidad de cada profesional y reservan turnos; el sistema valida horarios, evita solapamientos y controla las transiciones de estado de cada turno.

Proyecto para la materia **Ingeniería del Software 3** (UCC, 2026) — cátedra Ariel Schwindt. Cada trabajo práctico del cuatrimestre agrega una capa nueva sobre esta misma aplicación (CI, testing, CD, contenedores, IaC, seguridad, observabilidad), en vez de crear un proyecto distinto por TP.

---

## Puesta en marcha con Docker

**Requisitos:** Docker Desktop (o Docker + Docker Compose v2).

### 1. Clonar el repositorio

```bash
git clone https://github.com/canderojo/ingsoft3-tp01.git
cd ingsoft3-tp01
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Completá `DB_PASSWORD` en el `.env` generado.

### 3. Levantar el sistema (compilando localmente)

```bash
docker compose up -d --build
```

Esto construye las imágenes de backend y frontend desde los Dockerfiles del TP2, levanta PostgreSQL, espera a que esté `healthy` y recién entonces inicia el backend.

### 4. Acceder al sistema

| Componente | URL | Descripción |
|---|---|---|
| Frontend | http://localhost:3000 | Listado de profesionales, reserva y consulta de turnos |
| Backend (API) | http://localhost:8080 | API REST en Go |
| Health check | http://localhost:8080/health | Debería devolver `{"status":"ok","db":"ok"}` |
| PostgreSQL | localhost:5432 | Base de datos (dentro de la red de Docker: `db:5432`) |

### Levantar desde imágenes publicadas (GHCR)

Sin compilar localmente, usando las imágenes ya publicadas en GitHub Container Registry:

```bash
cp .env.example .env
docker compose -f docker-compose.registry.yml up -d
```

### Apagar el sistema

```bash
docker compose down       # conserva los datos (volumen db_data)
docker compose down -v    # borra también el volumen de la base
```

---

## Stack 

- **Backend**: Go + [chi] (router) + [sqlx] + [pgx] (driver de Postgres)
- **Frontend**: React + Vite
- **Base de datos**: PostgreSQL
- **Pruebas**: tests unitarios en Go y en el frontend, y pruebas de integración y e2e con Playwright
- **CI/CD**: GitHub Actions, imágenes en GHCR, despliegue en Render
- **Bases en la nube**: Neon (PostgreSQL)
- **Infraestructura como código**: Terraform
- **Seguridad y observabilidad**: Dependabot, CodeQL, UptimeRobot y Sentry

No se usa ORM: `sqlx` + `pgx` permiten escribir y mostrar el SQL real que se ejecuta contra la base, en vez de esconderlo detrás de una abstracción — decisión justificada en detalle en `decisiones.md`.

---

## Arquitectura

```
                React + Vite (SPA)
                       │
                       ▼
              Nginx (puerto 3000→80)
          sirve estáticos + proxy /api
                       │
                       ▼
              Backend Go + chi (8080)
           sqlx + pgx → SQL explícito
                       │
                       ▼
              PostgreSQL 16 (5432)
           volumen db_data (persistente)
```

El frontend usa rutas relativas en vez de una URL absoluta compilada (`VITE_API_URL`): nginx hace de proxy hacia el servicio `backend` en producción, y Vite hace lo mismo en desarrollo. Así la misma imagen del frontend funciona en cualquier entorno sin reconstruirse. El backend mantiene el middleware CORS igual, como capa de seguridad adicional.

---

## Estructura del repositorio

```
ingsoft3-tp01/
├── backend/
│   ├── Dockerfile              # multi-stage: golang:1.25-alpine → alpine:3.20
│   ├── go.mod / go.sum
│   ├── db/
│   │   └── init.sql            # creación de tablas + datos de ejemplo
│   └── ...                     # código Go (chi + sqlx + pgx) y tests
│
├── frontend/
│   ├── Dockerfile              # multi-stage: node:22-alpine → nginx:alpine
│   ├── nginx.conf
│   ├── package.json
│   ├── src/
│   └── e2e/                    # pruebas de integración y e2e (Playwright)
│
├── .github/
│   ├── dependabot.yml          # actualizaciones semanales: gomod, npm y github-actions
│   └── workflows/
│       └── ci.yml              # build, tests, integración, e2e y deploys
│
├── infra/                      # Terraform: entorno de preprod (Neon + Render)
│
├── docker-compose.yml          # build local
├── docker-compose.registry.yml # imágenes publicadas en GHCR
├── .env.example
├── decisiones.md               # decisiones técnicas de cada TP, con justificación
├── evidencias.md               # capturas y evidencia de cada TP
├── postmortem.md               # análisis de la caída provocada en producción (TP9)
└── README.md
```

---

## Ejecución local sin Docker

Requiere PostgreSQL 16 corriendo en el host (ver nota sobre el puerto 5432 en [Variables de entorno](#variables-de-entorno) si ya tenés una instancia nativa).

### Backend

```bash
cd backend
go mod download
go run .
# servidor disponible en http://localhost:8080
```

### Frontend

```bash
cd frontend
npm install
npm run dev
# frontend disponible en http://localhost:5173 (proxy a :8080 vía Vite)
```

---

## Reglas de negocio

Estas reglas viven en el **código Go del backend**, decisión intencional para que sean testeables unitariamente en TP5.

| # | Regla | Tipo |
|---|---|---|
| RN1 | No se pueden crear turnos en el pasado | Validación |
| RN2 | No se pueden crear turnos fuera del horario de atención | Validación |
| RN3 | No se pueden solapar turnos para el mismo profesional | Restricción |
| RN4 | No se pueden solapar turnos del mismo paciente entre distintos profesionales | Restricción |
| RN5 | No se puede cancelar un turno que ya empezó o que ya está completado | Restricción |
| RN6 | Transiciones manuales: `pendiente → confirmado`, `pendiente → cancelado` y `confirmado → cancelado`. Un turno pasa a `completado` sólo automáticamente, cuando estaba confirmado y su horario ya terminó | Transición de estado |
| RN7 | El precio del turno se copia de la tarifa del profesional al momento de la reserva (snapshot inmutable) | Cálculo |
| RN8 | Un turno se puede reservar hasta 10 minutos antes de su hora de inicio | Restricción |

---

## Endpoints disponibles

- `GET /health` — estado del servidor y de la conexión a la base.
- `GET /profesionales` (opcional `?especialidad=`) y `GET /profesionales/{id}`.
- `GET /profesionales/{id}/horarios-disponibles?fecha=YYYY-MM-DD` — huecos libres para reservar ese día.
- `POST /turnos` — reserva un turno (crea el paciente si no existe, identificándolo por DNI).
- `GET /turnos/{id}` y `GET /turnos?dni=...` / `?email=...` ("mis turnos").
- `PATCH /turnos/{id}/estado` — confirma o cancela un turno.

---

## Pantallas

- **Profesionales** (`/`) — listado filtrable por especialidad.
- **Detalle de profesional** (`/profesionales/:id`) — elegir fecha, ver horarios disponibles y reservar un turno.
- **Mis turnos** (`/mis-turnos`) — buscar turnos propios por DNI o email.
- **Detalle de turno** (`/turnos/:id`) — ver un turno y cambiar su estado según las transiciones permitidas.

---

## Tests

- **Backend**: `cd backend && go test ./... -cover`. Las pruebas se concentran en `internal/service`, donde viven las reglas de negocio; el umbral de cobertura del pipeline es 70 % de sentencias sobre ese paquete.
- **Frontend**: `cd frontend && npm test`. Cubre `src/utils` y `src/api/client.js`, con umbrales de 80 % de líneas y 75 % de ramas.
- **Integración y e2e**: las suites de Playwright están en `frontend/e2e/`. `api.spec.js` prueba la API directamente y `turnos.spec.js` recorre la aplicación en un navegador. En el pipeline corren contra el entorno de QA.

---

## Pipeline CI/CD

El workflow `.github/workflows/ci.yml` hace lo siguiente:

1. Construye backend y frontend en paralelo, con caché de capas de Docker, y corre los tests con sus umbrales de cobertura.
2. En `main`, publica las imágenes en GHCR con la etiqueta `sha-<commit>`.
3. Despliega automáticamente a **QA**.
4. Corre las pruebas de integración y e2e contra QA.
5. Si todo pasó, queda esperando la aprobación manual del environment `production` y recién entonces despliega a **PROD**.

Los despliegues a Render se hacen con deploy hooks que apuntan a la imagen publicada, con Auto-Deploy desactivado en los servicios. La rama `main` está protegida: exige PR y que pasen los checks requeridos. Para volver a una versión anterior de PROD se llama al deploy hook con el SHA del tag correspondiente.

---

## Entornos

| Entorno | Backend | Frontend | Base de datos | Cómo se crea |
|---|---|---|---|---|
| QA | Render | Render | Neon (`app_qa`) | A mano |
| PROD | Render | Render | Neon (`app_prod`) | A mano |
| Preprod | Render | Render | Neon | Terraform (`infra/`) |

Cada entorno usa su propia base de datos. Los servicios están en la región Ohio (AWS us-east-2) y usan el plan gratuito de Render, así que después de un tiempo sin uso se duermen y la primera carga puede tardar.

- Frontend de PROD: https://turnero-front-prod-pvlk.onrender.com
- API de PROD: https://turnero-api-prod-zwky.onrender.com/health

---

## Seguridad y observabilidad

- **Dependabot** revisa semanalmente las dependencias de Go, npm y GitHub Actions, y abre PRs de actualización.
- **CodeQL** analiza el código en cada cambio. Los workflows declaran permisos mínimos por job.
- **Push protection** de GitHub bloquea los pushes que contienen secretos.
- **UptimeRobot** consulta el `/health` de PROD cada 5 minutos y avisa por mail si cae.
- **Sentry** captura errores y tiempos de respuesta del backend y del frontend. No se envían query strings, cuerpos, cookies ni IP, porque `GET /turnos` lleva DNI y email en la URL.

---

## Variables de entorno

#### Backend

| Variable | Valor / default | Descripción |
|---|---|---|
| `DB_PASSWORD` | (definida en `.env`) | Contraseña de PostgreSQL, compartida entre `db` y `backend` (sólo en Docker Compose) |
| `DATABASE_URL` | `postgres://postgres:${DB_PASSWORD}@db:5432/turnos_centro_mujer?sslmode=disable` | Connection string armada por Compose |
| `FRONTEND_URL` | `http://localhost:3000` | Usado por el middleware CORS |
| `PORT` | `8080` | Puerto HTTP del backend |
| `SENTRY_DSN` | (vacía en local) | DSN del proyecto de Sentry del backend; se carga como variable de entorno del servicio en Render |

### Local (fuera de Docker)

Si corrés PostgreSQL en el host además del contenedor, recordá que **no pueden compartir el puerto 5432**: publicá el contenedor en otro puerto (por ejemplo 5433) y ajustá la connection string en consecuencia.

`.env` está en `.gitignore`; `.env.example` sí se versiona.

---

## Documentación adicional

- [`decisiones.md`](./decisiones.md) — decisiones técnicas y de diseño de cada TP, con problemas encontrados y justificaciones.
- [`evidencias.md`](./evidencias.md) — capturas y evidencia de cada TP.

---