package generator

import (
	"fmt"
	mathrand "math/rand"
	"strconv"
	"strings"
	"time"

	"fake-people-generator/internal/dataset"
)

type Sensitive struct {
	SSN           string         `json:"ssn,omitempty"`
	DriverLicense *DriverLicense `json:"driver_license,omitempty"`
	Passport      *Passport      `json:"passport,omitempty"`
	BankCard      *BankCard      `json:"bank_card,omitempty"`
}

type DriverLicense struct {
	Number    string `json:"number,omitempty"`
	StateCode string `json:"state_code,omitempty"`
}

type Passport struct {
	Number  string `json:"number,omitempty"`
	Country string `json:"country,omitempty"`
}

type BankCard struct {
	Brand       string `json:"brand,omitempty"`
	Number      string `json:"number,omitempty"`
	ExpiryMonth int    `json:"expiry_month,omitempty"`
	ExpiryYear  int    `json:"expiry_year,omitempty"`
	CVV         string `json:"cvv,omitempty"`
}

func buildSensitive(selection fieldSelection, rng *mathrand.Rand, generatedAt time.Time, state dataset.State) *Sensitive {
	if !selection.wantsSensitive() {
		return nil
	}

	sensitive := &Sensitive{}

	if selection.SensitiveSSN {
		sensitive.SSN = generateSSN(rng)
	}

	if selection.wantsDriverLicense() {
		driverLicense := generateDriverLicense(rng, state.Code)
		sensitive.DriverLicense = &DriverLicense{}
		if selection.SensitiveDriverLicenseNo {
			sensitive.DriverLicense.Number = driverLicense.Number
		}
		if selection.SensitiveDriverLicenseState {
			sensitive.DriverLicense.StateCode = driverLicense.StateCode
		}
	}

	if selection.wantsPassport() {
		passport := generatePassport(rng)
		sensitive.Passport = &Passport{}
		if selection.SensitivePassportNumber {
			sensitive.Passport.Number = passport.Number
		}
		if selection.SensitivePassportCountry {
			sensitive.Passport.Country = passport.Country
		}
	}

	if selection.wantsBankCard() {
		card := generateBankCard(rng, generatedAt)
		sensitive.BankCard = &BankCard{}
		if selection.SensitiveBankCardBrand {
			sensitive.BankCard.Brand = card.Brand
		}
		if selection.SensitiveBankCardNumber {
			sensitive.BankCard.Number = card.Number
		}
		if selection.SensitiveBankCardMonth {
			sensitive.BankCard.ExpiryMonth = card.ExpiryMonth
		}
		if selection.SensitiveBankCardYear {
			sensitive.BankCard.ExpiryYear = card.ExpiryYear
		}
		if selection.SensitiveBankCardCVV {
			sensitive.BankCard.CVV = card.CVV
		}
	}

	return sensitive
}

func generateSSN(rng *mathrand.Rand) string {
	area := 1 + rng.Intn(899)
	for area == 666 {
		area = 1 + rng.Intn(899)
	}
	group := 1 + rng.Intn(99)
	serial := 1 + rng.Intn(9999)
	return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
}

func generateDriverLicense(rng *mathrand.Rand, stateCode string) DriverLicense {
	switch stateCode {
	case "CA", "OR":
		return DriverLicense{Number: randLetters(rng, 1) + randDigits(rng, 7), StateCode: stateCode}
	case "NY", "CO", "KS", "ND", "SD":
		return DriverLicense{Number: randDigits(rng, 9), StateCode: stateCode}
	case "TX", "PA", "VA", "AZ", "NV", "UT", "NM", "LA", "MS", "AL", "AR", "TN", "KY", "IN", "OH":
		return DriverLicense{Number: randDigits(rng, 8), StateCode: stateCode}
	case "FL", "GA", "NC", "SC", "WA", "IL", "MI", "MN", "MA":
		return DriverLicense{Number: randLetters(rng, 1) + randDigits(rng, 11), StateCode: stateCode}
	case "NJ", "CT", "RI", "VT", "NH", "ME", "DE", "MD":
		return DriverLicense{Number: randLetters(rng, 1) + randDigits(rng, 8), StateCode: stateCode}
	default:
		return DriverLicense{Number: randLetters(rng, 1) + randDigits(rng, 7), StateCode: stateCode}
	}
}

func generatePassport(rng *mathrand.Rand) Passport {
	return Passport{
		Number:  randDigits(rng, 9),
		Country: "US",
	}
}

func generateBankCard(rng *mathrand.Rand, generatedAt time.Time) BankCard {
	type cardTemplate struct {
		Brand   string
		Prefix  string
		Length  int
		CVVSize int
	}

	templates := []cardTemplate{
		{Brand: "Visa", Prefix: "4", Length: 16, CVVSize: 3},
		{Brand: "Mastercard", Prefix: strconv.Itoa(51 + rng.Intn(5)), Length: 16, CVVSize: 3},
		{Brand: "American Express", Prefix: pickFromSet(rng, []string{"34", "37"}), Length: 15, CVVSize: 4},
		{Brand: "Discover", Prefix: "6011", Length: 16, CVVSize: 3},
	}

	template := templates[rng.Intn(len(templates))]
	number := generateLuhnNumber(rng, template.Prefix, template.Length)
	expiryMonth := 1 + rng.Intn(12)
	expiryYear := generatedAt.Year() + 1 + rng.Intn(6)

	return BankCard{
		Brand:       template.Brand,
		Number:      number,
		ExpiryMonth: expiryMonth,
		ExpiryYear:  expiryYear,
		CVV:         randDigits(rng, template.CVVSize),
	}
}

func generateLuhnNumber(rng *mathrand.Rand, prefix string, totalLength int) string {
	var builder strings.Builder
	builder.WriteString(prefix)

	for builder.Len() < totalLength-1 {
		builder.WriteString(randDigits(rng, 1))
	}

	base := builder.String()
	checkDigit := luhnCheckDigit(base)
	return base + strconv.Itoa(checkDigit)
}

func luhnCheckDigit(number string) int {
	sum := 0
	double := true

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

	return (10 - (sum % 10)) % 10
}

func randDigits(rng *mathrand.Rand, count int) string {
	var builder strings.Builder
	builder.Grow(count)
	for i := 0; i < count; i++ {
		builder.WriteByte(byte('0' + rng.Intn(10)))
	}
	return builder.String()
}

func randLetters(rng *mathrand.Rand, count int) string {
	var builder strings.Builder
	builder.Grow(count)
	for i := 0; i < count; i++ {
		builder.WriteByte(byte('A' + rng.Intn(26)))
	}
	return builder.String()
}
