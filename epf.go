package epf

import (
	"errors"
	"math"
	"time"

	"github.com/bearbin/go-age"
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

// Section representing a section inside the Third Schedule.
type Section struct {
	Name         string            // The name of the section
	Description  string            // Explanation of the section
	calculations []rateCalculation // Different ways to calculate the rate within the section
}

// Each section represent a different section documented in the Third Schedule.
var Sections []Section

// Rate represents a single row in a section.
type Rate struct {
	WagesFrom            float64
	WagesTo              float64
	ContributionEmployer float64
	ContributionEmployee float64
}

// Citizenship represents the citizenship status of an employee.
type Citizenship int

// Different types of citizenship status.
const (
	Unknown Citizenship = iota
	Malaysian
	PermanentResident
	NonMalaysian
)

// Employee contains employee information.
// ContributionBefore1August1998 is actually only required for Non Malaysian.
// The rest however are required for all employees.
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
		{0, 10, 10, notApplicable, 0, notApplicable, 0},
		{10, 20, 10, percentage, 0.15, percentage, 0.1},
		{20, 40, 20, percentage, 0.15, percentage, 0.1},
		{40, 5000, 20, percentage, 0.13, percentage, 0.08},
		{5000, 20000, 100, percentage, 0.12, percentage, 0.08},
	}
	descA := `1. The rate of monthly contributions specified in this Part shall apply to the following employees until the employees attain the age of sixty years:
a) employees who are Malaysian citizens;
b) employees who are not Malaysian citizens but are permanent residents of Malaysia; and
c) employees who are not Malaysian citizens who have elected to contribute before 1 August 1998.
2. The Board may prescribe any limit on the amount of wages and contributions for the purpose of determining any amount of contribution to be paid on the wages for the month by each employer for each employee.
3. For the purposes of subsection 43(3) and section 44A, the amount of contribution for the month is limited to any limit on the total contribution as prescribed by the Board.`
	a := Section{"A", descA, calculationsForA}
	calculationsForB := []rateCalculation{
		{0, 10, 10, notApplicable, 0, notApplicable, 0},
		{10, 20, 10, exactAmount, 5, percentage, 0.1},
		{20, 40, 20, exactAmount, 5, percentage, 0.1},
		{40, 5000, 20, exactAmount, 5, percentage, 0.08},
		{5000, 20000, 100, exactAmount, 5, percentage, 0.08},
	}
	descB := `1. The rate of monthly contributions specified in this Part shall apply to the following employees until the employees attain the age of sixty years:
a) employees who are not Malaysian citizens and have elected to contribute on or after 1 August 1998;
b) employees who are not Malaysian citizens and have elected to contribute under paragraph 3 of the First Schedule on or after 1 August 1998; and
c) employees who are not Malaysian citizens and have elected to contribute under paragraph 6 of the First Schedule on or after 1 August 2001.
2. The Board may prescribe any limit on the amount of wages and contributions for the purpose of determining any amount of contribution to be paid on the wages for the month by each employer for each employee.
3. For the purposes of subsection 43(3), the amount of contribution for the month is limited to any limit on the total contribution as prescribed by the Board.`
	b := Section{"B", descB, calculationsForB}
	calculationsForC := []rateCalculation{
		{0, 10, 10, notApplicable, 0, notApplicable, 0},
		{10, 20, 10, exactAmount, 2, exactAmount, 1},
		{20, 40, 20, exactAmount, 2, exactAmount, 1},
		{40, 5000, 20, percentage, 0.065, percentage, 0.04},
		{5000, 20000, 100, percentage, 0.06, percentage, 0.04},
	}
	descC := `1. The rate of monthly contributions specified in this Part shall apply to the following employees who have attained the age of sixty years:
a) employees who are Malaysian citizens;
b) employees who are not Malaysian citizens but are permanent residents in Malaysia; and
c) employees who are not Malaysian citizens who have elected to contribute before 1 August 1998.
2. The Board may prescribe any limit on the amount of wages and contributions for the purpose of determining any amount of contribution to be paid on the wages for the month by each employer for each employee.
3. For the purposes of subsection 43(3), the amount of contribution for the month is limited to any limit on the total contribution as prescribed by the Board.`
	c := Section{"C", descC, calculationsForC}
	calculationsForD := []rateCalculation{
		{0, 10, 10, notApplicable, 0, notApplicable, 0},
		{10, 20, 10, exactAmount, 5, exactAmount, 1},
		{20, 40, 20, exactAmount, 5, exactAmount, 2},
		{40, 5000, 20, exactAmount, 5, percentage, 0.04},
		{5000, 20000, 100, exactAmount, 5, percentage, 0.04},
	}
	descD := `1. The rate of monthly contributions specified in this Part shall apply to the following employees who have attained the age of sixty years:
a) employees who are not Malaysian citizens and have elected to contribute on or after 1 August 1998;
b) employees who are not Malaysian citizens and have elected to contribute under paragraph 3 of the First Schedule on or after 1 August 1998; and
c) employees who are not Malaysian citizens and have elected to contribute under paragraph 6 of the First Schedule on or after 1 August 2001.
2. The Board may prescribe any limit on the amount of wages and contributions for the purpose of determining the amount of contribution to be paid on the wages for the month by each employer for each employee.
3. For the purposes of subsection 43(3), the amount of contribution for the month is limited to any limit on the total contribution as prescribed by the Board.`
	d := Section{"D", descD, calculationsForD}
	Sections = []Section{a, b, c, d}
}

// Calculate what is the rate based on the method given.
func calculate(method calculationMethod, base float64, amount float64) float64 {
	switch method {
	case notApplicable:
		return amount
	case percentage:
		return math.Ceil(base * amount)
	default:
		return amount
	}
}

// SectionByName looks for a particular section based on its name.
func SectionByName(name string) (Section, error) {
	for _, sec := range Sections {
		if sec.Name == name {
			return sec, nil
		}
	}
	return Section{}, errors.New("Invalid section.")
}

// Rates returns all the rates within a particular section. These are the rates
// that are listed within the table inside Third Schedule.
func (s *Section) Rates() []Rate {
	rates := []Rate{}
	for _, calc := range s.calculations {
		for from := calc.Minimum; from < calc.Maximum; from += calc.Interval {
			to := from + calc.Interval
			employer := calculate(calc.EmployerMethod, to, calc.EmployerAmount)
			employee := calculate(calc.EmployeeMethod, to, calc.EmployeeAmount)
			rates = append(rates, Rate{from + 0.01, to, employer, employee})
		}
	}
	return rates
}

// Rate returns a particular rate for a given wages within the given section
func (s *Section) Rate(wages float64) Rate {
	// Check whether the rate is within the normal table. Return if it is.
	rates := s.Rates()
	for _, rate := range rates {
		if wages > rate.WagesFrom && wages <= rate.WagesTo {
			return rate
		}
	}
	// Calculate the rate since the rate is not is the normal table if it
	// reaches this stage.
	calc := s.calculations[len(s.calculations)-1]
	from, to := wages, wages
	employer := calculate(calc.EmployerMethod, to, calc.EmployerAmount)
	employee := calculate(calc.EmployeeMethod, to, calc.EmployeeAmount)
	return Rate{from, to, employer, employee}
}

// Section returns a particular section applicable to the given employee.
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
}

// Sections returns a list of applicable sections
func (e *Employee) Sections() []*Section {
	switch {
	case e.Citizenship != Unknown && !e.DateOfBirth.IsZero():
		// If both citizenship and date of birth is known, we can actually get
		// a precise section. So return a list of 1 element.
		return []*Section{e.Section()}
	case e.Citizenship != Unknown:
		// Get applicable sections if citizenship is not unknown.
		switch {
		case e.Citizenship == Malaysian ||
			e.Citizenship == PermanentResident ||
			(e.Citizenship == NonMalaysian && e.ContributionBefore1August1998):
			return []*Section{&Sections[0], &Sections[2]}
		default:
			return []*Section{&Sections[1], &Sections[3]}
		}
	case !e.DateOfBirth.IsZero():
		// Get applicable sections if date of birth is not unknown.
		age := age.Age(e.DateOfBirth)
		switch {
		case age > seniorAge:
			return []*Section{&Sections[2], &Sections[3]}
		default:
			return []*Section{&Sections[0], &Sections[1]}
		}
	default:
		return []*Section{&Sections[0], &Sections[1], &Sections[2], &Sections[3]}
	}
}

// Rate returns a particular rate applicable to the given employee.
func (e *Employee) Rate() Rate {
	sec := e.Section()
	return sec.Rate(e.Wages)
}

// ContributionTotal returns the total contribution from employer and employee
func (r *Rate) ContributionTotal() float64 {
	return r.ContributionEmployer + r.ContributionEmployee
}

// NewEmployeeMalaysian creates a new Malaysian Employee
func NewEmployeeMalaysian(dateOfBirth time.Time, wages float64) Employee {
	return Employee{
		Citizenship: Malaysian,
		DateOfBirth: dateOfBirth,
		Wages:       wages,
	}
}

// NewEmployeePermanentResident creates a new Permanent Resident Employee
func NewEmployeePermanentResident(dateOfBirth time.Time, wages float64) Employee {
	return Employee{
		Citizenship: PermanentResident,
		DateOfBirth: dateOfBirth,
		Wages:       wages,
	}
}

// NewEmployeeNonMalaysian creates a new Non Malaysian Employee
func NewEmployeeNonMalaysian(contributeBefore1August1998 bool, dateOfBirth time.Time, wages float64) Employee {
	return Employee{
		NonMalaysian,
		contributeBefore1August1998,
		dateOfBirth,
		wages,
	}
}
