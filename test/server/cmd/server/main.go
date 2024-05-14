package main

import (
	"log/slog"
	"os"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/test/server/api"
	"github.com/kmtym1998/graphql-tate/test/server/api/middleware"
	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(
		os.Stdout,
		&tint.Options{
			AddSource: true,
		},
	)))

	permission := middleware.NewPermission()
	tate, err := tate.New(permission)
	if err != nil {
		slog.Error(
			"failed to create tate",
			slog.String("error", err.Error()),
		)
	}

	router := api.Router{
		Tate: tate,
		Port: os.Getenv("PORT"),
	}
	if err := router.ListenAndServe(); err != nil {
		slog.Error(
			"failed to listen and serve",
			slog.String("error", err.Error()),
		)
	}
}
