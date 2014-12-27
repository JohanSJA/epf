package epf

import (
	"github.com/bearbin/go-age"
	"math"
	"time"
)

// Different methods in calculating a rate
type calculationMethod int

const (
	notApplicable calculationMethod = iota
	percentage
	exactAmount
)

type rateCalculation struct {
	Minimum        float64
	Maximum        float64
	Interval       float64
	EmployerMethod calculationMethod
	EmployerAmount float64
	EmployeeMethod calculationMethod
	EmployeeAmount float64
}

type Section struct {
	Name         string            // The name of the section
	calculations []rateCalculation // Different ways to calculate the rate within the section
}

// Each section represent a different section documented in the Third Schedule.
var Sections []Section

type Rate struct {
	Section              Section
	WagesFrom            float64
	WagesTo              float64
	ContributionEmployer float64
	ContributionEmployee float64
}

var Rates []Rate

type Citizenship int

const (
	Malaysian Citizenship = iota
	PermanentResident
	NonMalaysian
)

type Employee struct {
	Citizenship                   Citizenship
	ContributionBefore1August1998 bool
	DateOfBirth                   time.Time
	Wages                         float64
}

/*
The age when someone is treated as a senior citizenship according to the
third schedule.
*/
const seniorAge int = 60

func init() {
	calculationsForA := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, percentage, 0.13, percentage, 0.11},
		{20.0, 5000.0, 20.0, percentage, 0.13, percentage, 0.11},
		{5000.0, 20000.0, 100.0, percentage, 0.12, percentage, 0.11},
	}
	a := Section{"A", calculationsForA}
	calculationsForB := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, exactAmount, 5.0, percentage, 0.11},
		{20.0, 5000.0, 20.0, exactAmount, 5.0, percentage, 0.11},
		{5000.0, 20000.0, 100.0, exactAmount, 5.0, percentage, 0.11},
	}
	b := Section{"B", calculationsForB}
	calculationsForC := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, percentage, 0.065, percentage, 0.055},
		{20.0, 5000.0, 20.0, percentage, 0.065, percentage, 0.055},
		{5000.0, 20000.0, 100.0, percentage, 0.06, percentage, 0.055},
	}
	c := Section{"C", calculationsForC}
	calculationsForD := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, exactAmount, 5.0, percentage, 0.055},
		{20.0, 5000.0, 20.0, exactAmount, 5.0, percentage, 0.055},
		{5000.0, 20000.0, 100.0, exactAmount, 5.0, percentage, 0.055},
	}
	d := Section{"D", calculationsForD}
	Sections = []Section{a, b, c, d}

	Rates = []Rate{}
	// Calculating all the rates for all the sections
	for _, sect := range Sections {
		for _, calc := range sect.calculations {
			for from := calc.Minimum; from < calc.Maximum; from += calc.Interval {
				to := from + calc.Interval
				// Different ways to calculate the amount that need to be paid
				// by employer and employee.
				var employer float64
				switch calc.EmployerMethod {
				case notApplicable:
					employer = calc.EmployerAmount
				case percentage:
					employer = math.Ceil(to * calc.EmployerAmount)
				case exactAmount:
					employer = calc.EmployerAmount
				}
				var employee float64
				switch calc.EmployeeMethod {
				case notApplicable:
					employee = calc.EmployeeAmount
				case percentage:
					employee = math.Ceil(to * calc.EmployeeAmount)
				case exactAmount:
					employee = calc.EmployeeAmount
				}
				Rates = append(Rates, Rate{sect, from + 0.01, to, employer, employee})
			}
		}
	}
}

// Return a particular section applicable to the given employee
func (e *Employee) Section() *Section {
	age := age.Age(e.DateOfBirth)
	switch {
	case e.Citizenship == Malaysian ||
		e.Citizenship == PermanentResident ||
		(e.Citizenship == NonMalaysian && e.ContributionBefore1August1998):
		switch {
		case age > seniorAge:
			return &Sections[2]
		default:
			return &Sections[0]
		}
	default:
		switch {
		case age > seniorAge:
			return &Sections[3]
		default:
			return &Sections[1]
		}
	}
	return &Section{}
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

// Create a new Malaysian Employee
func NewEmployeeMalaysian(dateOfBirth time.Time, wages float64) Employee {
	return Employee{
		Citizenship: Malaysian,
		DateOfBirth: dateOfBirth,
		Wages:       wages,
	}
}

// Create a new Permanent Resident Employee
func NewEmployeePermanentResident(dateOfBirth time.Time, wages float64) Employee {
	return Employee{
		Citizenship: PermanentResident,
		DateOfBirth: dateOfBirth,
		Wages:       wages,
	}
}

// Create a new Non Malaysian Employee
func NewEmployeeNonMalaysian(contributeBefore1August1998 bool, dateOfBirth time.Time, wages float64) Employee {
	return Employee{NonMalaysian, contributeBefore1August1998, dateOfBirth, wages}
}
