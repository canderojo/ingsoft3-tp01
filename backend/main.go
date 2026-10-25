package main

import (
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/config"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/db"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/handlers"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/repository"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/service"
)

func main() {
	// godotenv.Load lee el archivo .env de esta carpeta y carga sus
	// variables como si fueran variables de entorno del sistema. Si no
	// existe .env (por ejemplo en un futuro deploy), no falla: seguimos
	// con las variables de entorno que ya estén seteadas afuera.
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, se usan variables de entorno del sistema")
	}

	cfg := config.Load()

	// Sentry recibe los errores (panics) y mide cuánto tarda cada pedido.
	// Con SENTRY_DSN vacía el SDK queda apagado: no manda nada y la app
	// funciona igual.
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0, // medir todos los pedidos
		// Centro de salud: nada de datos personales de pacientes en Sentry.
		DataCollection: &sentry.DataCollection{
			UserInfo:   sentry.Set(false),   // ni IP ni usuario
			HTTPBodies: []sentry.BodyType{}, // nunca cuerpos (el alta de turno trae DNI y email)
			Cookies:    &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
			// GET /turnos?dni=...&email=... lleva datos del paciente en el query string.
			QueryParams: &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
			// Render agrega headers con la IP del cliente (X-Forwarded-For, True-Client-IP).
			HTTPHeaders: &sentry.HeaderCollectionConfig{
				Request:  &sentry.KeyValueCollectionBehavior{Terms: headersConIP},
				Response: &sentry.KeyValueCollectionBehavior{Terms: headersConIP},
			},
		},
	}); err != nil {
		log.Printf("Advertencia: no se pudo iniciar Sentry (%v)", err)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Printf("Advertencia: no se pudo conectar a la base de datos (%v)", err)
		log.Println("El servidor levanta igual; /health va a marcar el error hasta que Postgres esté disponible")
	}

	// chi.NewRouter crea un router: el objeto que decide qué función de
	// Go se ejecuta según el método HTTP (GET, POST, ...) y la ruta
	// (/health, /turnos, ...) del request que llega.
	router := chi.NewRouter()
	router.Use(middleware.Logger) // loguea cada request en la consola (método, ruta, tiempo de respuesta)

	// Sin esto, el navegador bloquea los pedidos que el frontend (en
	// otro puerto, el dev server de Vite) le hace a esta API.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{cfg.FrontendURL},
		AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	// Sentry va después de CORS (los preflight OPTIONS no generan
	// transacciones) y antes de las rutas. Repanic: después de reportar el
	// panic lo vuelve a lanzar, para que lo siga manejando quien esté afuera.
	// Cada transacción sale con el patrón de chi (GET /turnos/{id}): chi
	// completa r.Pattern al rutear y sentryhttp lo usa como nombre.
	router.Use(sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle)

	healthHandler := handlers.HealthHandler{DB: database}
	router.Get("/health", healthHandler.Health)
	router.Head("/health", healthHandler.Health) // los monitores externos (UptimeRobot) preguntan con HEAD

	// El servicio de turnos recibe la base (envuelta en repository.Postgres)
	// y el reloj de Argentina. En los tests se reemplazan ambos por dobles.
	turnosService := service.NuevoTurnos(repository.Postgres{DB: database}, service.AhoraEnArgentina)

	profesionalesHandler := handlers.ProfesionalesHandler{DB: database, Servicio: turnosService}
	router.Get("/profesionales", profesionalesHandler.Listar)
	router.Get("/profesionales/{id}", profesionalesHandler.Obtener)
	router.Get("/profesionales/{id}/horarios-disponibles", profesionalesHandler.HorariosDisponibles)

	turnosHandler := handlers.TurnosHandler{DB: database, Servicio: turnosService}
	router.Post("/turnos", turnosHandler.Crear)
	router.Get("/turnos", turnosHandler.ListarDePaciente)
	router.Get("/turnos/{id}", turnosHandler.Obtener)
	router.Patch("/turnos/{id}/estado", turnosHandler.CambiarEstado)

	log.Printf("Servidor escuchando en http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		sentry.Flush(2 * time.Second) // log.Fatal sale con os.Exit, sin esperar a Sentry
		log.Fatal(err)
	}
}

// headersConIP son partes de nombres de headers que llevan la IP del
// cliente; Sentry reemplaza su valor por "[Filtered]".
var headersConIP = []string{"forwarded", "-ip", "remote-", "via", "-user"}
