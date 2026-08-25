// Package config centraliza la lectura de variables de entorno.
// Idea: nada de valores hardcodeados en el resto del código — todo lo
// que pueda cambiar entre tu máquina, otra máquina, o un futuro
// contenedor Docker, se lee de acá.
package config

import "os"

// En Go, un "struct" es como una clase sin métodos ni herencia: solo
// agrupa datos relacionados bajo un mismo nombre.
type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
}

// Load arma un Config leyendo variables de entorno, con valores por
// defecto razonables para desarrollo local.
func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		// Origen permitido para CORS: el dev server de Vite usa el
		// puerto 5173 por defecto.
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

// getEnv es una función chica que envuelve os.LookupEnv (de la
// librería estándar de Go) para poder darle un valor por defecto
// cuando la variable no está seteada.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
