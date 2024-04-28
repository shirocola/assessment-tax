package http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestExample(t *testing.T) {
	t.Log("Example test runs")
}

func TestUploadTaxCalculationCSVHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/tax/calculations/upload-csv", UploadTaxCalculationCSVHandler).Methods("POST")

	// Simulate CSV file content
	csvContent := `totalIncome,wht,donation
500000,0,0
600000,40000,20000
750000,50000,15000
`
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(part, strings.NewReader(csvContent))
	writer.Close()

	req, err := http.NewRequest("POST", "/tax/calculations/upload-csv", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var actualResponse struct {
		Taxes []taxResponse `json:"taxes"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	expectedTaxes := []taxResponse{
		{TotalIncome: 500000.0, Tax: 29000.0},
		{TotalIncome: 600000.0, Tax: 0.0},
		{TotalIncome: 750000.0, Tax: 11250.0},
	}

	if len(actualResponse.Taxes) != len(expectedTaxes) {
		t.Fatalf("expected %d tax entries, got %d", len(expectedTaxes), len(actualResponse.Taxes))
	}

	for i, tax := range actualResponse.Taxes {
		if tax.TotalIncome != expectedTaxes[i].TotalIncome || tax.Tax != expectedTaxes[i].Tax {
			t.Errorf("entry %d: expected %+v, got %+v", i, expectedTaxes[i], tax)
		}
	}
}

func TestUploadTaxCalculationCSVHandler_InvalidData(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/tax/calculations/upload-csv", UploadTaxCalculationCSVHandler).Methods("POST")

	// Simulate invalid CSV file content
	invalidCsvContent := `totalIncome,wht,donation
500000,not_a_number,0
600000,40000,20000
750000,,15000
`
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "invalid_taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(part, strings.NewReader(invalidCsvContent))
	writer.Close()

	req, err := http.NewRequest("POST", "/tax/calculations/upload-csv", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Expecting a bad request or similar error status due to invalid input
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid data: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSetKReceiptDeductionHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/admin/deductions/k-receipt", SetKReceiptDeductionHandler).Methods("POST")

	body := strings.NewReader(`{"amount": 70000}`)
	req, err := http.NewRequest("POST", "/admin/deductions/k-receipt", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Using a struct to match expected JSON structure
	var resp struct {
		KReceipt float64 `json:"kReceipt"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedKReceipt := 70000.0
	if resp.KReceipt != expectedKReceipt {
		t.Errorf("Expected kReceipt to be %.2f, got %.2f", expectedKReceipt, resp.KReceipt)
	}
}
