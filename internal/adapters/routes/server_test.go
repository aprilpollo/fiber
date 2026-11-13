package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/routes/handler"
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// stubUsers is an input.UserService whose answers each test sets, so the HTTP
// adapter can be checked on its own: binding, status mapping, response shape.
type stubUsers struct {
	created input.CreateUser
	err     error
}

func (s *stubUsers) List(context.Context, domain.ListQuery) ([]domain.User, int64, error) {
	return nil, 0, s.err
}

func (s *stubUsers) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &domain.User{ID: id, Name: "Alice", Email: "alice@example.com"}, nil
}

func (s *stubUsers) Create(_ context.Context, in input.CreateUser) (*domain.User, error) {
	s.created = in
	if s.err != nil {
		return nil, s.err
	}
	return domain.NewUser(in.Name, in.Email)
}

func (s *stubUsers) Update(context.Context, uuid.UUID, input.UpdateUser) (*domain.User, error) {
	return nil, s.err
}

func (s *stubUsers) Delete(context.Context, uuid.UUID) error { return s.err }

func testApp(users input.UserService) *fiber.App {
	cfg := &config.App{AppName: "test", RateLimitMax: 1000, RateLimitWindowSeconds: 60}
	return NewApp(cfg, handler.NewUserHandler(users), handler.NewBookHandler(nil),
		func(context.Context) error { return nil })
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
	Payload json.RawMessage `json:"payload"`
}

func do(t *testing.T, app *fiber.App, method, target, body string) (int, envelope) {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	}
	res, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	var env envelope
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &env); err != nil {
			t.Fatalf("%s %s: body is not the JSON envelope: %s", method, target, raw)
		}
	}
	return res.StatusCode, env
}

func TestCreateIgnoresFieldsTheClientMayNotSet(t *testing.T) {
	forged := uuid.New()
	status, env := do(t, testApp(&stubUsers{}), "POST", "/api/v1/users",
		`{"id":"`+forged.String()+`","created_at":"2000-01-01T00:00:00Z","name":"Alice","email":"alice@example.com"}`)
	if status != fiber.StatusCreated {
		t.Fatalf("status %d: %+v", status, env)
	}
	var got struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
	_ = json.Unmarshal(env.Payload, &got)
	if got.ID == forged || got.CreatedAt.Year() == 2000 {
		t.Fatalf("client-supplied id/created_at reached the entity: %+v", got)
	}
}

func TestDomainErrorsMapToStatus(t *testing.T) {
	id := uuid.New().String()
	tests := []struct {
		name, method, target, body string
		err                        error
		status                     int
		message                    string
	}{
		{"not found", "GET", "/api/v1/users/" + id, "", domain.ErrNotFound, 404, "NOT_FOUND"},
		{"conflict", "POST", "/api/v1/users", `{"name":"A","email":"a@b.co"}`, domain.ErrConflict, 409, "CONFLICT"},
		{"domain validation", "PUT", "/api/v1/users/" + id, `{"name":""}`, &domain.ValidationError{Field: "name", Message: "is required"}, 400, "BAD_REQUEST"},
		{"unknown filter", "GET", "/api/v1/users?password_eq=x", "", &domain.ValidationError{Field: "password", Message: "cannot be filtered on"}, 400, "BAD_REQUEST"},
		{"delete missing", "DELETE", "/api/v1/users/" + id, "", domain.ErrNotFound, 404, "NOT_FOUND"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, env := do(t, testApp(&stubUsers{err: tc.err}), tc.method, tc.target, tc.body)
			if status != tc.status || env.Message != tc.message {
				t.Fatalf("got %d %+v, want %d %s", status, env, tc.status, tc.message)
			}
		})
	}
}

func TestInternalErrorsAreNotShownToTheClient(t *testing.T) {
	leak := errors.New(`pq: relation "users" does not exist`)
	status, env := do(t, testApp(&stubUsers{err: leak}), "GET", "/api/v1/users/"+uuid.New().String(), "")
	if status != fiber.StatusInternalServerError {
		t.Fatalf("status %d", status)
	}
	if strings.Contains(env.Error, "relation") {
		t.Fatalf("storage error leaked to the client: %q", env.Error)
	}
}

func TestRequestShapeIsCheckedBeforeTheCore(t *testing.T) {
	stub := &stubUsers{}
	app := testApp(stub)
	for _, tc := range []struct{ method, target, body string }{
		{"GET", "/api/v1/users/not-a-uuid", ""},
		{"POST", "/api/v1/users", `{"name":"Alice","email":"nope"}`},
		{"POST", "/api/v1/users", `{"name":`},
		{"GET", "/api/v1/users?_limit=abc", ""},
	} {
		if status, env := do(t, app, tc.method, tc.target, tc.body); status != 400 {
			t.Errorf("%s %s: got %d %+v, want 400", tc.method, tc.target, status, env)
		}
	}
	if stub.created != (input.CreateUser{}) {
		t.Fatalf("an invalid request reached the service: %+v", stub.created)
	}
}
