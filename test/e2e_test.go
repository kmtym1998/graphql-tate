package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"

	tate "github.com/kmtym1998/graphql-tate"
	"github.com/kmtym1998/graphql-tate/example/api"
	"github.com/kmtym1998/graphql-tate/example/api/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

type testCase struct {
	name     string
	query    string
	role     string
	expected string
}

func TestE2E(t *testing.T) {
	permission := tate.RootFieldPermission{
		ast.Query: tate.ChildFieldPermission{
			"user": tate.ChildFieldPermission{
				"id":   tate.OR(middleware.IsEditor, middleware.IsAdmin),
				"name": tate.OR(middleware.IsViewer, middleware.IsEditor, middleware.IsAdmin),
			},
			"todos": tate.OR(middleware.IsViewer, middleware.IsEditor, middleware.IsAdmin),
		},
		ast.Mutation: tate.ChildFieldPermission{
			"createTodo": middleware.IsAnonymous,
		},
	}

	t.Run("with default tate", func(t *testing.T) {
		tate, err := tate.New(permission)
		if err != nil {
			t.Fatalf("failed to create tate: %v", err)
		}

		port := findAvailablePort()

		r := api.Router{
			Port: port,
			Tate: tate,
		}
		go func() {
			r.ListenAndServe() // nolint:errcheck
		}()

		for _, tc := range []testCase{
			{
				name:     "query user as admin",
				query:    `query { user(id: "U1") { id name } }`,
				role:     "admin",
				expected: `{"data":{"user":{"id":"U1","name":"user1"}}}`,
			}, {
				name:     "query user as viewer",
				query:    `query { user(id: "U1") { id name } }`,
				role:     "viewer",
				expected: `{"errors":[{"message":"permission denied for user","path":["user","id"],"extensions":{"fieldName":"user"}}],"data":null}`,
			}, {
				name:     "query todos as editor",
				query:    `query { todos { id text } }`,
				role:     "editor",
				expected: `{"data":{"todos":[{"id":"1","text":"todo1"},{"id":"2","text":"todo2"},{"id":"3","text":"todo2"}]}}`,
			}, {
				name:     "query todos as anonymous",
				query:    `query { todos { id text } }`,
				role:     "anonymous",
				expected: `{"errors":[{"message":"permission denied for todos","path":["todos"],"extensions":{"fieldName":"todos"}}],"data":null}`,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				reqBody, err := json.Marshal(map[string]interface{}{
					"query": tc.query,
				})
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}

				req, err := http.NewRequest(http.MethodPost, "http://localhost:"+port+"/v1/graphql", bytes.NewBuffer(reqBody))
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
				req.Header = http.Header{
					"Content-Type": []string{"application/json"},
					"X-Role":       []string{tc.role},
				}

				t.Logf("request to %s", "http://localhost:"+port+"/v1/graphql")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("failed to send request: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					t.Errorf("unexpected status code: %d", resp.StatusCode)
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}

				if string(respBody) != tc.expected {
					t.Errorf("unexpected response body:\n\texpected: %s\n\tactual: %s", tc.expected, respBody)
				}
			})
		}
	})
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
