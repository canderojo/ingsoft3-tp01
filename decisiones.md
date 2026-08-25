# Decisiones

## 1. Por qué Git no pudo resolver el conflicto automáticamente

Git no pudo resolver el conflicto automáticamente porque las ramas `feature/titulo-a` y `feature/titulo-b` modificaron el mismo contenido del archivo de maneras diferentes. Git no podía determinar por sí solo cuál de las dos versiones era la correcta.

El conflicto podría haberse evitado si las dos ramas no hubieran modificado las mismas líneas del archivo, o si los cambios se hubieran integrado antes de que las ramas se alejaran demasiado de `main`.

Para resolverlo, revisé las dos versiones, elegí qué contenido debía quedar y eliminé los marcadores `<<<<<<<`, `=======` y `>>>>>>>`.

## 2. Problemas encontrados y cómo los solucioné

El primer problema que surgio fue el intento de hacer `push` directamente a `main`. GitHub lo rechazó porque la rama estaba protegida, por lo que los cambios debían hacerse mediante un Pull Request. Esto permitió comprobar que la protección estaba funcionando correctamente.

Finalmente, al trabajar con dos ramas se produjo un conflicto. Lo resolví desde GitHub revisando las diferencias entre ambas versiones y seleccionando manualmente el contenido que debía conservarse.

## 3. Declaración de uso de IA

Utilicé inteligencia artificial como herramienta de apoyo durante la realización del TP, principalmente para escribir los contenidos de los archivos y avanzar de manera guiada en la configuración del repositorio, ramas, Pull Requests, resolución de conflictos y versionado.

No utilicé la IA para reemplazar la ejecución del trabajo. Los comandos y procedimientos fueron ejecutados y comprobados en mi propio entorno.

Cuando recibí indicaciones de la IA, las contrasté con el enunciado del TP y comprobé que los cambios realizados en el repositorio fueran los esperados.


---

## TP2 — Contenedores

### 1. Qué app elegí y por qué

La app que elegí para el semestre es un turnero para un centro de salud de la mujer: profesionales de 4 especialidades (dermatología, nutrición, ecografía, endocrinología) que atienden turnos reservados por pacientes. La desarrollé específicamente para esta materia.

Contra los criterios de la guía: corre localmente sin problemas (Go + chi + sqlx/pgx en el backend, React + Vite en el frontend, contra PostgreSQL); tiene reglas de negocio no triviales y testeables unitariamente (validación de horarios, superposición de turnos, transiciones de estado, snapshot de precio al momento de reservar); la entiendo en profundidad porque es desarrollo propio; y el tamaño (4 pantallas) está dentro del rango de "CRUD + 2-3 pantallas" que pide la consigna.

Elegí Go porque es un lenguaje que ya habia utiizado anteriormente, por lo que tengo mayor manejo y comprension del mismo. Para el router usé chi en vez de un framework más grande como Gin: es minimalista y cercano a la librería estándar, más fácil de explicar en la defensa oral. En vez de un ORM usé sqlx + pgx, para poder mostrar y justificar el SQL real que se ejecuta contra la base, en vez de una abstracción que lo esconda.

React + Vite para el frontend por ser el stack SPA más estándar. 

PostgreSQL porque lo sugiere la cátedra.

### 2. Decisiones de contenerización

Backend y frontend usan Dockerfiles multi-stage: una etapa de build con el compilador, y una etapa final que solo contiene lo necesario para ejecutar.

**Imágenes base:**

| Servicio | Etapa build | Etapa final |
|----------|-------------|-------------|
| Backend | `golang:1.25-alpine` | `alpine:3.20` |
| Frontend | `node:22-alpine` | `nginx:alpine` |
| Base de datos | — | `postgres:16-alpine` (oficial, sin build propio) |

Como Go compila a un binario nativo, la imagen final del backend no necesita ningún runtime instalado — solo el binario y el sistema base mínimo. Resultado: 35.7MB, contra los 329MB que pesa la imagen que incluye el compilador. La imagen del frontend, con nginx sirviendo los estáticos ya compilados, queda en 93.7MB.

En ambos Dockerfiles copio primero los archivos de dependencias (`go.mod`/`go.sum` o `package.json`) e instalo antes de copiar el resto del código, para aprovechar el cache de capas de Docker: si solo cambia código, no se vuelven a descargar dependencias.

**Qué persiste y qué no:** los datos de PostgreSQL viven en un volumen nombrado (`db_data:/var/lib/postgresql/data`), que sobrevive a un `docker compose down` sin `-v`. Lo verifiqué creando un turno, reiniciando el sistema completo, y confirmando que seguía presente; con `down -v` en cambio el volumen se borra intencionalmente (verificado: la consulta de turnos vuelve vacía). Los contenedores de `backend` y `frontend` no guardan estado — son recreables sin pérdida de información.

A diferencia del sample de la cátedra (que crea las tablas con código al arrancar la app), mi backend no lo hace. Uso un bind mount (`./backend/db/init.sql:/docker-entrypoint-initdb.d/init.sql`) para que PostgreSQL ejecute el script de creación de tablas y datos de ejemplo la primera vez que inicializa el volumen.

**Rutas relativas + proxy nginx en vez de URL absoluta + CORS**: el frontend originalmente llamaba al backend con una URL fija (`VITE_API_URL`), compilada dentro del JavaScript en el momento del build. Lo migré a rutas relativas, con nginx haciendo de proxy hacia el servicio `backend` en producción (y el proxy de Vite en desarrollo), para que la misma imagen del frontend sirva en cualquier entorno sin necesitar reconstruirse. El backend mantiene el middleware CORS configurado de todas formas, como capa de seguridad adicional.

### 3. Problemas encontrados y cómo los resolví

- **Puerto 5432 ocupado por un PostgreSQL nativo de Windows**: tenía instalado PostgreSQL como servicio nativo, escuchando en el mismo puerto que el contenedor. Lo resolví publicando el contenedor de Postgres en el puerto 5433 y ajustando la connection string (`Host=localhost;Port=5433;...` durante las pruebas locales, `Host=db` dentro del compose).
- **Puerto 8080 ocupado por el propio backend corriendo suelto**: al intentar correr el contenedor del backend mientras el binario corría en paralelo fuera de Docker (para pruebas), hubo conflicto de puertos. Lo identifiqué con `Get-NetTCPConnection -LocalPort 8080` y lo resolví deteniendo el proceso con `Stop-Process`.

### 4. Declaración de uso de IA

Utilicé inteligencia artificial (Claude) como herramienta de apoyo durante la realización del TP2, principalmente para entender los conceptos de Docker (imágenes, contenedores, multi-stage builds, redes, volúmenes) y para guiarme paso a paso en la escritura de los Dockerfiles, el `docker-compose.yml, la publicación en el registry y la redaccion de los archivos decisiones.md, README.md y evidencias.md.

No utilicé la IA para reemplazar la ejecución del trabajo. Todos los comandos fueron ejecutados y verificados en mi propio entorno (builds, levantado de contenedores, pruebas de persistencia, publicación en ghcr.io). Practiqué primero el flujo completo sobre el sample de la cátedra (`demo-fullstack`) antes de aplicarlo a mi propia app, para entender cada paso antes de repetirlo.


