# Decisiones

## Enlaces de este TP (TP6)

- [Corrida de un PR](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38297699730/job/114972683928?pr=33): los tests en verde y "Entrar al registry" salteado.
- [Corrida de main](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38298297081/job/114974664060): "Construir y publicar la imagen" es el último paso.
- Commit del merge etiquetado en el registry: `5014b8c0f84c2512e36f33f9f5c012af72edb6c7`.
- QA: [front](https://turnero-front-qa-j4zw.onrender.com) y [API](https://turnero-api-qa-y8o0.onrender.com).
- PROD: [front](https://turnero-front-prod-pvlk.onrender.com) y [API](https://turnero-api-prod-zwky.onrender.com).

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

## TP5 — Testing y calidad

### Qué lógica elegí testear y por qué

Testeé la capa `internal/service` del backend, que es donde están todas las reglas de negocio del turnero:

- No reservar en el pasado ni con menos de 10 minutos de anticipación.
- Que el turno entre en el horario de atención del profesional.
- Que no se superponga con otro turno del profesional ni con otro de la paciente.
- La máquina de estados (qué cambios de estado se permiten).
- El auto-completado de los turnos confirmados que ya pasaron.
- El snapshot del precio al momento de reservar.
- Los horarios disponibles de cada profesional.

Es donde más duele un bug: un error ahí le da a una paciente un turno que no existe, pisa el turno de otra o le cobra un precio que no corresponde.

En el frontend testeé la lógica pura que tiene:

- `src/utils/format.js`: cómo se muestran las horas, los precios y las fechas, incluida la regla de no correr la hora 3 horas.
- `src/api/client.js`: cómo convierte los errores del backend en el mensaje que ve la paciente.
- `src/utils/resumen.js`: el texto de resumen de los turnos de la paciente, que agregué en el PR #29.

Al escribir los tests encontré que el código no cumplía dos reglas que yo misma había escrito en el README: RN5 decía que no se podía cancelar un turno pasado, pero no estaba implementado, y RN6 decía que un turno nunca pasaba a `completado` a mano, pero el backend lo permitía por la API. Las corregí, agregué la regla RN8 (anticipación mínima de 10 minutos) y actualicé el README. Es justo lo que plantea la guía: un test protege la regla tal como está en el código, y sólo yo sabía lo que se quería.

PRs: el refactor en [#22](https://github.com/canderojo/ingsoft3-tp01/pull/22), la suite del backend con estas correcciones en [#23](https://github.com/canderojo/ingsoft3-tp01/pull/23) y la del frontend en [#26](https://github.com/canderojo/ingsoft3-tp01/pull/26).

### El refactor para poder mockear

Antes del TP5, las funciones de `service` recibían `*sqlx.DB` y llamaban directo a las funciones de `repository`, y además usaban `time.Now()` adentro. No había forma de testearlas sin un Postgres levantado, y cualquier test que dependiera de la fecha iba a dar distinto según el día en que se corriera.

El refactor hizo que esas dos dependencias entren desde afuera: `service` ahora define una interfaz `Repositorio` con los 9 métodos que necesita, y `service.Turnos` se arma con `NuevoTurnos(repo, ahora)`. En la app real, `main.go` le pasa `repository.Postgres` (que sólo envuelve las funciones que ya existían) y el reloj; en los tests, un doble y una hora fija. Las reglas no cambiaron, sólo de dónde vienen la base y la hora. El paquete `service` ya ni siquiera importa a `repository`.

### Herramientas que usé (mi stack no es el de la cátedra)

| Qué hace falta | Backend (Go) | Frontend (React + Vite) |
|---|---|---|
| Test parametrizado | Tabla de casos con `t.Run` | `it.each` |
| El doble (mock/stub) | `repoDoble`, escrito a mano | `vi.fn()` + `vi.stubGlobal("fetch", …)` |
| Medir la cobertura | `go test -coverprofile` + `go tool cover` | `@vitest/coverage-v8` |
| Umbral que frena el build | Script `scripts/test-con-umbral.sh` | `thresholds` en `vite.config.js` |
| Qué entra en la cuenta | Sólo el paquete `./internal/service/...` | `include: ['src/utils/**', 'src/api/client.js']` |

En el backend, el doble (`repoDoble`) es una base de datos de mentira que escribí a mano: un struct que implementa la interfaz `Repositorio`. En Go es lo más común hacerlo así, y no hace falta instalar ninguna librería. Hace dos cosas:

- **Responde lo que el test le pide** (stub). Por ejemplo, "este profesional atiende de 9 a 13".
- **Anota lo que el servicio le pide a la base** (mock): qué turnos se intentaron guardar y qué cambios de estado se pidieron. Así el test puede revisar después qué hizo el servicio.

En el frontend, el mock está en el primer test de `client.test.js`. Ahí se reemplaza `fetch` (la función que le habla al backend) por una de mentira. El test revisa dos cosas: que el error del backend llegue bien a la pantalla, y que `fetch` se haya llamado una sola vez, con `POST`, a una dirección que termina en `/turnos`.

El umbral del backend va en un script (`scripts/test-con-umbral.sh`) porque Go no trae una opción para eso, como sí la tienen vitest o coverlet. El script corre los tests, lee el porcentaje total y, si queda abajo del umbral, hace fallar el build.


### Qué dejé afuera de la cuenta de cobertura

**Backend.** En Go no se puede excluir archivo por archivo como en .NET: lo que se elige es qué paquetes se miden. Elegí medir sólo `./internal/service/...`, porque ahí están todas las reglas de negocio. Esto tiene una ventaja: cualquier archivo nuevo que se agregue a `service` entra solo en la cuenta, sin tocar ninguna configuración. La desventaja es que si algún día creo un paquete nuevo con reglas, no se va a medir hasta que lo agregue a mano al comando.
Medí la cobertura de todo el backend con `go test ./... -coverprofile` y me dio 20,4 %. Es bajo porque cuenta paquetes que no tienen tests ni reglas de negocio, así que el umbral lo aplico sólo sobre `internal/service`. Dejé afuera:

- `main.go`, `internal/config` e `internal/db`: son el arranque. Leen las variables de entorno, abren la conexión a Postgres y conectan las piezas. No tienen reglas, y si algo ahí está mal la app no levanta.
- `internal/models`: son structs de datos. Las únicas funciones (`Scan`, `Value` y `MarshalJSON` de `HoraDelDia`) sólo traducen la hora entre Postgres y JSON.
- `internal/repository`: son las consultas SQL. Para probarlas de verdad hace falta un Postgres real, y eso ya es un test de integración (TP7)
- `internal/handlers`: reciben los pedidos que llegan desde el frontend, se los pasan al servicio y devuelven la respuesta. Los dejé afuera de la cuenta porque el umbral lo puse para controlar que las **reglas de negocio** estén testeadas, y los handlers no tienen reglas: sólo traducen entre HTTP y el servicio.
Un error en un handler y uno en el service no son igual de peligrosos. Si un handler falla, el pedido no funciona y se nota enseguida: la pantalla muestra un error. Si falla una regla del service, la app sigue funcionando, pero hace algo mal sin que nadie se dé cuenta, como dar un turno superpuesto o cobrar un precio equivocado. Por eso el esfuerzo de testear lo puse en el service.
Además, lo que hacen los handlers se prueba mejor con la app completa y una base de datos real, que es lo que se hace en el TP7. Si los contara ahora, sin tests, el porcentaje bajaría tanto que tendría que poner un umbral muy bajo, y dejaría de servir para controlar las reglas.
Igual no están vacíos: revisan que lleguen todos los datos y que la fecha tenga el formato correcto. Por eso lo dejo anotado como pendiente: para sumarlos habría que testearlos con `httptest`.

**Frontend.** Medí sólo `src/utils/**` y `src/api/client.js`, que es donde está la lógica. Dejé afuera:

- **Las páginas y los componentes:** muestran datos y reaccionan a clics. Se prueban mejor con la app completa en el navegador, que es lo que se hace en el TP7.
- **`mockData.js`:** son datos inventados para la demo, no hay nada que probar.
- **`api/turnos.js` y `api/profesionales.js`:** sólo arman la dirección del pedido y llaman a `client.js`, que sí está testeado.

### Umbral de cobertura

**Backend: 70 % de sentencias.**

Antes del PR #29, el pipeline medía el servicio en 78,6 % (77 de 98 sentencias cubiertas). El umbral lo elegí pensando en cuánto código sin tests tiene que entrar para que el build se frene:

| Umbral | Se frena cuando entran… | Problema |
|---|---|---|
| 80 % | 3 sentencias sin tests | Ya falla hoy casi sin margen: cualquier cambio lo rompe |
| 75 % | 5 sentencias sin tests | Un solo `if` con dos líneas lo rompe |
| **70 %** | **13 sentencias sin tests** | Frena una función mediana sin tests, pero deja hacer cambios chicos |

Me quedé con 70 % porque cumple lo que quiero del umbral: que no se pueda agregar una funcionalidad nueva sin tests, pero sin que el build falle por cualquier detalle. Un umbral que falla todo el tiempo termina desactivándose. Primero había elegido 75 %, pero lo había calculado con los números de mi máquina; cuando vi los del pipeline, me di cuenta de que dejaba muy poco margen y lo bajé.

Para subirlo habría que testear los handlers. Después del PR #29 el servicio mide 87,3 %.

La métrica es de sentencias porque es la única que mide Go: no tiene cobertura de rama, así que en el backend no hay número de rama para reportar. Para compensarlo, al escribir los tests pensé los dos caminos de cada `if` importante.

**Frontend: 80 % de líneas y 75 % de ramas.**

Antes del PR #29 medía 96,4 % de líneas (27 de 28) y 90,9 % de ramas (10 de 11). Con estos umbrales, el build se frena si entran **6 líneas** o **3 caminos de `if`** sin tests. Si los subía a 90 % y 85 %, una sola función nueva de dos líneas ya lo rompía.

Uso dos números porque en el frontend sí se mide la cobertura de rama, y es la más honesta: dice si se probaron los dos lados de cada `if`, no sólo si se ejecutó la línea. El de ramas es más bajo porque son pocas y cada una mueve mucho el porcentaje.

Después del PR #29 mide 97,4 % de líneas y 96 % de ramas.

La corrida con el resumen de cobertura de los dos lados y los reportes para descargar: [corrida 36939280662](https://github.com/canderojo/ingsoft3-tp01/actions/runs/36939280662).

### Por qué cobertura alta no garantiza calidad

La cobertura mide qué código se ejecutó, no qué se comprobó. Un test que llame a `dentroDeHorarioAtencion(inicio, fin, profesional)` y no revise el resultado suma cobertura, pero no detecta nada: si la función dijera que un turno de las 3 de la mañana entra en horario, el test seguiría en verde.

### Ejercicio del camino sin cubrir

En el reporte HTML (`go tool cover -html`) quedó en rojo esta parte de `ObtenerTurno`:

```go
if err == sql.ErrNoRows {
    return nil, ErrTurnoNoExiste
}
```
`ObtenerTurno` busca un turno por su número. Este `if` es el camino de cuando el turno **no existe**. Los tests que había sólo probaban turnos que sí existían, porque siempre le cargaban un turno al doble. El otro camino, el del turno inexistente, no lo probaba nadie.

**Qué entrada lo recorre.** Pedir un turno que no existe, por ejemplo el turno 99, con el doble vacío (sin ningún turno cargado). Como no encuentra nada, el doble contesta `sql.ErrNoRows`, que es lo mismo que contesta Postgres cuando busca una fila que no está.

**Qué decidí.** Agregar el test `TestObtenerTurno_QueNoExiste_DevuelveErrTurnoNoExiste`, en el PR [#24](https://github.com/canderojo/ingsoft3-tp01/pull/24). Lo agregué porque esa línea decide qué ve la paciente cuando entra a un turno que no existe: gracias a ella la app contesta "el turno no existe" (error 404). Si alguien la borrara, la app contestaría "error interno del servidor" (error 500), y antes ningún test lo hubiera detectado.

### El umbral bloqueando un merge

Para demostrar que el umbral frena de verdad hice dos PRs.

**Primer PR: [#29](https://github.com/canderojo/ingsoft3-tp01/pull/29) (mergeado).**

Agregué una funcionalidad nueva: un resumen de los turnos de la paciente, que dice cuál es su próximo turno, cuántos tiene activos y cuántos le falta confirmar.

1. **Primer commit, sin tests.** Todo compilaba y todos los tests pasaban, pero los dos checks se pusieron en rojo porque la cobertura quedó abajo del umbral. Como los dos son required checks desde el TP4, GitHub no dejó mergear. 
Se ve en los logs de la [corrida roja](https://github.com/canderojo/ingsoft3-tp01/actions/runs/36919025218).
2. **Segundo commit, con los tests.** Escribí un test por cada caso que contempla el código nuevo, para que no quedara ningún camino sin probar.

   En el backend, tres tests para `ResumenDeTurnos`:
   - **Una paciente sin turnos:** el resumen no tiene próximo turno y marca 0 activos.
   - **Una paciente con turnos de todo tipo:** le cargué cinco turnos. Tres no tienen que contar (uno cancelado, uno completado y uno pendiente que ya empezó) y dos sí (un confirmado de mañana y un pendiente de hoy a la tarde). El test comprueba que cuente sólo esos dos, que marque uno sin confirmar, que elija como próximo el de hoy porque es el más cercano, y que avise que tiene un turno hoy.
   - **Una paciente sin turnos hoy:** su único turno es mañana, así que el resumen no tiene que avisar que tiene turno hoy.

   En el frontend, un solo test con `it.each` para `textoResumenTurnos`, con un caso por cada mensaje que puede mostrar: sin turnos (lista vacía o sin lista), sólo turnos cancelados o completados, un turno activo, varios activos con uno sin confirmar, y varios sin confirmar. Así se prueban tanto los textos en singular como en plural.

   Con esos tests la cobertura volvió a pasar el umbral, los dos checks se pusieron en verde y pude mergear.

| | Umbral | Sin tests (rojo) | Con tests (verde) |
|---|---|---|---|
| Backend (sentencias) | 70 % | 65,3 % | 87,3 % |
| Frontend (líneas) | 80 % | 71,05 % | 97,4 % |
| Frontend (ramas) | 75 % | 40 % | 96 % |

**Segundo PR: [#30](https://github.com/canderojo/ingsoft3-tp01/pull/30) (queda abierto hasta la defensa).**

Agrega validaciones de DNI y email para el formulario de reserva, sin tests. Compila y los tests pasan, pero `build-frontend` está en rojo: la cobertura cae a 71,15 % de líneas y 55,81 % de ramas, abajo de los umbrales de 80 % y 75 %. `build-backend` queda en verde, pero alcanza con que uno de los dos esté en rojo para que no se pueda mergear.

Lo hice sólo en el frontend porque el backend, después del PR #29, quedó en 87,3 %. Para bajarlo del 70 % hubiera tenido que agregar unas 30 líneas de código sin tests, y en el frontend alcanzaba con un archivo chico. 
Este PR no lo arreglo: queda así, en rojo, para mostrar en la defensa que el freno funciona.

### Problemas encontrados y cómo los resolví

- **Bug de zona horaria.** Los turnos de la tarde del día aparecían todos tachados. La app guarda los horarios con la hora de Argentina pero etiquetada como UTC, y el backend los comparaba con `time.Now()` del contenedor, que es la hora UTC real, 3 horas adelantada. Lo arreglé pasándole al servicio un reloj en hora de Argentina (`service.AhoraEnArgentina`, con UTC-3 fijo) desde `main.go`. Gracias al reloj inyectado no tuve que tocar ninguna regla.
- **La cobertura del backend no daba igual en mi máquina y en el pipeline.** Con el mismo código y los mismos tests, en local daba 82,4 % y en el contenedor 78,6 %. La diferencia es la versión de Go: tengo instalada la 1.27 y el Dockerfile usa la 1.25, y las dos cuentan distinto las sentencias de `HorariosDisponibles` y `CrearTurno`. Tomé como válido el número del pipeline, porque es donde mide el gate y la 1.25 es la versión del proyecto, y recalculé el umbral con esos números.
- **Un test del front dependía de mi `.env`.** El test del mock comparaba la URL con `"/turnos"`, pero en mi máquina el `.env` tenía `VITE_API_URL=http://localhost:8080`, así que fallaba en local y en el pipeline no. Cambié el assert para que sólo compruebe que la URL termina en `/turnos`, y saqué la variable del `.env` y del `.env.example`, porque ya no se usaba (el proxy de Vite hace ese trabajo)

### Declaración de uso de IA

Usé Claude en el chat para entender la teoría de la guía (AAA, stubs y mocks, cobertura de línea y de rama, quality gates), para adaptar cada paso a Go y a mi app, y para escribir los tests, que me mostró y explicó antes de aplicarlos. También la use para redactar esta sección.


## TP6: entrega continua con environments, aprobaciones y deployment patterns

## Enlaces de este TP (TP6)

- Paquetes públicos en ghcr.io, etiquetados con el commit del merge `5014b8c0f84c2512e36f33f9f5c012af72edb6c7` (se bajan con `docker pull`, sin credenciales): [backend](https://github.com/users/canderojo/packages/container/package/ingsoft3-tp01-backend) y [frontend](https://github.com/users/canderojo/packages/container/package/ingsoft3-tp01-frontend).
- [Corrida de un PR](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38297699730/job/114972683928?pr=33) (PR #33): los tests en verde y "Entrar al registry" salteado.
- [Corrida de main](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38298297081/job/114974664060): "Construir y publicar la imagen" es el último paso.
- - [Corrida aprobada](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38332145805) (#47: QA, espera de revisión, aprobación y PROD en verde).
- [Corrida rechazada](https://github.com/canderojo/ingsoft3-tp01/actions/runs/38334184933) (#49): el rechazo con su motivo.
- [Release `v6.0.0`](https://github.com/canderojo/ingsoft3-tp01/releases/tag/v6.0.0).
- QA: [front](https://turnero-front-qa-j4zw.onrender.com) y [API](https://turnero-api-qa-y8o0.onrender.com).
- PROD: [front](https://turnero-front-prod-pvlk.onrender.com) y [API](https://turnero-api-prod-zwky.onrender.com).

### 1. Por qué la imagen se publica solo si la verificación pasó

La imagen se publica en ghcr.io únicamente si la verificación pasó, y eso lo garantiza el orden de los pasos: "Construir y publicar la imagen" es el último paso de cada job, después de construir la etapa de tests, correrla y armar el reporte de cobertura. Los pasos de un job corren en orden y se cortan en el primer error, así que si los tests fallan el job termina ahí y la publicación nunca se ejecuta. No hace falta ningún `if` que diga "si los tests pasaron": la condición está en el orden. 

La garantía sale de encadenar tres cosas. Nada entra a `main` sin los checks en verde (la protección de rama del TP4 más el umbral de cobertura del TP5). Solo lo que entra a `main` se publica: el paso de publicar lleva `push` condicionado a que el evento sea un push a `main`. Y publicar es el último paso del mismo job que corrió los tests. Lo comprobé con dos corridas: en la del PR #33 los tests pasan y "Entrar al registry" sale salteado, y en la corrida de `main` "Construir y publicar la imagen" es el último paso. Las imágenes quedan como `ghcr.io/canderojo/ingsoft3-tp01-backend` y `ghcr.io/canderojo/ingsoft3-tp01-frontend`, etiquetadas con el commit que las produjo (`sha-<commit>`).

Esto garantiza lo que publica el pipeline, no lo que físicamente puede entrar al registry: con mis credenciales podría hacer un `docker push` a mano desde mi máquina y nada lo impediría. La cadena evita que algo llegue por accidente, pero no es un control de acceso.

### 2. Continuous Delivery y no Continuous Deployment

Implementé Continuous Delivery: lo que entra a `main` con el CI en verde se despliega solo a QA, y producción queda esperando la aprobación manual del environment `production`. No implementé Continuous Deployment porque hoy no tengo la red de seguridad que haría falta para sacar a la persona de ese paso. Mis tests cubren la lógica de turnos del backend y las utilidades del frontend, pero no hay tests de integración ni end-to-end contra el entorno desplegado, no tengo monitoreo ni alertas, y el rollback es manual. Si todo lo que pasa los tests llegara solo a producción, un error que los tests no ven llegaría directo a las pacientes sin que nadie lo mire. Para una app de turnos de un centro de salud me parece correcto que una persona autorice. Para pasar a Continuous Deployment me faltarían tests contra el entorno desplegado, monitoreo con alertas y un rollback automático.

### 3. Proveedor y entornos

Usé Render para correr la app y Neon para las bases de datos, y con eso cumplo los cinco puntos del contrato del TP. Hay dos entornos separados, QA y PROD, cada uno con su front y su API como servicios web de Render con URL pública (cuatro servicios en total, los cuatro en Ohio). Cada entorno tiene su propia base en Neon, `app_qa` y `app_prod`, dentro del mismo proyecto. Front y back corren como contenedores construidos a partir de mis Dockerfiles (runtime Docker). El deploy lo dispara el pipeline con el commit que verificó: el Auto-Deploy de Render está apagado en los cuatro servicios y el workflow llama a los deploy hooks. Y producción está detrás de la aprobación del environment `production` de GitHub.

Elegí Neon y no el Postgres de Render porque el Postgres gratuito de Render expira a los 30 días de creado, que ni siquiera cubre el mes del TP, mientras que el plan gratuito de Neon es permanente y no pide tarjeta. Es una decisión de leer la letra chica del free tier. Puse Neon y Render en la misma región (Ohio) para que cada consulta de la API a la base no cruce de región.

Mi backend no crea las tablas al arrancar y no uso migraciones: en local las crea el `init.sql` que monta el compose. En Neon ese archivo no corre solo, así que ejecuté el contenido de `backend/db/init.sql` a mano en el SQL Editor, primero en `app_qa` y después en `app_prod`, y verifiqué con `select count(*) from profesionales` que las dos devolvieran 8. No es la forma profesional de hacerlo, que sería una herramienta de migraciones como Flyway o golang-migrate, y por eso lo cuento como una limitación.

### 4. Qué quedó en variables y qué quedó adentro de la imagen

La dirección del backend no está escrita adentro de la imagen del front. El `nginx.conf` del TP2 tenía `http://backend:8080` fijo (el nombre del servicio en compose), que en Render no existe. Lo convertí en una plantilla (`frontend/default.conf.template`) donde el DNS y la dirección del backend salen de las variables `DNS_RESOLVER` y `BACKEND_URL`, que la imagen oficial de nginx reemplaza al arrancar. El Dockerfile deja como valores por defecto los de compose (`127.0.0.11` y `http://backend:8080`), así que `docker compose up` sigue funcionando sin tocar nada, y Render los pisa con las variables de cada servicio. Lo comprobé en mi máquina levantando la misma imagen con y sin variables. En Render, el front de QA y el de PROD se construyen del mismo código y solo cambia `BACKEND_URL`.

Por variable quedó también la conexión a la base: cada API recibe `DATABASE_URL` con la URI de Postgres de su base de Neon, cargada en el servicio de Render y no en el repositorio. El puerto lo define Render con `PORT` y `config.go` lo lee. Adentro de la imagen quedó la estructura de nginx: las tres rutas que se reenvían al backend (`/profesionales`, `/turnos` y `/health`) y la línea `proxy_ssl_server_name on`, necesaria porque la API de Render es https.

Para cada punto de la tabla de equivalencias de la guía, esto es lo que usé:

| Qué tenía que lograr | Qué usé |
|---|---|
| Publicar la imagen desde el pipeline | `docker/build-push-action` hacia ghcr.io |
| Que solo la rama principal publique | `push` condicionado a `github.event_name == 'push' && github.ref == 'refs/heads/main'` |
| Etiquetar con el commit | `sha-${{ github.sha }}` |
| Un entorno con su propia configuración | environments `qa` y `production` de GitHub |
| Que el deploy lo dispare el pipeline | deploy hooks de Render, con Auto-Deploy apagado |
| Desplegar el commit verificado y no la punta de la rama | `&ref=$GITHUB_SHA` en cada hook |
| Un gate humano antes de producción | required reviewers del environment `production` |
| Secrets con alcance de entorno | environment secrets |
| Una base por entorno | dos bases en Neon, `app_qa` y `app_prod` |
| Crear el esquema en cada base | `init.sql` ejecutado a mano en el SQL Editor |
| Smoke test con reintentos | bucle de `curl` en bash |
| Release versionada con notas | tag `v6.0.0` y release de GitHub con notas generadas |
| Pasarle la conexión a la app | variable `DATABASE_URL` (URI de Postgres, driver pgx) |
| Endpoint que diga que está vivo | `/health` |
| Endpoint que pruebe que la base responde | `/profesionales` |
| Servir el front construido | imagen de nginx con el `dist` de Vite |
| Que el front sepa dónde está su backend | variable `BACKEND_URL` en la plantilla de nginx |

### 5. El diseño de la cadena y el alcance de cada secret

El `ci.yml` tiene cuatro jobs. `build-backend` y `build-frontend` verifican y publican la imagen. `deploy-qa` tiene `needs: [build-backend, build-frontend]`, así que no arranca si alguno falló, y `if: github.ref == 'refs/heads/main'`, así que en un Pull Request queda salteado: los PRs verifican pero no despliegan. `deploy-prod` tiene `needs: deploy-qa`, de modo que solo corre si QA quedó vivo. No repite el `if` de la rama porque lo hereda: si `deploy-qa` se saltea, `deploy-prod` también.

Cada job de deploy declara su environment. `qa` no tiene reglas de protección, porque QA es automático a propósito. `production` tiene required reviewers (yo misma) con "Prevent self-review" desactivado, que es lo que me permite aprobar mis propios deploys. Los secrets de cada entorno viven dentro de su environment: `RENDER_HOOK_API_QA` y `RENDER_HOOK_FRONT_QA` en `qa`, y `RENDER_HOOK_API_PROD` y `RENDER_HOOK_FRONT_PROD` en `production`. Los de PROD tienen que estar ahí y no a nivel repositorio porque así solo los puede leer un job que pasó la aprobación: si estuvieran en el repo, cualquier job del workflow podría usarlos y desplegar a producción sin pasar por el gate. Para publicar las imágenes no hace falta ningún secret nuevo, porque uso el `GITHUB_TOKEN` de cada corrida, que dura lo que dura el job. El permiso `packages: write` lo declaro dentro de cada job que publica y no a nivel workflow, para darle el mínimo privilegio solo a quien lo necesita.

Un deploy hook es una URL de Render que, cuando el pipeline la llama, le ordena a Render construir y publicar ese servicio. Si se llama sola, Render despliega lo último que haya en `main` en ese momento, y eso no necesariamente es lo que el pipeline verificó. Por eso a cada llamada le agrego `&ref=$GITHUB_SHA`, que le indica a Render qué commit tiene que desplegar: el de la corrida que acaba de pasar los tests. Sin el `ref`, con dos merges seguidos la corrida del primero podría terminar desplegando el segundo, que todavía no pasó por ninguna verificación, y la idea de que se promueve lo mismo que se verificó dejaría de ser cierta. En producción esto pesa todavía más, porque la corrida puede quedar esperando mi aprobación durante horas y `main` puede cambiar en ese tiempo: yo aprobaría el commit A y podría terminar desplegándose otro. Esto también es lo que hace que, aunque Render reconstruya la imagen y no use la que publiqué, al menos el código que se construye sea el verificado. Lo comprobé con el rollback: `main` estaba en `289d62c` (el commit del footer), llamé a los hooks de PROD con el SHA del tag `v6.0.0` y Render desplegó `d6cb538` en los dos servicios, con disparador "Deploy Hook". Desplegó el commit que le indiqué y no el último de `main`.

El job `deploy-prod` también tiene `concurrency` con `cancel-in-progress: false`. Sirve para que no haya dos deploys a producción corriendo al mismo tiempo y para que, si ya hay uno en curso, no se lo cancele a la mitad. Lo que no resuelve es el orden de las aprobaciones. Si mergeo un cambio A y después un cambio B, las dos corridas quedan esperando mi aprobación, y como B se hizo sobre A, incluye los dos cambios. Si apruebo primero B y después A, producción termina en el código de A y pierde lo de B, aunque nadie lo haya borrado. La guía aclara que la cátedra no midió cómo se comporta `concurrency` con jobs detenidos esperando aprobación. Por eso la regla es mía y manual: si hay dos esperando, apruebo en orden de la más vieja a la más nueva, o rechazo la vieja si la nueva ya la incluye.

### 6. Qué mira el aprobador y el rechazo documentado

Antes de aprobar un deploy a producción miro que la corrida tenga los dos builds y `deploy-qa` en verde, y abro el front de QA para comprobar que el cambio ya se ve ahí: un verde del smoke solo no me alcanza. Miro qué cambia ese commit, en particular si toca la conexión a la base, las variables o el workflow. Compruebo que no haya otra corrida esperando aprobación, para no aprobar una vieja después de una nueva. Y cuando el cambio es de código, miro en Render que el deploy de QA de ese commit figure como live.

Hay cosas que mi aprobador no puede ver. El smoke de QA no dice qué versión está corriendo, así que un verde puede ser de la versión anterior mientras Render todavía construye. Tampoco veo cómo se comporta la app con usuarias reales, porque no tengo monitoreo ni métricas de errores. Ese es el tipo de información que me faltaría para decidir mejor.

El flujo completo quedó en la corrida #47 (el PR #35 con el job `deploy-prod`): QA verde, "Waiting for review", aprobación y PROD en verde con su smoke. El rechazo es la corrida #49, que salió del PR #36, un commit que solo agregaba la sección de enlaces a `decisiones.md`. El motivo que escribí al rechazar fue: "Este commit solo agrega enlaces a decisiones.md y no cambia el código: PROD ya corre este mismo código desde la corrida #47, así que no hay nada nuevo para llevar. Rechazo y espero el próximo cambio real de código". Es un motivo que sale de lo que vi en esa corrida: ese commit no toca código y producción ya corría el mismo código desde la corrida anterior, así que no había nada nuevo para llevar. `deploy-prod` quedó en rojo y producción no se tocó.

### 7. La letra chica del free tier y cómo la maneja el pipeline

En Render, un servicio gratuito se duerme tras unos 15 minutos sin tráfico y el primer pedido lo despierta con un cold start de hasta cerca de un minuto. Por eso el smoke test no es un `curl` único: reintenta hasta 30 veces cada 20 segundos, 10 minutos de margen que cubren el build de Render más el cold start, y cada `curl` lleva `--max-time 10` para que un servicio que acepta la conexión y no contesta no deje colgado el job. Las horas de instancia y los minutos de build son del workspace y no de cada servicio: 750 horas de instancia y 500 minutos de build por mes, compartidos por los cuatro servicios. Para no gastarlas no dejé ningún ping que mantenga los servicios despiertos. 

Neon suspende el cómputo tras unos 5 minutos sin uso y se despierta solo, mucho más rápido que Render; el plan gratuito limita el almacenamiento a 0,5 GB y las horas de cómputo mensuales. Que la base pueda estar despertando cuando el backend arranca fue una de las razones para dos cambios en el backend antes de empezar el despliegue. `db.Connect` usa `sqlx.Open` en vez de `sqlx.Connect`, así que no prueba la conexión al arrancar y el backend se recupera solo cuando la base responde. Y `/health` devuelve 503 cuando la base no responde, en vez de 200 con el error escondido en el cuerpo, para que el smoke test pueda confiar en el código de estado.

### 8. Qué garantía pierdo porque Render reconstruye desde el repositorio

La imagen que el pipeline publica en ghcr.io es la que pasó los tests, pero no es la que corre en QA ni en PROD. Render construye las dos imágenes de nuevo desde el repositorio a partir del commit que le indica el hook. Entonces promuevo el mismo commit, pero no la misma imagen. Dos construcciones del mismo código pueden salir distintas: una dependencia que cambió entre un build y el otro, una imagen base que se actualizó. Lo que corre en producción podría no ser exactamente lo que verificaron mis tests, y si algo saliera distinto sería difícil de explicar. Lo hice así porque en este TP Render despliega desde el repositorio; hacer que ejecute la imagen publicada en el registry es lo que cierra el TP siguiente, y la imagen etiquetada con el commit que ya publico es lo que lo hace posible. El `&ref` reduce el riesgo, porque al menos el commit desplegado es el verificado.

### 9. Qué prueba mi smoke test y qué no

En cada intento el smoke test comprueba tres cosas. `/health` de la API responde con éxito, y como el handler devuelve 503 si la base no responde, un éxito implica que la API está viva y llegó a la base. `/profesionales` de la API devuelve datos, lo que prueba que la base tiene las tablas y el backend consulta. Y `/profesionales` pedido a través del front llega al backend, lo que prueba que el proxy de nginx y `BACKEND_URL` están bien. Esa tercera comprobación la agregué yo: la guía solo pide que el front abra, y con una dirección del backend mal escrita el front abre igual y esa prueba daría verde.

Lo que no prueba es qué versión está corriendo. El hook responde al instante y Render construye en segundo plano; mientras construye, o si ese build falla, sigue sirviendo la versión anterior, así que el smoke puede dar verde contra la versión vieja. La mitigación sería que `/health` informe el commit que corre y que el smoke lo compare con `github.sha`, no la implementé porque la guía la marca como opcional en este TP y obligatoria en el siguiente, y prefirí no mezclar un cambio del backend con el resto de la entrega. Lo dejo como una limitación conocida de mi pipeline: hoy un verde del smoke dice que el entorno responde, pero no que esté corriendo el commit que acabo de desplegar.

### 10. Cómo comprobé que cada entorno usa su propia base

El error que quería descartar es que PROD estuviera apuntando a la base de QA. Las dos cadenas de conexión difieren en una sola palabra (`app_qa` contra `app_prod`), se copian una después de la otra, y si pego la equivocada no se pone nada en rojo: PROD levanta, contesta y el smoke da verde, pero escribe en la base de QA. Que `/health` y `/profesionales` den lo mismo en los dos entornos no prueba nada, porque las dos bases tienen los mismos ocho profesionales. Para comprobarlo inserté un profesional con un nombre inconfundible solo en `app_prod` desde el SQL Editor de Neon ("Dra. Paula PROD"). Apareció en la API y en el front de PROD y no apareció en los de QA. Eso muestra que cada API lee su propia base. Dejé esa fila en `app_prod` para mostrarla en la defensa.

### 11. Patrón de deployment para una producción real, y plan de rollback

Para esta app en una producción real con usuarios usaría blue-green. Mantendría dos entornos completos, uno con la versión actual y otro con la nueva, y un router que cambia el tráfico de uno al otro cuando el nuevo está probado. Lo elegiría porque el rollback es instantáneo (alcanza con devolver el cambio del router), no hay tiempo de caída durante el despliegue, y en un sistema de turnos de salud pesa más poder volver atrás rápido que el costo de la infraestructura. 
El costo es el doble de infraestructura y el riesgo principal es que las dos versiones comparten la misma base de datos, así que un cambio de esquema que la versión vieja no entienda complica el rollback. No elegiría canary porque necesita métricas de tráfico real para decidir si el canario está fallando, y esa observabilidad hoy no la tengo. Lo que me falta para ejecutar cualquiera de los patrones avanzados es monitoreo: métricas de errores y latencia, alertas y logs centralizados. Hoy hago un recreate, con una ventana en la que el servicio se reconstruye.

Mi plan de rollback actual, que probé de verdad, es redesplegar la versión anterior conocida. Primero obtengo el commit del tag con `git rev-list -n1 v6.0.0`. Después llamo a los dos deploy hooks de PROD con `&ref=<ese commit>`; los hooks los copio de Render. Espero a que Render muestre ese commit como live en los dos servicios y compruebo en el front que volvió la versión anterior.

Lo medí así: con producción en el commit del footer (PR #38), llamé a los hooks con el commit del tag `v6.0.0` (`d6cb538`) a las 00:20:30 y el rollback tardó 45 segundos hasta que los dos servicios figuraron como live. En Render, los deploys duraron 18,6 segundos (API) y 28,3 segundos (front). Después del rollback el footer de PROD volvió a decir la frase anterior. Esta vía llama a los hooks directamente, así que no pasa por la aprobación del environment `production` ni deja una corrida en Actions; la evidencia es la pantalla Deploys de Render.

El rollback de código no deshace los datos. La fila "Dra. Paula PROD" sigue en `app_prod` aunque haya vuelto a un commit anterior, y lo mismo pasaría con una migración que hubiera borrado una columna: el código vuelve atrás y la base no.

### 12. Problemas encontrados y cómo los resolví

El primer `turnero-api-qa` quedó en Singapur en vez de Ohio; la región no se puede cambiar, así que lo borré y lo recreé. 
`/health` contestaba siempre 200 aunque la base estuviera caída, lo que habría dejado pasar el smoke con la base rota; lo cambié a 503 y de paso usé `sqlx.Open` para que el backend no quede roto si la base despierta tarde. Al abrir el front de QA con el backend dormido vi un 502: era el cold start de la API, y se resolvió despertándola primero.

Sobre el orden de la evidencia: la guía indica hacer primero el rechazo y después mergear otro Pull Request para aprobar. Yo lo hice en otro orden. Primero aprobé la corrida #47, porque era la primera vez que probaba el camino completo a producción y llevaba el cambio del nginx del front. Después hice el rechazo en la corrida #49, y recién después aprobé otro deploy con un cambio visible, el del subtítulo de la pantalla principal (PR #37, commit `d6cb538`), que es el que etiqueté como `v6.0.0` y que se ve en la pantalla de PROD. Más adelante aprobé también el cambio de la frase del footer (PR #38) para poder practicar el rollback. Así que hay una aprobación posterior al rechazo, como pide la guía, y las dos evidencias obligatorias existen, cada una en su propia corrida. La única diferencia con lo que sugiere la guía es que la primera aprobación (#47) fue anterior al rechazo.

### 13. Declaración de uso de IA

Usé Claude para entender la guía antes de aplicarla, para armar los archivos que tenia que modificar en cada paso y para redactar un borrador de esta sección, que después revisé y corregí.