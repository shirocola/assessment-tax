package tax

import (
	"math"
	"sync"
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

var (
	personalDeductionLock sync.Mutex
	personalDeduction     float64 = 60000.0
)

func SetPersonalDeduction(amount float64) (float64, error) {
	personalDeductionLock.Lock()
	defer personalDeductionLock.Unlock()

	personalDeduction = amount
	return personalDeduction, nil
}

func CalculateTax(input TaxCalculationInput) (TaxCalculationOutput, error) {
	personalAllowance := 60000.0
	netIncome := input.TotalIncome - personalAllowance

	for _, allowance := range input.Allowances {
		switch allowance.AllowanceType {
		case "personal":
			if allowance.Amount > 0 {
				netIncome += personalAllowance
				netIncome -= allowance.Amount
			}
		case "donation":
			if allowance.Amount > 100000 {
				allowance.Amount = 100000
			}
			netIncome -= allowance.Amount
		case "k-receipt":
			if allowance.Amount > 50000 {
				allowance.Amount = 50000
			}
			netIncome -= allowance.Amount
		}
	}

	taxLevels := []struct {
		Level   string
		Min     float64
		Max     float64
		TaxRate float64
	}{
		{"0-150,000", 0, 150000, 0.00},
		{"150,001-500,000", 150001, 500000, 0.10},
		{"500,001-1,000,000", 500001, 1000000, 0.15},
		{"1,000,001-2,000,000", 1000001, 2000000, 0.20},
		{"2,000,001 and up", 2000001, math.MaxFloat64, 0.35},
	}

	var details []TaxDetail
	totalTax := 0.0

	// Calculate tax for each level
	for _, level := range taxLevels {
		if netIncome > level.Min {
			taxableIncome := math.Min(netIncome, level.Max) - level.Min
			tax := taxableIncome * level.TaxRate
			tax = math.Round(tax)
			totalTax += tax
			details = append(details, TaxDetail{Level: level.Level, Tax: tax})
		} else {
			details = append(details, TaxDetail{Level: level.Level, Tax: 0})
		}
	}

	totalTax -= input.WHT
	// totalTax = math.Round(totalTax)
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

func GetPersonalDeduction() float64 {
	personalDeductionLock.Lock()
	defer personalDeductionLock.Unlock()
	return personalDeduction
}
