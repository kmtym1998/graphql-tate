package test

import (
	"net"
	"testing"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/example/api"
)

func TestE2E(t *testing.T) {
	permission := tate.RootFieldPermission{}
	tate, err := tate.New(permission)
	if err != nil {
		t.Fatalf("failed to create tate: %v", err)
	}

	r := api.Router{
		Port: findAvailablePort(),
		Tate: tate,
	}
	go func() {
		r.ListenAndServe() // nolint:errcheck
	}()

	t.Skip("TODO: implement e2e test")
}

// 空いているポートを探して返す
func findAvailablePort() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	addr := l.Addr().String()

	return addr[len(addr)-4:]
}
