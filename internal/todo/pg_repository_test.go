package todo

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresRepository_CRUD(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set; skipping integration test")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, "it")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("mismatch")
	}
	list, err := repo.List(ctx, 10, 0)
	if err != nil || len(list) == 0 {
		t.Fatalf("List: %v len=%d", err, len(list))
	}
	title := "changed"
	done := true
	updated, err := repo.Update(ctx, created.ID, UpdateTodoRequest{Title: &title, Completed: &done})
	if err != nil || updated.Title != "changed" || !updated.Completed {
		t.Fatalf("Update: %v %+v", err, updated)
	}
	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(ctx, created.ID); err == nil {
		t.Fatalf("expected not found")
	}
}
