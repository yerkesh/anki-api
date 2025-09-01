package httpServer

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"anki-api/internal/handlers"
)

// Deps wraps things the router needs (db pool, etc.).
type Deps struct {
	DB *pgxpool.Pool
	H  *handlers.Handler
}

func New(deps Deps) chi.Router {
	r := chi.NewRouter()

	// Chi goodies
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Liveness / readiness
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := deps.DB.Ping(ctx); err != nil {
			http.Error(w, "db down", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ok"))
	})

	r.Post("/user", deps.H.CreateUser)
	r.Post("/{userID}/collection", deps.H.CreateCollection)
	r.Post("/{collectionID}/flashcard", deps.H.CreateFlashcard)

	r.Get("/{userID}/collections", deps.H.GetCollections)
	r.Get("/{collectionID}/flashcard", deps.H.GetFlashcards)
	r.Patch("/flashcards/{flashcardID}/status", deps.H.UpdateFlashcardStatus)
	r.Delete("/flashcards/{flashcardID}", deps.H.DeleteFlashcard)
	r.Get("/flashcards/{flashcardID}", deps.H.GetFlashcard)
	// Versioned API
	r.Mount("/api/v1", r)

	return r
}
