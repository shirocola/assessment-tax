package tax

import (
	"testing"
)

func TestCalculateTaxBasic(t *testing.T) {
	t.Run("Basic Tax calculation", func(t *testing.T) {
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
	})
}

func TestCalculateTaxWithWHT(t *testing.T) {
	t.Run("Tax calculation with Withholding Tax", func(t *testing.T) {
		input := TaxCalculationInput{
			TotalIncome: 500000,
			WHT:         25000, // Withholding Tax
			Allowances:  []Allowance{},
		}
		// The expected tax calculation explanation:
		// Personal allowance is 60,000, reducing taxable income to 440,000.
		// Tax calculation:
		// - 0 to 150,000 taxed at 0% -> Tax = 0
		// - 150,001 to 440,000 taxed at 10% -> Tax = (440,000 - 150,000) * 10% = 29,000
		// - Subtract WHT of 25,000 from the calculated tax -> Remaining tax = 29,000 - 25,000 = 4,000
		expectedTax := 4000.0 // Expected tax after WHT is accounted
		result, _ := CalculateTax(input)
		if result.TotalTax != expectedTax {
			t.Errorf("Expected tax after WHT to be %.1f, got %.1f", expectedTax, result.TotalTax)
		}
	})
}

func TestCalculateTaxWithDetailedLevelsAndDonations(t *testing.T) {
	t.Run("Detailed Levels with Donations", func(t *testing.T) {
		input := TaxCalculationInput{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 200000.0},
			},
		}
		// Expected results based on the given inputs
		expectedTotalTax := 19000.0 // Expected total tax after considering allowances and tax rates.
		expectedDetails := []TaxDetail{
			{Level: "0-150,000", Tax: 0},
			{Level: "150,001-500,000", Tax: 19000},
			{Level: "500,001-1,000,000", Tax: 0},
			{Level: "1,000,001-2,000,000", Tax: 0},
			{Level: "2,000,001 and up", Tax: 0},
		}

		result, _ := CalculateTax(input)
		// Check the total tax calculated
		if result.TotalTax != expectedTotalTax {
			t.Errorf("Expected total tax to be %.1f, got %.1f", expectedTotalTax, result.TotalTax)
		}
		// Check each tax level detail
		for i, detail := range result.Details {
			if detail.Tax != expectedDetails[i].Tax {
				t.Errorf("Expected tax for level %s to be %.1f, got %.1f", detail.Level, expectedDetails[i].Tax, detail.Tax)
			}
		}
	})
}

func TestCalculateTaxWithDonations(t *testing.T) {
	input := TaxCalculationInput{
		TotalIncome: 500000,
		WHT:         0,
		Allowances: []Allowance{
			{AllowanceType: "donation", Amount: 200000},
		},
	}
	expectedTax := 19000.0 // Adjusted tax considering donation deductions
	result, _ := CalculateTax(input)
	if result.TotalTax != expectedTax {
		t.Errorf("Expected tax with donations to be %.1f, got %.1f", expectedTax, result.TotalTax)
	}
}

func TestSetPersonalDeduction(t *testing.T) {
	initialDeduction := GetPersonalDeduction()
	newDeduction := 70000.0

	// Set new personal deduction
	updatedDeduction, err := SetPersonalDeduction(newDeduction)
	if err != nil {
		t.Errorf("Error setting personal deduction: %v", err)
	}

	// Check if the deduction was updated correctly
	if updatedDeduction != newDeduction {
		t.Errorf("Expected personal deduction to be %.2f, got %.2f", newDeduction, updatedDeduction)
	}

	// Reset to initial value for other tests
	SetPersonalDeduction(initialDeduction)
}
