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