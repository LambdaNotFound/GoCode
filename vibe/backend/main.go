package main

import (
	"encoding/json"
	"log"
	"net/http"

	"vibe/db"
	"vibe/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	database, err := db.Init("./expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	h := handlers.New(database)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /expenses", h.ListExpenses)
	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("DELETE /expenses/{id}", h.DeleteExpense)
	mux.HandleFunc("GET /expenses/summary", h.GetSummary)

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
