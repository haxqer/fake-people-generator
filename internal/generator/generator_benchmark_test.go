package generator_test

import (
	"testing"
	"time"

	"fake-people-generator/internal/generator"
)

func BenchmarkGenerateDefaultBatch(b *testing.B) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ageMin := 18
	ageMax := 65
	input := generator.GenerateInput{
		Country: "US",
		Count:   100,
		Filters: generator.Filters{
			Gender:    generator.GenderRandom,
			AgeMin:    &ageMin,
			AgeMax:    &ageMax,
			StateCode: "CA",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := gen.Generate(input); err != nil {
			b.Fatalf("generate default batch: %v", err)
		}
	}
}

func BenchmarkGenerateExpandedBatch(b *testing.B) {
	gen := generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	})

	ageMin := 25
	ageMax := 45
	input := generator.GenerateInput{
		Country: "US",
		Count:   100,
		Filters: generator.Filters{
			Gender:    generator.GenderRandom,
			AgeMin:    &ageMin,
			AgeMax:    &ageMax,
			StateCode: "TX",
		},
		IncludeFields: []string{
			"education.level",
			"education.school_name",
			"education.major",
			"education.graduation_year",
			"employment.industry",
			"employment.seniority",
			"employment.start_date",
			"employment.work_mode",
			"sensitive.ssn",
			"sensitive.driver_license.number",
			"sensitive.driver_license.state_code",
			"sensitive.bank_card.brand",
			"sensitive.bank_card.number",
			"sensitive.bank_card.expiry_month",
			"sensitive.bank_card.expiry_year",
			"sensitive.bank_card.cvv",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := gen.Generate(input); err != nil {
			b.Fatalf("generate expanded batch: %v", err)
		}
	}
}
