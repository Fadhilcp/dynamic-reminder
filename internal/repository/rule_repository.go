package repository

import (
	"database/sql"
	"fmt"
	"dynamic-reminder/internal/models"
)

type RuleRepository struct {
	DB *sql.DB
}

func NewRuleRepository(db *sql.DB) *RuleRepository {
	return &RuleRepository{DB: db}
}

// LogAudit helper
func (r *RuleRepository) LogAudit(actionType, entityType string, entityID int, details string) error {
	query := `INSERT INTO audit_logs (action_type, entity_type, entity_id, details) VALUES (?, ?, ?, ?)`
	_, err := r.DB.Exec(query, actionType, entityType, entityID, details)
	return err
}

func (r *RuleRepository) CreateRule(rule *models.ReminderRule) error {
	query := `INSERT INTO reminder_rules (name, condition_type, condition_value, is_active) VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, rule.Name, rule.ConditionType, rule.ConditionValue, rule.IsActive)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rule.ID = int(id)

	details := fmt.Sprintf("Created rule '%s' (Type: %s, Value: %s)", rule.Name, rule.ConditionType, rule.ConditionValue)
	return r.LogAudit("rule_created", "rule", rule.ID, details)
}

func (r *RuleRepository) GetRules() ([]models.ReminderRule, error) {
	query := `SELECT id, name, condition_type, condition_value, is_active, created_at, updated_at FROM reminder_rules`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.ReminderRule
	for rows.Next() {
		var rule models.ReminderRule
		err := rows.Scan(&rule.ID, &rule.Name, &rule.ConditionType, &rule.ConditionValue, &rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *RuleRepository) UpdateRule(id int, rule *models.ReminderRule) error {
	query := `UPDATE reminder_rules SET name = ?, condition_type = ?, condition_value = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.DB.Exec(query, rule.Name, rule.ConditionType, rule.ConditionValue, id)
	if err != nil {
		return err
	}

	details := fmt.Sprintf("Updated rule ID %d: Name='%s', Type='%s', Value='%s'", id, rule.Name, rule.ConditionType, rule.ConditionValue)
	return r.LogAudit("rule_updated", "rule", id, details)
}

func (r *RuleRepository) ToggleRuleStatus(id int, isActive bool) error {
	query := `UPDATE reminder_rules SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.DB.Exec(query, isActive, id)
	if err != nil {
		return err
	}

	details := fmt.Sprintf("Changed active status of rule ID %d to %t", id, isActive)
	return r.LogAudit("rule_status_changed", "rule", id, details)
}

func (r *RuleRepository) DeleteRule(id int) error {
	query := `DELETE FROM reminder_rules WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	details := fmt.Sprintf("Deleted rule ID %d", id)
	return r.LogAudit("rule_deleted", "rule", id, details)
}