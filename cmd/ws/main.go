package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/fquinteiro/ama-room-server/internal/api"
	"github.com/fquinteiro/ama-room-server/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WS_DATABASE_USER"),
		os.Getenv("WS_DATABASE_PASSWORD"),
		os.Getenv("WS_DATABASE_HOST"),
		os.Getenv("WS_DATABASE_PORT"),
		os.Getenv("WS_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
