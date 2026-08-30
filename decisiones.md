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


## TP3 — Planificación y trazabilidad (GitHub Projects)

### Duración del sprint

Elegí sprints de **1 semana**, porque la cátedra entrega un TP por semana. Cada sprint cierra
alineado con cada entrega, lo que me permite tener siempre un objetivo claro y verificable al
final de la semana, en vez de un ciclo desacoplado del calendario real de la materia.

### Límite de trabajo en progreso (WIP)

Elegí un límite de **2** para la columna "In Progress" (regla de arranque: personas + 1,
trabajando sola = 2). Esto me permite tener una tarea activa y una segunda en caso de que la
primera quede esperando algo (una revisión, una respuesta) sin perder el foco de avanzar. Si en
la práctica nunca lo alcanzo, es señal de que quedó demasiado alto y debería bajarlo a 1.

### Diagnóstico de la historia mal escrita

La historia "Como desarrollador quiero crear la tabla usuarios" está mal escrita porque es una
**tarea técnica disfrazada de historia de usuario**: nadie "quiere" una tabla, es un medio para lograr algo, no algo que alguien pueda notar o valorar. Le falta el "para qué..." que justifica por
qué importa hacerla. La reescribiría como: "Como paciente quiero registrar mis datos para poder reservar un turno" y "crear la tabla usuarios" pasaría a ser una de las tareas técnicas dentro de esa historia.

### Problemas encontrados y cómo los resolví

No encontré problemas durante el desarrollo del TP.

### Declaración de uso de IA

Usé Claude como asistente durante todo el TP: para entender la teoría (jerarquía épica/historia/tarea, INVEST, criterios de aceptación, WIP limit, trazabilidad) antes de ejecutar cada paso; para guiarme comando a comando en la creación del proyecto, las labels, los issues, el board, el sprint y el PR de trazabilidad. Verifiqué cada paso ejecutándolo yo misma en mi cuenta de GitHub y revisando el resultado real en la web antes de avanzar al siguiente paso; y para la redacción de decisiones.md

## TP4 — CI: Pipelines as Code (GitHub Actions)

### Estructura del pipeline

Elegí dos jobs en paralelo (`build-backend` y `build-frontend`), uno por cada Dockerfile del TP2. Los separé porque backend y frontend tienen stacks completamente distintos (Go vs Node/nginx) y ciclos de vida propios: un cambio en el frontend no debería esperar a que compile el backend para saber si está bien, y viceversa. Como corren en runners independientes, un job puede fallar sin bloquear al otro — lo verifiqué directamente en la demo de §3.4: al romper el backend, `build-frontend` siguió en verde.

El pipeline no compila por su cuenta: usa `docker/build-push-action` apuntando a los mismos Dockerfiles del TP2 (`context: ./backend`, `context: ./frontend`), con `push: false`. La decisión de fondo es no duplicar la definición de build: si el pipeline compilara con `go build`/`npm run build` directamente, tendría dos formas distintas de construir la app que tarde o temprano divergen, y estaría verificando algo distinto de lo que después se despliega. El Dockerfile es la única fuente de verdad sobre cómo se arma la app.

### Cache de capas

Cacheo las capas de Docker con `type=gha` (el almacén de GitHub Actions), usando `docker/setup-buildx-action` en los dos jobs — el driver de Docker que viene por defecto en los runners no sabe exportar cache, así que sin ese paso el build falla directamente. Cada job usa un `scope` distinto (`backend` / `frontend`) para no pisarse el cache entre sí.

Lo confirmé corriendo el pipeline dos veces sobre el mismo PR: en la segunda corrida, todas las capas del backend y del frontend salieron `CACHED`. El tiempo bajó de 45s/59s a 19s/18s en backend/frontend respectivamente; aunque, como marca la guía, la ganancia de tiempo no es la evidencia real (con un proyecto de este tamaño puede no notarse, o incluso ser más lento por el costo de subir el cache); lo que importa es que las capas se reutilizaron.

Si el cache desaparece (GitHub lo desaloja guardando espacio o simplemente expira), el pipeline sigue funcionando exactamente igual, solo que reconstruye todo desde cero, no es una dependencia, es una optimización.

### El pipeline como gate

Activé `required_status_checks` sobre `main` exigiendo `build-backend` y `build-frontend`, con `strict: true` y 0 approvals requeridos, ya que el trabajo es individual y GitHub nunca deja aprobar el propio PR.

Lo demostré rompiendo el backend a propósito: agregué un import a un paquete inexistente (`_ "github.com/canderojo/paquete-que-no-existe"`), verifiqué que fallaba también en local con `docker build ./backend`, y lo subí en un PR (#17). El check `build-backend` se puso en rojo, el botón de merge quedó bloqueado, y `build-frontend` siguió en verde de forma independiente. Saqué el import roto en un commit de fix, el pipeline volvió a correr solo y se puso en verde, y el botón de merge se habilitó. 
Para verificar el `strict: true` (que exige tener la rama actualizada, no solo el check en verde) necesitaba dos PRs abiertos al mismo tiempo. Usé el PR #16 (un cambio mínimo en el README, abierto previamente para el checkpoint de §3.3) como ese segundo PR: tras mergear la demo del gate (#17), volví al #16 y apareció el botón "Update branch", confirmando que quedó desactualizado contra el nuevo `main`, y que GitHub exige traer los cambios antes de permitir el merge. Con la evidencia ya capturada, cerré el PR #16 sin mergearlo porque su único propósito era servir de "PR de relleno" para esta prueba, su contenido no aportaba nada al proyecto en sí, y la guía habilita explícitamente cualquiera de las dos opciones (mergear o cerrar) para este caso puntual.

### Problemas encontrados y cómo los resolví

- **VS Code borraba el import "roto" al guardar**: el formateador automático de Go (`goimports`) detecta imports no usados y los elimina al guardar. Lo resolví agregando el prefijo `_` antes de la ruta del import (`_ "paquete"`), que le indica a Go que el import es intencional aunque no se use en el código, y el formateador dejó de tocarlo.

### Declaración de uso de IA

Usé Claude como asistente durante todo el TP: para entender la teoría (integración continua como práctica, pipeline as code, anatomía de un workflow, triggers, cache de capas, secrets, el pipeline como gate) antes de ejecutar cada paso; y para guiarme paso a paso en la escritura del `ci.yml`, la configuración del gate vía GitHub Settings, y la demostración de la rotura/fix del build. Verifiqué cada paso ejecutándolo yo misma (corridas del pipeline, builds locales con `docker build` para confirmar que los errores eran reales y no artificios del pipeline, configuración de la protección de rama) antes de avanzar al siguiente. Usé Claude también para la redacción de esta sección de `decisiones.md`.