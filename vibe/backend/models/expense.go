package models

type Expense struct {
	ID          int64   `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"created_at"`
}

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// UpdateExpenseRequest only permits Amount and Description to be changed.
// Any other fields in the request body will be rejected.
type UpdateExpenseRequest struct {
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
}
