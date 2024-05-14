package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/test/server/api/handler"
	"github.com/kmtym1998/graphql-tate/test/server/api/middleware"
)

type Router struct {
	Port string
	Tate *tate.Tate
}

func (r Router) ListenAndServe() error {
	mux := chi.NewRouter()
	mux.Use(middleware.InjectRole)

	mux.Post("/v1/graphql", handler.PostV1GraphQLHandler(r.Tate))

	port := func() string {
		if r.Port == "" {
			return "8080"
		}

		return r.Port
	}()

	slog.Info(
		"Server is running",
		slog.String("port", port),
		"address", "http://localhost:"+port+"/v1/graphql",
	)

	return http.ListenAndServe(":"+port, mux)
}
