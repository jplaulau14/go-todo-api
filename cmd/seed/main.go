package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}
	count := 25
	if v := os.Getenv("SEED_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			count = n
		}
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	titles := []string{
		"Buy groceries", "Read a book", "Write blog post", "Exercise", "Call a friend",
		"Plan vacation", "Clean kitchen", "Fix bug #42", "Review PR", "Learn Go generics",
	}
	tags := []string{"home", "work", "study", "health", "fun"}

	now := time.Now().UTC()

	stmt, err := db.Prepare(`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	defer func() { _ = stmt.Close() }()

	inserted := 0
	for i := 0; i < count; i++ {
		id := uuid.NewString()
		title := fmt.Sprintf("%s [%s] #%d", titles[rand.Intn(len(titles))], tags[rand.Intn(len(tags))], i+1)
		completed := i%3 == 0
		createdAt := now.Add(-time.Duration(rand.Intn(96)) * time.Hour) // within last 4 days
		updatedAt := createdAt
		if completed {
			updatedAt = createdAt.Add(time.Duration(rand.Intn(12)) * time.Hour)
		}
		if _, err := stmt.Exec(id, title, completed, createdAt, updatedAt); err != nil {
			log.Fatalf("insert: %v", err)
		}
		inserted++
	}

	log.Printf("seeded %d todos", inserted)
}
