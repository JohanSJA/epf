package epf

import (
	"math"
)

// Different methods in calculating a rate
type CalculationMethod int

const (
	NotApplicable CalculationMethod = iota
	Percentage
	ExactAmount
)

type RateCalculation struct {
	Minimum        float64
	Maximum        float64
	Interval       float64
	EmployerMethod CalculationMethod
	EmployerAmount float64
	EmployeeMethod CalculationMethod
	EmployeeAmount float64
}

type Section struct {
	Name         string            // The name of the section
	Calculations []RateCalculation // Different ways to calculate the rate within the section
}

// Each section represent a different section documented in the Third Schedule.
var EpfSections []Section

type Rate struct {
	Section              Section
	WagesFrom            float64
	WagesTo              float64
	ContributionEmployer float64
	ContributionEmployee float64
}

var Rates []Rate

func init() {
	calculationsForA := []RateCalculation{
		{0.0, 10.0, 10.0, NotApplicable, 0.0, NotApplicable, 0.0},
		{10.0, 20.0, 10.0, Percentage, 0.13, Percentage, 0.11},
		{20.0, 5000.0, 20.0, Percentage, 0.13, Percentage, 0.11},
		{5000.0, 20000.0, 100.0, Percentage, 0.12, Percentage, 0.11},
	}
	a := Section{"A", calculationsForA}
	calculationsForB := []RateCalculation{
		{0.0, 10.0, 10.0, NotApplicable, 0.0, NotApplicable, 0.0},
		{10.0, 20.0, 10.0, ExactAmount, 5.0, Percentage, 0.11},
		{20.0, 5000.0, 20.0, ExactAmount, 5.0, Percentage, 0.11},
		{5000.0, 20000.0, 100.0, ExactAmount, 5.0, Percentage, 0.11},
	}
	b := Section{"B", calculationsForB}
	calculationsForC := []RateCalculation{
		{0.0, 10.0, 10.0, NotApplicable, 0.0, NotApplicable, 0.0},
		{10.0, 20.0, 10.0, Percentage, 0.065, Percentage, 0.055},
		{20.0, 5000.0, 20.0, Percentage, 0.065, Percentage, 0.055},
		{5000.0, 20000.0, 100.0, Percentage, 0.06, Percentage, 0.055},
	}
	c := Section{"C", calculationsForC}
	calculationsForD := []RateCalculation{
		{0.0, 10.0, 10.0, NotApplicable, 0.0, NotApplicable, 0.0},
		{10.0, 20.0, 10.0, ExactAmount, 5.0, Percentage, 0.055},
		{20.0, 5000.0, 20.0, ExactAmount, 5.0, Percentage, 0.055},
		{5000.0, 20000.0, 100.0, ExactAmount, 5.0, Percentage, 0.055},
	}
	d := Section{"D", calculationsForD}
	EpfSections = []Section{a, b, c, d}

	Rates = []Rate{}
	// Calculating all the rates for all the sections
	for _, sect := range EpfSections {
		for _, calc := range sect.Calculations {
			for from := calc.Minimum; from < calc.Maximum; from += calc.Interval {
				to := from + calc.Interval
				// Different ways to calculate the amount that need to be paid
				// by employer and employee.
				var employer float64
				switch calc.EmployerMethod {
				case NotApplicable:
					employer = calc.EmployerAmount
				case Percentage:
					employer = math.Ceil(to * calc.EmployerAmount)
				case ExactAmount:
					employer = calc.EmployerAmount
				}
				var employee float64
				switch calc.EmployeeMethod {
				case NotApplicable:
					employee = calc.EmployeeAmount
				case Percentage:
					employee = math.Ceil(to * calc.EmployeeAmount)
				case ExactAmount:
					employee = calc.EmployeeAmount
				}
				Rates = append(Rates, Rate{sect, from + 0.01, to, employer, employee})
			}
		}
	}
}

func rate10(section Section) Rate {
	return Rate{Section: section, WagesFrom: 0.01, WagesTo: 10}
}

// Return a particular rate for a given wages within the given section
func SectionRate(sectionName string, wages float64) Rate {
	rates := SectionRates(sectionName)
	for _, rate := range rates {
		if wages > rate.WagesFrom && wages <= rate.WagesTo {
			return rate
		}
	}
	return Rate{}
}

// Return all the rates within the given section
func SectionRates(sectionName string) []Rate {
	sRates := []Rate{}
	for _, value := range Rates {
		if value.Section.Name == sectionName {
			sRates = append(sRates, value)
		}
	}
	return sRates
}

// Return the total contribution from employer and employee
func (r *Rate) ContributionTotal() float64 {
	return r.ContributionEmployer + r.ContributionEmployee
}
