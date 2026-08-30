# Proyecto ingsoft3-tp01 - Version B

git clone https://github.com/canderojo/ingsoft3-tp01.git

## Turnero Centro de Salud de la Mujer

Sistema de reserva de turnos para un centro de salud con profesionales de 4 especialidades: dermatología, nutrición, ecografía y endocrinología.

Proyecto para la materia **Ingeniería del Software 3** (UCC 2026). Cada trabajo práctico del cuatrimestre agrega una capa nueva sobre esta misma aplicación (CI, testing, CD, contenedores, IaC, seguridad, observabilidad).

## Stack

- **Backend**: Go + [chi] (router) + [sqlx] + [pgx] (driver de Postgres)
- **Frontend**: React + Vite
- **Base de datos**: PostgreSQL


## Cómo levantar el proyecto con Docker

**Requisitos**: Docker Desktop.

1. Cloná el repositorio:
```bash
   git clone https://github.com/canderojo/ingsoft3-tp01.git
   cd ingsoft3-tp01
```
2. Copiá la plantilla de variables de entorno y completá la contraseña de la base:
```bash
   cp .env.example .env
```

3. Levantá el sistema completo (compila las imágenes desde el código):
```bash
   docker compose up -d --build
```

4. Abrí `http://localhost:3000` en el navegador.

**Levantar desde las imágenes ya publicadas** (sin compilar, usando lo publicado en GitHub Container Registry):
```bash
docker compose -f docker-compose.registry.yml up -d
```
(también necesita el `.env` con `DB_PASSWORD` del paso 2)

**Verificar que el backend está sano:**
```bash
curl http://localhost:8080/health
```
Debería devolver `{"status":"ok","db":"ok"}`.

**Apagar el sistema:**
```bash
docker compose down        # conserva los datos de la base
docker compose down -v     # borra también el volumen de la base
```

## Endpoints disponibles

- `GET /health` — estado del servidor y de la conexión a la base.
- `GET /profesionales` (opcional `?especialidad=`) y `GET /profesionales/{id}`.
- `GET /profesionales/{id}/horarios-disponibles?fecha=YYYY-MM-DD` — huecos libres para reservar ese día.
- `POST /turnos` — reserva un turno (crea el paciente si no existe, identificándolo por DNI).
- `GET /turnos/{id}` y `GET /turnos?dni=...` / `?email=...` ("mis turnos").
- `PATCH /turnos/{id}/estado` — confirma o cancela un turno.

## Pantallas

- **Profesionales** (`/`) — listado filtrable por especialidad.
- **Detalle de profesional** (`/profesionales/:id`) — elegir fecha, ver horarios disponibles y reservar un turno.
- **Mis turnos** (`/mis-turnos`) — buscar turnos propios por DNI o email.
- **Detalle de turno** (`/turnos/:id`) — ver un turno y cambiar su estado según las transiciones permitidas.

## Documentación

- [decisiones.md](./decisiones.md) — decisiones técnicas y de diseño, con justificaciones.
- [evidencias.md](./evidencias.md) — capturas y evidencia de cada TP.

