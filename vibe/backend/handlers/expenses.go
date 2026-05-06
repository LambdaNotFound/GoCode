package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"vibe/models"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	query := "SELECT id, amount, category, description, date, created_at FROM expenses"
	args := []any{}
	if category != "" {
		query += " WHERE category = ?"
		args = append(args, category)
	}
	query += " ORDER BY date DESC, id DESC"

	rows, err := h.db.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.CreatedAt); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		expenses = append(expenses, e)
	}

	writeJSON(w, http.StatusOK, expenses)
}

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Amount <= 0 || req.Category == "" || req.Date == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "amount, category, and date are required"})
		return
	}

	result, err := h.db.ExecContext(r.Context(),
		"INSERT INTO expenses (amount, category, description, date) VALUES (?, ?, ?, ?)",
		req.Amount, req.Category, req.Description, req.Date,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	var e models.Expense
	h.db.QueryRowContext(r.Context(),
		"SELECT id, amount, category, description, date, created_at FROM expenses WHERE id = ?", id,
	).Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.CreatedAt)

	writeJSON(w, http.StatusCreated, e)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	result, err := h.db.ExecContext(r.Context(), "DELETE FROM expenses WHERE id = ?", id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if n, _ := result.RowsAffected(); n == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "expense not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var req models.UpdateExpenseRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // reject category, date, or any other field
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if req.Amount == nil && req.Description == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "provide at least one of: amount, description"})
		return
	}
	if req.Amount != nil && *req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "amount must be greater than zero"})
		return
	}

	// Build SET clause from whichever permitted fields were provided.
	setClauses := []string{}
	args := []any{}
	if req.Amount != nil {
		setClauses = append(setClauses, "amount = ?")
		args = append(args, *req.Amount)
	}
	if req.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *req.Description)
	}

	query := "UPDATE expenses SET "
	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := h.db.ExecContext(r.Context(), query, args...)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if n, _ := result.RowsAffected(); n == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "expense not found"})
		return
	}

	// Return the updated record.
	var e models.Expense
	row := h.db.QueryRowContext(r.Context(),
		"SELECT id, amount, category, description, date, created_at FROM expenses WHERE id = ?", id,
	)
	if err := row.Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.CreatedAt); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch updated expense"})
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(),
		"SELECT category, SUM(amount) as total FROM expenses GROUP BY category ORDER BY total DESC",
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	type CategoryTotal struct {
		Category string  `json:"category"`
		Total    float64 `json:"total"`
	}

	totals := []CategoryTotal{}
	for rows.Next() {
		var ct CategoryTotal
		rows.Scan(&ct.Category, &ct.Total)
		totals = append(totals, ct)
	}

	writeJSON(w, http.StatusOK, totals)
}
