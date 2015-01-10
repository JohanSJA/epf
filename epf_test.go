package epf

import (
	"reflect"
	"testing"
	"time"
)

func TestRateContributionTotal0(t *testing.T) {
	rate := Rate{ContributionEmployer: 0.0, ContributionEmployee: 0.0}
	if rate.ContributionTotal() != 0.0+0.0 {
		t.Fail()
	}
}

func TestRateContributionTotal2400(t *testing.T) {
	rate := Rate{ContributionEmployer: 1300.0, ContributionEmployee: 1100.0}
	if rate.ContributionTotal() != 1300.0+1100.0 {
		t.Fail()
	}
}

func TestSectionByNameValid(t *testing.T) {
	name := "A"
	_, err := SectionByName(name)
	if err != nil {
		t.Fail()
	}
}

func TestSectionByNameInvalid(t *testing.T) {
	_, err := SectionByName("Invalid")
	if err == nil {
		t.Fail()
	}
}

func TestSectionARate(t *testing.T) {
	sec, err := SectionByName("A")
	if err != nil {
		t.Skip()
	}
	rate := sec.Rate(550.0)
	expected := 73.0
	if rate.ContributionEmployer != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployer)
	}
	rate = sec.Rate(25000.0)
	expected = 3000.0
	if rate.ContributionEmployer != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployer)
	}
}

func TestSectionBRate(t *testing.T) {
	sec, err := SectionByName("B")
	if err != nil {
		t.Skip()
	}
	rate := sec.Rate(720.0)
	expected := 80.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
}

func TestSectionCRate(t *testing.T) {
	sec, err := SectionByName("C")
	if err != nil {
		t.Skip()
	}
	rate := sec.Rate(1050.0)
	expected := 59.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
}

func TestSectionDRate(t *testing.T) {
	sec, err := SectionByName("D")
	if err != nil {
		t.Skip()
	}
	rate := sec.Rate(1150.0)
	expected := 64.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
}

func TestEmployeeSectionJuniorMalaysian(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-30, 0, 0)
	emp := NewEmployeeMalaysian(age, wages)
	section := emp.Section()
	if section.Name != "A" {
		t.Errorf("Expecting: %v , Gotten: %v", "A", section.Name)
	}
}

func TestEmployerSectionSeniorMalaysian(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-70, 0, 0)
	emp := NewEmployeeMalaysian(age, wages)
	section := emp.Section()
	if section.Name != "C" {
		t.Errorf("Expecting: %v , Gotten: %v", "C", section.Name)
	}
}

func TestEmployerSectionJuniorNonMalaysian(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-30, 0, 0)
	emp := NewEmployeeNonMalaysian(false, age, wages)
	section := emp.Section()
	if section.Name != "B" {
		t.Errorf("Expecting: %v , Gotten: %v", "B", section.Name)
	}
}

func TestEmployerSectionSeniorNonMalaysian(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-70, 0, 0)
	emp := NewEmployeeNonMalaysian(false, age, wages)
	section := emp.Section()
	if section.Name != "D" {
		t.Errorf("Expecting: %v , Gotten: %v", "D", section.Name)
	}
}

func TestEmployeeSectionJuniorPR(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-30, 0, 0)
	emp := NewEmployeePermanentResident(age, wages)
	section := emp.Section()
	if section.Name != "A" {
		t.Errorf("Expecting: %v , Gotten: %v", "A", section.Name)
	}
}

func TestEmployeeSectionSeniorPR(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-70, 0, 0)
	emp := NewEmployeePermanentResident(age, wages)
	section := emp.Section()
	if section.Name != "C" {
		t.Errorf("Expecting: %v , Gotten: %v", "C", section.Name)
	}
}

func TestEmployeeSectionsMalaysian(t *testing.T) {
	emp := Employee{Citizenship: Malaysian}
	sections := emp.Sections()
	expected := []*Section{&Sections[0], &Sections[2]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeUnknown(t *testing.T) {
	emp := Employee{}
	sections := emp.Sections()
	expected := []*Section{&Sections[0], &Sections[1], &Sections[2], &Sections[3]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeSectionsPR(t *testing.T) {
	emp := Employee{Citizenship: PermanentResident}
	sections := emp.Sections()
	expected := []*Section{&Sections[0], &Sections[2]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeSectionsNonMalaysian(t *testing.T) {
	emp := Employee{Citizenship: NonMalaysian}
	sections := emp.Sections()
	expected := []*Section{&Sections[1], &Sections[3]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeSectionsJunior(t *testing.T) {
	age := time.Now().AddDate(-30, 0, 0)
	emp := Employee{DateOfBirth: age}
	sections := emp.Sections()
	expected := []*Section{&Sections[0], &Sections[1]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeSectionsSenior(t *testing.T) {
	age := time.Now().AddDate(-65, 0, 0)
	emp := Employee{DateOfBirth: age}
	sections := emp.Sections()
	expected := []*Section{&Sections[2], &Sections[3]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeSectionsJuniorMalaysian(t *testing.T) {
	wages := 1500.0
	age := time.Now().AddDate(-30, 0, 0)
	emp := NewEmployeeMalaysian(age, wages)
	sections := emp.Sections()
	expected := []*Section{&Sections[0]}
	if !reflect.DeepEqual(sections, expected) {
		t.Logf("Expecting: %v", expected)
		t.Logf("Gotten: %v", sections)
		t.Fail()
	}
}

func TestEmployeeRate(t *testing.T) {
	empAge := time.Now().AddDate(-30, 0, 0)
	emp := NewEmployeeMalaysian(empAge, 1500.0)
	rate := emp.Rate()
	expected := 165.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
}
