package epf

import (
	"testing"
	"time"
)

func TestRateContributionTotal(t *testing.T) {
	rate := Rate{ContributionEmployer: 0.0, ContributionEmployee: 0.0}
	if rate.ContributionTotal() != 0.0+0.0 {
		t.Fail()
	}
	rate = Rate{ContributionEmployer: 1300.0, ContributionEmployee: 1100.0}
	if rate.ContributionTotal() != 1300.0+1100.0 {
		t.Fail()
	}
}

func TestSectionRate(t *testing.T) {
	rate := SectionRate("A", 550.0)
	expected := 73.0
	if rate.ContributionEmployer != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployer)
	}
	rate = SectionRate("B", 720.0)
	expected = 80.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
	rate = SectionRate("C", 1050.0)
	expected = 59.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
	rate = SectionRate("D", 1150.0)
	expected = 64.0
	if rate.ContributionEmployee != expected {
		t.Errorf("Expecting: %v , Gotten: %v", expected, rate.ContributionEmployee)
	}
}

func TestEmployeeSection(t *testing.T) {
	wages := 1500.0
	seniorAge := time.Now().AddDate(-65, 0, 0)
	juniorAge := time.Now().AddDate(-30, 0, 0)
	juniorMalaysian := NewEmployeeMalaysian(juniorAge, wages)
	section := juniorMalaysian.Section()
	if section.Name != "A" {
		t.Errorf("Expecting: %v , Gotten: %v", "A", section.Name)
	}
	seniorMalaysian := NewEmployeeMalaysian(seniorAge, wages)
	section = seniorMalaysian.Section()
	if section.Name != "C" {
		t.Errorf("Expecting: %v , Gotten: %v", "C", section.Name)
	}
	juniorNonMalaysian := NewEmployeeNonMalaysian(false, juniorAge, wages)
	section = juniorNonMalaysian.Section()
	if section.Name != "B" {
		t.Errorf("Expecting: %v , Gotten: %v", "B", section.Name)
	}
	seniorNonMalaysian := NewEmployeeNonMalaysian(false, seniorAge, wages)
	section = seniorNonMalaysian.Section()
	if section.Name != "D" {
		t.Errorf("Expecting: %v , Gotten: %v", "D", section.Name)
	}
	juniorPR := NewEmployeePermanentResident(juniorAge, wages)
	section = juniorPR.Section()
	if section.Name != "A" {
		t.Errorf("Expecting: %v , Gotten: %v", "A", section.Name)
	}
	seniorPR := NewEmployeePermanentResident(seniorAge, wages)
	section = seniorPR.Section()
	if section.Name != "C" {
		t.Errorf("Expecting: %v , Gotten: %v", "C", section.Name)
	}
}
