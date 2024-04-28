package http

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shirocola/assessment-tax/internal/tax"
)

type taxResponse struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

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
	json.NewEncoder(w).Encode(result)
}

func SetPersonalDeductionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedDeduction, err := tax.SetPersonalDeduction(req.Amount)
	if err != nil {
		http.Error(w, "Failed to set personal deduction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		PersonalDeduction float64 `json:"personalDeduction"`
	}{
		PersonalDeduction: updatedDeduction,
	})
}

func UploadTaxCalculationCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("taxFile")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var results []taxResponse

	if _, err := reader.Read(); err != nil { // Skip the header
		http.Error(w, "CSV file read error", http.StatusBadRequest)
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
			return
		}

		if len(record) != 3 {
			http.Error(w, "CSV format error", http.StatusBadRequest)
			return
		}

		totalIncome, err1 := strconv.ParseFloat(record[0], 64)
		wht, err2 := strconv.ParseFloat(record[1], 64)
		donation, err3 := strconv.ParseFloat(record[2], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			http.Error(w, "Invalid data in CSV", http.StatusBadRequest)
			return
		}

		taxOutput, err := tax.CalculateTax(tax.TaxCalculationInput{
			TotalIncome: totalIncome,
			WHT:         wht,
			Allowances:  []tax.Allowance{{AllowanceType: "donation", Amount: donation}},
		})
		if err != nil {
			http.Error(w, "Error calculating tax", http.StatusInternalServerError)
			return
		}
		results = append(results, taxResponse{
			TotalIncome: totalIncome,
			Tax:         taxOutput.TotalTax,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Taxes []taxResponse `json:"taxes"`
	}{
		Taxes: results,
	})
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tax/calculations", TaxCalculationHandler).Methods("POST")
	r.HandleFunc("/admin/deductions/personal", SetPersonalDeductionHandler).Methods("POST")
	r.HandleFunc("/tax/calculations/upload-csv", UploadTaxCalculationCSVHandler).Methods("POST")
}
