package database

import (
	"database/sql"
	"log"
	"time"
)

func SeedTasks(db *sql.DB) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check task count: %v", err)
	}

	if count > 0 {
		log.Println("Tasks already seeded, skipping...")
		return
	}

	log.Println("Seeding 5 sample tasks...")

	now := time.Now()
	tasks := []struct {
		Title       string
		Description string
		Status      string
		DueDate     time.Time
	}{
		{"Task 1", "Complete documentation", "pending", now.AddDate(0, 0, 3)},
		{"Task 2", "Fix critical bug", "pending", now.Add(-24 * time.Hour)},
		{"Task 3", "Deploy to production", "pending", now.Add(1 * time.Hour)},
		{"Task 4", "Review pull requests", "completed", now.AddDate(0, 0, -1)},
		{"Task 5", "Weekly team sync", "pending", now.AddDate(0, 0, 1)},
	}
	
	query := `INSERT INTO tasks (title, description, status, due_date) VALUES (?, ?, ?, ?)`

	for _, t := range tasks {
		_, err := db.Exec(query, t.Title, t.Description, t.Status, t.DueDate)
		if err != nil {
			log.Fatalf("Failed to seed task %s: %v", t.Title, err)
		}
	}
	
	log.Println("Successfully seeded 5 tasks.")
}