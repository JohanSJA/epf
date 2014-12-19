package epf

import (
	"testing"
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
