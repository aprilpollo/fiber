package services

import (
	"context"
	"errors"
	"testing"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"

	"github.com/google/uuid"
)

// fakeUserRepo is an in-memory output.UserRepository — the point of the port
// is that the core can be exercised without a database.
type fakeUserRepo struct {
	rows    map[uuid.UUID]domain.User
	updates int
}

func newFakeUserRepo(users ...domain.User) *fakeUserRepo {
	r := &fakeUserRepo{rows: map[uuid.UUID]domain.User{}}
	for _, u := range users {
		r.rows[u.ID] = u
	}
	return r
}

func (r *fakeUserRepo) FindAll(context.Context, domain.ListQuery) ([]domain.User, int64, error) {
	out := make([]domain.User, 0, len(r.rows))
	for _, u := range r.rows {
		out = append(out, u)
	}
	return out, int64(len(out)), nil
}

func (r *fakeUserRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.rows[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &u, nil
}

func (r *fakeUserRepo) Save(_ context.Context, u *domain.User) error {
	for _, existing := range r.rows {
		if existing.Email == u.Email {
			return domain.ErrConflict
		}
	}
	r.rows[u.ID] = *u
	return nil
}

func (r *fakeUserRepo) Update(_ context.Context, u *domain.User) error {
	if _, ok := r.rows[u.ID]; !ok {
		return domain.ErrNotFound
	}
	r.updates++
	r.rows[u.ID] = *u
	return nil
}

func (r *fakeUserRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := r.rows[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.rows, id)
	return nil
}

func ptr(s string) *string { return &s }

func TestCreateUserStoresAValidUser(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	u, err := svc.Create(context.Background(), input.CreateUser{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.rows[u.ID]; !ok {
		t.Fatal("created user was not saved")
	}
}

func TestCreateUserRejectsInvalidInputBeforeStorage(t *testing.T) {
	repo := newFakeUserRepo()
	_, err := NewUserService(repo).Create(context.Background(), input.CreateUser{Name: "Alice", Email: "nope"})

	var verr *domain.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("want ValidationError, got %v", err)
	}
	if len(repo.rows) != 0 {
		t.Fatal("invalid user reached the repository")
	}
}

func TestCreateUserPassesConflictThrough(t *testing.T) {
	existing, _ := domain.NewUser("Alice", "alice@example.com")
	svc := NewUserService(newFakeUserRepo(*existing))

	_, err := svc.Create(context.Background(), input.CreateUser{Name: "Other", Email: "alice@example.com"})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("want ErrConflict, got %v", err)
	}
}

func TestUpdateUserChangesOnlyGivenFields(t *testing.T) {
	existing, _ := domain.NewUser("Alice", "alice@example.com")
	repo := newFakeUserRepo(*existing)

	u, err := NewUserService(repo).Update(context.Background(), existing.ID, input.UpdateUser{Name: ptr(" Alicia ")})
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alicia" || u.Email != "alice@example.com" {
		t.Fatalf("got %+v", u)
	}
}

func TestUpdateUserValidatesTheResult(t *testing.T) {
	existing, _ := domain.NewUser("Alice", "alice@example.com")
	repo := newFakeUserRepo(*existing)

	_, err := NewUserService(repo).Update(context.Background(), existing.ID, input.UpdateUser{Name: ptr("")})
	var verr *domain.ValidationError
	if !errors.As(err, &verr) || verr.Field != "name" {
		t.Fatalf("want ValidationError on name, got %v", err)
	}
	if repo.updates != 0 || repo.rows[existing.ID].Name != "Alice" {
		t.Fatal("an invalid update reached the repository")
	}
}

func TestUpdateMissingUserIsNotFound(t *testing.T) {
	_, err := NewUserService(newFakeUserRepo()).Update(context.Background(), uuid.New(), input.UpdateUser{Name: ptr("x")})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
