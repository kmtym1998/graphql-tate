package middleware

import (
	"context"
	"net/http"
)

func InjectRole(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roleName := r.Header.Get("X-Role")
		if roleName == "" {
			roleName = "anonymous"
		}

		ctx = RoleWith(ctx, roleName)

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

type RoleKey struct{}

func RoleWith(ctx context.Context, roleName string) context.Context {
	return context.WithValue(ctx, RoleKey{}, roleName)
}

func RoleFrom(ctx context.Context) string {
	roleName, _ := ctx.Value(RoleKey{}).(string)
	return roleName
}
