package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"dynamic-reminder/internal/models"
)

func Run(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for t := range ticker.C {
			log.Printf("[SCHEDULER] Running evaluation cycle at %v", t.Format(time.RFC3339))
			evaluateRules(db)
		}
	}()
}

func evaluateRules(db *sql.DB) {
	
	rules := getActiveRules(db)
	if len(rules) == 0 {
		return
	}

	tasks := getPendingTasks(db)
	if len(tasks) == 0 {
		return
	}

	now := time.Now()
	for _, task := range tasks {

		
		if task.LastRemindedAt.Valid {
			if time.Since(task.LastRemindedAt.Time) < 30*time.Second {
				continue 
			}
		}

		for _, rule := range rules {
			
			triggered := false
			
			switch rule.ConditionType {
			case "days_before_due":
				days, err := strconv.Atoi(rule.ConditionValue)
				if err == nil {
					targetDate := now.AddDate(0, 0, days)
					if task.DueDate.Year() == targetDate.Year() && task.DueDate.YearDay() == targetDate.YearDay() {
						triggered = true
					}
				}
			case "status_overdue":
				if task.DueDate.Before(now) {
					triggered = true
				}
			}

			if triggered {
				message := fmt.Sprintf("REMINDER: Task '%s' triggered by rule '%s'", task.Title, rule.Name)
				
				fmt.Println("--------------------------------------------------")
				fmt.Println(message)
				fmt.Println("--------------------------------------------------")
				
				logAudit(db, "reminder_triggered", "task", task.ID, message)

				// to update last reminded at in task 
				_, err := db.Exec(`UPDATE tasks SET last_reminded_at = CURRENT_TIMESTAMP WHERE id = ?`, task.ID)
				if err != nil {
					log.Printf("Failed to update last_reminded_at for task %d: %v", task.ID, err)
				}
				// to avoid remind the same task
				break
			}
		}
	}
}

func getActiveRules(db *sql.DB) []models.ReminderRule {
	rows, err := db.Query(`SELECT id, name, condition_type, condition_value FROM reminder_rules WHERE is_active = 1`)
	if err != nil {
		log.Printf("Scheduler error fetching rules: %v", err)
		return nil
	}
	defer rows.Close()

	var rules []models.ReminderRule
	for rows.Next() {
		var r models.ReminderRule
		if err := rows.Scan(&r.ID, &r.Name, &r.ConditionType, &r.ConditionValue); err == nil {
			rules = append(rules, r)
		}
	}
	return rules
}

func getPendingTasks(db *sql.DB) []models.Task {
	rows, err := db.Query(`SELECT id, title, status, due_date, last_reminded_at FROM tasks WHERE status = 'pending'`)
	if err != nil {
		log.Printf("Scheduler error fetching tasks: %v", err)
		return nil
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.DueDate, &t.LastRemindedAt); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}


func logAudit(db *sql.DB, actionType, entityType string, entityID int, details string) {
	query := `INSERT INTO audit_logs (action_type, entity_type, entity_id, details) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, actionType, entityType, entityID, details)
	if err != nil {
		log.Printf("Scheduler failed to log audit: %v", err)
	}
}