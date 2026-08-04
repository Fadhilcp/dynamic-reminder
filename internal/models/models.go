package models

import (
	"database/sql"
	"time"
)
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // 'pending', 'completed'
	DueDate     time.Time `json:"due_date"`
	LastRemindedAt sql.NullTime `json:"last_reminded_at"`	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReminderRule struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	ConditionType  string    `json:"condition_type"`  // 'days_before_due', 'status_overdue'
	ConditionValue string    `json:"condition_value"` // '3', 'pending'
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         int       `json:"id"`
	ActionType string    `json:"action_type"` // 'rule_created', 'reminder_triggered'
	EntityType string    `json:"entity_type"` // 'rule', 'task'
	EntityID   int       `json:"entity_id"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}