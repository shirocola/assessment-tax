package http

import (
	"encoding/json"
	"net/http"

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
