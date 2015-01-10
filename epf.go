package epf

import (
	"errors"
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

// Type representing a section inside the Third Schedule.
type Section struct {
	Name         string            // The name of the section
	Description  string            // Explanation of the section
	calculations []rateCalculation // Different ways to calculate the rate within the section
}

// Each section represent a different section documented in the Third Schedule.
var Sections []Section

type Rate struct {
	WagesFrom            float64
	WagesTo              float64
	ContributionEmployer float64
	ContributionEmployee float64
}

type Citizenship int

const (
	Unknown Citizenship = iota
	Malaysian
	PermanentResident
	NonMalaysian
)

/*
Type representing an employee. ContributionBefore1August1998 is actually only
required for Non Malaysian. The rest however are required for all employees.
*/
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
	descA := `The rate of monthly contributions specified in this Part shall apply to:-
- employees who are Malaysian citizens;
- employees who are not Malaysian citizens but are permanent residents of Malaysia; and
- employees who are not Malaysian citizens who have elected to contribute before 1 August 1998.
until the the employee reached the age of sixty years

In this Part:-
- the amount of wages for the month which shall be contributd to the Fund by each employer for each employee shall be according to any limit on the amount of wages and contributions as prescribed by the Board; and
- the amount of contributions for the month for the purpose of subsection 43(3) and section 44A is limited to any limit on the total contributions as prescribed by the Board.`
	a := Section{"A", descA, calculationsForA}
	calculationsForB := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, exactAmount, 5.0, percentage, 0.11},
		{20.0, 5000.0, 20.0, exactAmount, 5.0, percentage, 0.11},
		{5000.0, 20000.0, 100.0, exactAmount, 5.0, percentage, 0.11},
	}
	descB := `The rate of monthly contributions specified in this Part shall apply to employees who are not Malaysian citizens:-
- who elect to contribute on or after 1 August 1998;
- who elect to contribute under subsection 54(3) on or after 1 August 1998; and
- who elect to contribute under paragraph 6 of the First Schedule on or after 1 August 2001.
until the the employee reached the age of sixty years

In this Part:-
- the amount of wages for the month which shall be contributd to the Fund by each employer for each employee shall be according to any limit on the amount of wages and contributions as prescribed by the Board; and
- the amount of contributions for the month for the purpose of subsection 43(3) and section 44A is limited to any limit on the total contributions as prescribed by the Board.`
	b := Section{"B", descB, calculationsForB}
	calculationsForC := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, percentage, 0.065, percentage, 0.055},
		{20.0, 5000.0, 20.0, percentage, 0.065, percentage, 0.055},
		{5000.0, 20000.0, 100.0, percentage, 0.06, percentage, 0.055},
	}
	descC := `The rate of monthly contributions specified in this Part shall apply to:-
- employees who are Malaysian citizens;
- employees who are not Malaysian citizens but are permanent residents of Malaysia; and
- employees who are not Malaysian citizens who have elected to contribute before 1 August 1998.
who have attained the age of sixty years.

In this Part:-
- the amount of wages for the month which shall be contributd to the Fund by each employer for each employee shall be according to any limit on the amount of wages and contributions as prescribed by the Board; and
- the amount of contributions for the month for the purpose of subsection 43(3) and section 44A is limited to any limit on the total contributions as prescribed by the Board.`
	c := Section{"C", descC, calculationsForC}
	calculationsForD := []rateCalculation{
		{0.0, 10.0, 10.0, notApplicable, 0.0, notApplicable, 0.0},
		{10.0, 20.0, 10.0, exactAmount, 5.0, percentage, 0.055},
		{20.0, 5000.0, 20.0, exactAmount, 5.0, percentage, 0.055},
		{5000.0, 20000.0, 100.0, exactAmount, 5.0, percentage, 0.055},
	}
	descD := `The rate of monthly contributions specified in this Part shall apply to employees who are not Malaysian citizens:-
- who elect to contribute on or after 1 August 1998;
- who elect to contribute under subsection 54(3) on or after 1 August 1998; and
- who elect to contribute under paragraph 3 of the First Schedule on or after 1 August 1998; and
- who elect to contribute under paragraph 6 of the First Schedule on or after 1 August 2001.
who have attained the arge of sixty years.

In this Part:-
- the amount of wages for the month which shall be contributd to the Fund by each employer for each employee shall be according to any limit on the amount of wages and contributions as prescribed by the Board; and
- the amount of contributions for the month for the purpose of subsection 43(3) and section 44A is limited to any limit on the total contributions as prescribed by the Board.`
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

// Look for a particular section based on its name.
func SectionByName(name string) (Section, error) {
	for _, sec := range Sections {
		if sec.Name == name {
			return sec, nil
		}
	}
	return Section{}, errors.New("Invalid section.")
}

/*
Return all the rates within a particular section. These are the rates that are
listed within the table inside Third Schedule.
*/
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

// Return a particular rate for a given wages within the given section
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

// Return a particular section applicable to the given employee.
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

// Return a list of applicable sections
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

// Return a particular rate applicable to the given employee.
func (e *Employee) Rate() Rate {
	sec := e.Section()
	return sec.Rate(e.Wages)
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
	return Employee{
		NonMalaysian,
		contributeBefore1August1998,
		dateOfBirth,
		wages,
	}
}
