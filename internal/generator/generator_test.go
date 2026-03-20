package generator_test

import (
	"testing"
	"time"

	"fake-people-generator/internal/generator"
)

func TestGenerateConsistency(t *testing.T) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ageMin := 18
	ageMax := 65
	result, err := gen.Generate(generator.GenerateInput{
		Country: "US",
		Count:   100,
		Filters: generator.Filters{
			Gender:    generator.GenderRandom,
			AgeMin:    &ageMin,
			AgeMax:    &ageMax,
			StateCode: "NY",
		},
	})
	if err != nil {
		t.Fatalf("generate profiles: %v", err)
	}

	if len(result.Profiles) != 100 {
		t.Fatalf("expected 100 profiles, got %d", len(result.Profiles))
	}

	for _, profile := range result.Profiles {
		if profile.Country != "US" {
			t.Fatalf("expected US country, got %s", profile.Country)
		}
		if profile.Address.StateCode != "NY" {
			t.Fatalf("expected NY state code, got %s", profile.Address.StateCode)
		}
		if profile.Age < ageMin || profile.Age > ageMax {
			t.Fatalf("expected age within range, got %d", profile.Age)
		}
		birthDate, err := time.Parse("2006-01-02", profile.BirthDate)
		if err != nil {
			t.Fatalf("parse birth date: %v", err)
		}
		if age := calcAge(birthDate, result.GeneratedAt); age != profile.Age {
			t.Fatalf("expected age %d to match birth date %s, got %d", profile.Age, profile.BirthDate, age)
		}
		if profile.Phone.AreaCode == "" {
			t.Fatal("expected area code")
		}
		if profile.Address.City == "" || profile.Address.PostalCode == "" {
			t.Fatal("expected city and postal code")
		}
	}
}

func calcAge(birthDate time.Time, reference time.Time) int {
	age := reference.Year() - birthDate.Year()
	if reference.Month() < birthDate.Month() || (reference.Month() == birthDate.Month() && reference.Day() < birthDate.Day()) {
		age--
	}
	return age
}
