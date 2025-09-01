package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	nethttp "net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"

	"anki-api/config"
	"anki-api/internal/clients/chatgpt"
	"anki-api/internal/handlers"
	httpRouter "anki-api/internal/httpServer"
	"anki-api/internal/repository"
	"anki-api/internal/repository/generated"
	"anki-api/internal/usecase"
)

func main() {
	//// optional: allow custom .env path
	//var envFile string
	//flag.StringVar(&envFile, "env-file", ".env", "path to env file")
	//flag.Parse()
	//
	//// Load .env (ignore missing by default, fail if user explicitly provided a path)
	//if envFile != "" {
	//	if err := godotenv.Load(envFile); err != nil {
	//		log.Fatalf("load .env: %v", err)
	//	}
	//} else {
	//	// tries ".env" in CWD; ignore error if not present
	//	_ = godotenv.Load()
	//}

	var cfg config.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}

	if err := run(cfg); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run(cfg config.Config) error {
	ctx := context.Background()

	pool, err := newPgxPool(ctx, cfg.DatabaseConfig)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	log.Printf("connected to %s", cfg.DatabaseConfig.Name)

	queries := generated.New(pool)
	gptClient := chatgpt.New(cfg.ChatGPTKey, time.Second*20)
	collectionsRepo := repository.NewCollectionsRepo(queries)
	h := handlers.NewHandler(handlers.NewOptions(
		handlers.WithValidator(validator.New()),
		handlers.WithUsers(usecase.NewUsersUsecase(repository.NewUsersRepo(queries))),
		handlers.WithCollections(usecase.NewCollectionsUsecase(collectionsRepo)),
		handlers.WithFlashcards(usecase.NewFlashcardsUsecase(repository.NewFlashcardsRepo(queries), collectionsRepo, gptClient)),
	))

	r := httpRouter.New(httpRouter.Deps{DB: pool, H: h})
	srv := httpRouter.NewServer(fmt.Sprintf(":%s", cfg.HTTP.Port), r)

	ln, err := net.Listen("tcp", srv.Addr())
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("HTTP server listening on %s", ln.Addr())

	go func() {
		if err = srv.Listen(ln); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-ctx.Done()

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(shCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

	log.Println("shutting down")
	return nil
}

func newPgxPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	return pool, nil
}
