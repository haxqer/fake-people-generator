package generator_test

import (
	"context"
	"errors"
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
		if profile.MiddleName == "" {
			t.Fatal("expected middle name")
		}
		if profile.MaritalStatus == "" {
			t.Fatal("expected marital status")
		}
		if profile.Occupation == "" {
			t.Fatal("expected occupation")
		}
		if profile.Username == "" {
			t.Fatal("expected username")
		}
		if profile.Email == "" {
			t.Fatal("expected email")
		}
		if profile.TimeZone == "" {
			t.Fatal("expected time zone")
		}
		if profile.Education != nil {
			t.Fatal("expected education to be omitted by default")
		}
		if profile.Sensitive != nil {
			t.Fatal("expected sensitive to be omitted by default")
		}
		if profile.Employment.Status != "employed" {
			t.Fatalf("expected employed status, got %s", profile.Employment.Status)
		}
		if profile.Employment.CompanyName == "" {
			t.Fatal("expected company name")
		}
		if profile.Employment.JobTitle == "" {
			t.Fatal("expected job title")
		}
		if profile.Employment.CompanyDomain == "" {
			t.Fatal("expected company domain")
		}
		if profile.Employment.WorkEmail == "" {
			t.Fatal("expected work email")
		}
		if profile.Employment.Industry != "" || profile.Employment.StartDate != "" {
			t.Fatalf("expected employment enhancements to be omitted by default, got %+v", profile.Employment)
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

func TestGenerateMinorOccupation(t *testing.T) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ageMin := 10
	ageMax := 10
	result, err := gen.Generate(generator.GenerateInput{
		Country: "US",
		Count:   1,
		Filters: generator.Filters{
			AgeMin: &ageMin,
			AgeMax: &ageMax,
		},
	})
	if err != nil {
		t.Fatalf("generate minor profile: %v", err)
	}

	if got := result.Profiles[0].Occupation; got != "Student" {
		t.Fatalf("expected Student occupation for minor, got %s", got)
	}
	if got := result.Profiles[0].MaritalStatus; got != "single" {
		t.Fatalf("expected single marital status for minor, got %s", got)
	}
	if got := result.Profiles[0].Employment.Status; got != "student" {
		t.Fatalf("expected student employment status for minor, got %s", got)
	}
	if got := result.Profiles[0].Employment.CompanyName; got != "" {
		t.Fatalf("expected empty company name for minor, got %s", got)
	}
}

func TestGenerateWithIncludeFields(t *testing.T) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ageMin := 31
	ageMax := 31
	result, err := gen.Generate(generator.GenerateInput{
		Country: "US",
		Count:   1,
		Filters: generator.Filters{
			Gender:    generator.GenderMale,
			AgeMin:    &ageMin,
			AgeMax:    &ageMax,
			StateCode: "TX",
		},
		IncludeFields: []string{
			"education.level",
			"education.school_name",
			"education.graduation_year",
			"employment.industry",
			"employment.seniority",
			"employment.start_date",
			"employment.work_mode",
			"employment.manager_name",
			"employment.office_city",
			"employment.office_state",
			"sensitive.ssn",
			"sensitive.driver_license.number",
			"sensitive.driver_license.state_code",
			"sensitive.passport.number",
			"sensitive.passport.country",
			"sensitive.bank_card.brand",
			"sensitive.bank_card.number",
			"sensitive.bank_card.expiry_month",
			"sensitive.bank_card.expiry_year",
			"sensitive.bank_card.cvv",
		},
	})
	if err != nil {
		t.Fatalf("generate profiles with include fields: %v", err)
	}

	profile := result.Profiles[0]
	if profile.Education == nil || profile.Education.Level == "" || profile.Education.SchoolName == "" {
		t.Fatalf("expected education payload, got %+v", profile.Education)
	}
	if profile.Employment.Industry == "" || profile.Employment.Seniority == "" || profile.Employment.StartDate == "" {
		t.Fatalf("expected employment enhancement payload, got %+v", profile.Employment)
	}
	if profile.Employment.ManagerName == "" || profile.Employment.OfficeCity == "" || profile.Employment.OfficeState == "" {
		t.Fatalf("expected management/office payload, got %+v", profile.Employment)
	}
	if profile.Sensitive == nil || profile.Sensitive.SSN == "" {
		t.Fatal("expected sensitive payload")
	}
	if profile.Sensitive.DriverLicense == nil || profile.Sensitive.DriverLicense.StateCode != "TX" {
		t.Fatalf("expected TX driver license, got %+v", profile.Sensitive.DriverLicense)
	}
	if profile.Sensitive.Passport == nil || profile.Sensitive.Passport.Country != "US" {
		t.Fatalf("expected US passport, got %+v", profile.Sensitive.Passport)
	}
	if profile.Sensitive.BankCard == nil || !isLuhn(profile.Sensitive.BankCard.Number) {
		t.Fatalf("expected Luhn-valid card, got %+v", profile.Sensitive.BankCard)
	}
}

func TestGenerateInvalidIncludeField(t *testing.T) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	_, err := gen.Generate(generator.GenerateInput{
		Country:       "US",
		Count:         1,
		IncludeFields: []string{"education.unknown"},
	})
	if err == nil {
		t.Fatal("expected validation error for invalid include_fields")
	}
}

func TestGenerateContextCanceled(t *testing.T) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := gen.GenerateContext(ctx, generator.GenerateInput{
		Country: "US",
		Count:   100,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func calcAge(birthDate time.Time, reference time.Time) int {
	age := reference.Year() - birthDate.Year()
	if reference.Month() < birthDate.Month() || (reference.Month() == birthDate.Month() && reference.Day() < birthDate.Day()) {
		age--
	}
	return age
}

func isLuhn(number string) bool {
	sum := 0
	double := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return sum%10 == 0
}
