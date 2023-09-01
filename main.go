package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/madraceee/url-shortner/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

	// DB connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB URL not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handleHealth)
	v1Router.Post("/shortURL", apiCfg.handleUrlShorten)
	v1Router.Get("/{shortURL}", apiCfg.handleFetchShortUrl)

	router.Mount("/v1", v1Router)

	// Server frontend
	fs := http.FileServer(http.Dir("frontend"))
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Detect file extension to set the correct content type
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}
		fs.ServeHTTP(w, r)
	}))

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server is starting at port %v", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Could not start server ", err)
	}
}
