// In `tax.go`
package tax

import (
	"math"
)

type TaxCalculationInput struct {
	TotalIncome float64
	WHT         float64
	Allowances  []Allowance
}

type Allowance struct {
	AllowanceType string
	Amount        float64
}

type TaxCalculationOutput struct {
	TotalTax float64 `json:"tax"`
}

func CalculateTax(input TaxCalculationInput) (TaxCalculationOutput, error) {
	personalAllowance := 60000.0
	netIncome := input.TotalIncome - personalAllowance

	var totalTax float64
	taxLevels := []struct {
		Level   string
		Min     float64
		Max     float64
		TaxRate float64
	}{
		{"0-150,000", 0, 150000, 0},
		{"150,001-500,000", 150001, 500000, 0.10},
	}

	for _, level := range taxLevels {
		if netIncome > level.Min {
			taxableIncome := math.Min(netIncome, level.Max) - level.Min
			tax := taxableIncome * level.TaxRate
			tax = math.Round(tax)
			totalTax += tax
		}
	}

	// Subtract WHT from the calculated total tax
	totalTax = math.Max(0, totalTax-input.WHT)
	return TaxCalculationOutput{TotalTax: totalTax}, nil
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}
