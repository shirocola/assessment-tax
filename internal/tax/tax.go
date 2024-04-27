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

type TaxDetail struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxCalculationOutput struct {
	TotalTax float64     `json:"tax"`
	Details  []TaxDetail `json:"taxLevel"`
}

func CalculateTax(input TaxCalculationInput) (TaxCalculationOutput, error) {
	personalAllowance := 60000.0
	netIncome := input.TotalIncome - personalAllowance

	for _, a := range input.Allowances {
		if a.AllowanceType == "personal" && a.Amount > 0 {
			netIncome += personalAllowance
			netIncome -= a.Amount
		} else if a.AllowanceType == "donation" {
			donationDeduction := a.Amount
			if donationDeduction > 100000 {
				donationDeduction = 100000
			}
			netIncome -= donationDeduction
		} else {
			netIncome -= a.Amount
		}
	}

	var totalTax float64
	details := []TaxDetail{}

	taxLevels := []struct {
		Level   string
		Min     float64
		Max     float64
		TaxRate float64
	}{
		{"0-150,000", 0, 150000, 0},
		{"150,001-500,000", 150001, 500000, 0.10},
		{"500,001-1,000,000", 500001, 1000000, 0.15},
		{"1,000,001-2,000,000", 1000001, 2000000, 0.20},
		{"2,000,001 and up", 2000001, math.MaxFloat64, 0.35},
	}

	for _, level := range taxLevels {
		if netIncome > level.Min {
			taxableIncome := min(netIncome, level.Max) - level.Min
			tax := taxableIncome * level.TaxRate
			tax = math.Round(tax)
			totalTax += tax
			details = append(details, TaxDetail{Level: level.Level, Tax: tax})
		} else {
			details = append(details, TaxDetail{Level: level.Level, Tax: 0})
		}
	}

	totalTax -= input.WHT
	totalTax = math.Round(totalTax)
	if totalTax < 0 {
		totalTax = 0
	}

	return TaxCalculationOutput{TotalTax: totalTax, Details: details}, nil
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
