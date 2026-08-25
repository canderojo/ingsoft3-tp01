package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/canderojo/turnos-centro-mujer/backend/internal/config"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/db"
	"github.com/canderojo/turnos-centro-mujer/backend/internal/handlers"
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

	healthHandler := handlers.HealthHandler{DB: database}
	router.Get("/health", healthHandler.Health)

	profesionalesHandler := handlers.ProfesionalesHandler{DB: database}
	router.Get("/profesionales", profesionalesHandler.Listar)
	router.Get("/profesionales/{id}", profesionalesHandler.Obtener)
	router.Get("/profesionales/{id}/horarios-disponibles", profesionalesHandler.HorariosDisponibles)

	turnosHandler := handlers.TurnosHandler{DB: database}
	router.Post("/turnos", turnosHandler.Crear)
	router.Get("/turnos", turnosHandler.ListarDePaciente)
	router.Get("/turnos/{id}", turnosHandler.Obtener)
	router.Patch("/turnos/{id}/estado", turnosHandler.CambiarEstado)

	log.Printf("Servidor escuchando en http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
