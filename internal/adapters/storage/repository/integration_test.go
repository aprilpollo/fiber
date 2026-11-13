package repository

import (
	"context"
	"errors"
	"os"
	"testing"

	"aprilpollo/internal/adapters/storage/orm/models"
	"aprilpollo/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// openTestDB connects to TEST_POSTGRES_URI and gives each test fresh tables.
// Without the variable the test is skipped, so `go test ./...` needs no database.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	uri := os.Getenv("TEST_POSTGRES_URI")
	if uri == "" {
		t.Skip("TEST_POSTGRES_URI not set")
	}
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrator().DropTable(&models.UserModel{}, &models.BookModel{}); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.UserModel{}, &models.BookModel{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestUserRepositoryAgainstPostgres(t *testing.T) {
	db := openTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	alice, _ := domain.NewUser("Alice", "alice@example.com")
	if err := repo.Save(ctx, alice); err != nil {
		t.Fatal(err)
	}
	if alice.CreatedAt.IsZero() {
		t.Fatal("Save must sync timestamps back to the domain")
	}
	bob, _ := domain.NewUser("Bob", "bob@example.com")
	if err := repo.Save(ctx, bob); err != nil {
		t.Fatal(err)
	}

	t.Run("duplicate email is ErrConflict", func(t *testing.T) {
		dup, _ := domain.NewUser("Other", "alice@example.com")
		if err := repo.Save(ctx, dup); !errors.Is(err, domain.ErrConflict) {
			t.Fatalf("want ErrConflict, got %v", err)
		}
	})

	t.Run("missing id is ErrNotFound", func(t *testing.T) {
		if _, err := repo.FindByID(ctx, uuid.New()); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("FindByID: want ErrNotFound, got %v", err)
		}
		ghost, _ := domain.NewUser("Ghost", "ghost@example.com")
		if err := repo.Update(ctx, ghost); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("Update: want ErrNotFound, got %v", err)
		}
		if _, err := repo.FindByID(ctx, ghost.ID); !errors.Is(err, domain.ErrNotFound) {
			t.Fatal("Update of a missing user must not insert it")
		}
		if err := repo.Delete(ctx, uuid.New()); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("Delete: want ErrNotFound, got %v", err)
		}
	})

	t.Run("total counts the filtered rows, not the table", func(t *testing.T) {
		users, total, err := repo.FindAll(ctx, domain.ListQuery{
			Filters: []domain.Filter{{Field: "name", Op: domain.OpContains, Value: "li"}},
			Limit:   10,
		})
		if err != nil {
			t.Fatal(err)
		}
		if total != 1 || len(users) != 1 || users[0].Email != "alice@example.com" {
			t.Fatalf("got total=%d users=%+v", total, users)
		}
	})

	t.Run("sort and page", func(t *testing.T) {
		users, total, err := repo.FindAll(ctx, domain.ListQuery{Sort: "name", Desc: true, Limit: 1, Offset: 1})
		if err != nil {
			t.Fatal(err)
		}
		if total != 2 || len(users) != 1 || users[0].Name != "Alice" {
			t.Fatalf("got total=%d users=%+v", total, users)
		}
	})

	t.Run("update keeps created_at and moves updated_at", func(t *testing.T) {
		got, _ := repo.FindByID(ctx, bob.ID)
		created, updated := got.CreatedAt, got.UpdatedAt
		got.Name = "Robert"
		if err := repo.Update(ctx, got); err != nil {
			t.Fatal(err)
		}
		again, _ := repo.FindByID(ctx, bob.ID)
		if again.Name != "Robert" || !again.CreatedAt.Equal(created) || !again.UpdatedAt.After(updated) {
			t.Fatalf("got %+v (created %v, updated %v)", again, created, updated)
		}
		if !got.UpdatedAt.Equal(again.UpdatedAt) {
			t.Fatalf("Update must sync updated_at back: domain %v, stored %v", got.UpdatedAt, again.UpdatedAt)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := repo.Delete(ctx, bob.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.FindByID(ctx, bob.ID); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("want ErrNotFound after delete, got %v", err)
		}
	})
}
