package todo

import (
	"context"
	"testing"
	"time"
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

	list, err := repo.List(ctx, 10, 0)
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

func TestInMemoryRepository_ListPagination(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()

	a, _ := repo.Create(ctx, "a")
	time.Sleep(5 * time.Millisecond)
	b, _ := repo.Create(ctx, "b")
	time.Sleep(5 * time.Millisecond)
	c, _ := repo.Create(ctx, "c")

	// Default order is newest first by CreatedAt -> c, b, a
	list, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3, got %d", len(list))
	}
	if list[0].ID != c.ID || list[1].ID != b.ID || list[2].ID != a.ID {
		t.Fatalf("unexpected order: %+v", list)
	}

	// Limit 2
	list, _ = repo.List(ctx, 2, 0)
	if len(list) != 2 || list[0].ID != c.ID || list[1].ID != b.ID {
		t.Fatalf("limit not applied: %+v", list)
	}

	// Offset 1
	list, _ = repo.List(ctx, 2, 1)
	if len(list) != 2 || list[0].ID != b.ID || list[1].ID != a.ID {
		t.Fatalf("offset not applied: %+v", list)
	}

	// Offset beyond size
	list, _ = repo.List(ctx, 2, 5)
	if len(list) != 0 {
		t.Fatalf("expected empty with large offset, got %d", len(list))
	}
}
