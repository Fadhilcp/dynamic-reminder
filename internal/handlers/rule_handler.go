package handlers

import (
	"net/http"
	"strconv"

	"dynamic-reminder/internal/models"
	"dynamic-reminder/internal/repository"
	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	Repo *repository.RuleRepository
}

func NewRuleHandler(repo *repository.RuleRepository) *RuleHandler {
	return &RuleHandler{Repo: repo}
}

func (h *RuleHandler) CreateRule(c *gin.Context) {
	var rule models.ReminderRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *RuleHandler) GetRules(c *gin.Context) {
	rules, err := h.Repo.GetRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func (h *RuleHandler) UpdateRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	var rule models.ReminderRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateRule(id, &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule updated successfully"})
}

func (h *RuleHandler) ToggleRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.ToggleRuleStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle rule status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule status updated successfully"})
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	if err := h.Repo.DeleteRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}