package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shirocola/assessment-tax/internal/tax" // Assuming this is the correct import path
)

func TaxCalculationHandler(w http.ResponseWriter, r *http.Request) {
	var input tax.TaxCalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := tax.CalculateTax(input)
	if err != nil {
		http.Error(w, "Error calculating tax", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func SetPersonalDeductionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define a struct to decode the request body
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the personal deduction (assuming the method SetPersonalDeduction exists)
	updatedDeduction, err := tax.SetPersonalDeduction(req.Amount)
	if err != nil {
		http.Error(w, "Failed to set personal deduction", http.StatusInternalServerError)
		return
	}

	// Respond with the updated personal deduction
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		PersonalDeduction float64 `json:"personalDeduction"`
	}{
		PersonalDeduction: updatedDeduction,
	})
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tax/calculation", TaxCalculationHandler).Methods("POST")
	r.HandleFunc("/admin/deductions/personal", SetPersonalDeductionHandler).Methods("POST")
}
