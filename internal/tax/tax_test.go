package tax

import (
	"testing"
)

func TestCalculateTaxBasic(t *testing.T) {
	input := TaxCalculationInput{
		TotalIncome: 500000,        // Total income before deductions
		WHT:         0,             // No withholding tax for this basic scenario
		Allowances:  []Allowance{}, // No additional allowances specified
	}
	// The expected tax calculation explanation:
	// Personal allowance is 60,000, so taxable income is 500,000 - 60,000 = 440,000.
	// Tax rates are applied as follows:
	// - First 150,000 is taxed at 0% -> Tax = 0
	// - Next 290,000 (from 150,001 to 440,000) is taxed at 10% -> Tax = 290,000 * 10% = 29,000
	expectedTax := 29000.0 // Expected tax considering the personal allowance and tax rates.
	result, _ := CalculateTax(input)
	if result.TotalTax != expectedTax {
		t.Errorf("Expected tax to be %.1f, got %.1f. Total income: %.1f, Allowance applied: 60000", expectedTax, result.TotalTax, input.TotalIncome)
	}
}
