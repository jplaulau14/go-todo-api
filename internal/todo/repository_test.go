package todo

import (
	"context"
	"testing"
)

func TestInMemoryRepository_CRUD(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, "test")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if created.ID == "" || created.Title != "test" || created.Completed {
		t.Fatalf("unexpected created: %+v", created)
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("Get mismatch: got %s want %s", got.ID, created.ID)
	}

	list, err := repo.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("List error or wrong len: %v len=%d", err, len(list))
	}

	newTitle := "updated"
	done := true
	updated, err := repo.Update(ctx, created.ID, UpdateTodoRequest{Title: &newTitle, Completed: &done})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if updated.Title != "updated" || !updated.Completed {
		t.Fatalf("unexpected updated: %+v", updated)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if _, err := repo.Get(ctx, created.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}
