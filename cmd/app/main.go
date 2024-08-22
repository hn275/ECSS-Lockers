package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/zvdv/ECSS-Lockers/internal"
	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/email"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router"
	"github.com/zvdv/ECSS-Lockers/internal/router/auth"
	"github.com/zvdv/ECSS-Lockers/internal/router/dash"
)

const addr string = "127.0.0.1:8080"

func init() {
	godotenv.Load()

	dbURL := fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN"))
	database.Connect(dbURL)
	internal.Initialize()
	email.Initialize()
    crypto.Initialize()
}

func main() {
	app := chi.NewRouter()

	app.Use(middleware.RealIP)
	app.Use(requestLogger)
	app.Use(middleware.Recoverer)

	app.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	app.Handle("/", http.HandlerFunc(router.Home))

	app.Route("/auth", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(auth.Auth))
		r.Handle("/api/login", http.HandlerFunc(auth.AuthApiLogin))
		r.Handle("/api/token", http.HandlerFunc(auth.AuthApiToken))
	})

	app.Route("/dash", func(r chi.Router) {
		r.Use(auth.AuthenticatedUserOnly)
		r.Handle("/", http.HandlerFunc(dash.Dash))
		r.Handle("/locker/register", http.HandlerFunc(dash.DashLockerRegister))
		r.Handle("/api/locker", http.HandlerFunc(dash.ApiLocker))
	})

	logger.Info("Listening at http://%s", addr)
	logger.Info("for local dev, use http://127.0.0.1:8080, for more information, see: https://stackoverflow.com/a/1188145/19114163")
	logger.Fatal(http.ListenAndServe(addr, app))
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		next.ServeHTTP(ww, r)
		end := time.Since(start)

		url := r.URL.Path

		statusCode := ww.Status()
		statusString := color.GreenString("%d", statusCode)
		if statusCode >= 400 {
			statusString = color.RedString("%d", statusCode)
		}

		logger.Trace("%s %s %s from %s - %s %dB in %v",
			r.Method,
			url,
			r.Proto,
			r.RemoteAddr,
			statusString,
			ww.BytesWritten(),
			end)
	})
}