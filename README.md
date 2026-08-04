# Dynamic Reminder & Background Scheduler

A robust, Go-based rule engine and background background task scheduler. This system continuously monitors pending tasks and evaluates them against dynamic, user-defined rules to trigger automated reminders. 

It was built with a strong focus on data integrity, system idempotency, and auditability.

## 🚀 Key Features

- **Dynamic Rule Engine:** Create, update, and delete reminder rules on the fly (e.g., "Remind 2 days before due", "Remind when overdue").
- **Background Scheduler:** A standalone Goroutine that safely polls the database at configurable intervals without blocking the main API thread.
- **Idempotency & Rate Limiting:** Implements a strict `last_reminded_at` cooldown mechanism to ensure the scheduler never spams users with duplicate reminders for the same task.
- **Audit Trails:** Every triggered reminder, rule creation, update, and deletion is permanently logged for system observability.
- **Graceful Error Handling:** Strict database `RowsAffected` checks prevent false-positive success responses when interacting with non-existent records.

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Gin HTTP Framework (`github.com/gin-gonic/gin`)
- **Database:** SQLite3 (`github.com/mattn/go-sqlite3`)

## 📦 Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go           # Application entry point and router setup
├── internal/
│   ├── database/             # SQLite connection and table migrations
│   ├── models/               # Data structures (Task, Rule, AuditLog)
│   ├── repository/           # Database execution layer
│   ├── handlers/             # HTTP controllers and JSON validation
│   └── scheduler/            # Background polling and rule evaluation logic
├── reminder.db               # SQLite database file (auto-generated)
└── README.md