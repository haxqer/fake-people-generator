package generator

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"

	"fake-people-generator/internal/dataset"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderRandom Gender = "random"
)

type Filters struct {
	Gender    Gender
	AgeMin    *int
	AgeMax    *int
	StateCode string
}

type GenerateInput struct {
	Country       string
	Count         int
	Filters       Filters
	IncludeFields []string
}

type Result struct {
	Country     string    `json:"country"`
	Count       int       `json:"count"`
	GeneratedAt time.Time `json:"generated_at"`
	Profiles    []Profile `json:"profiles"`
}

type Profile struct {
	ID            string     `json:"id"`
	Country       string     `json:"country"`
	Gender        Gender     `json:"gender"`
	FirstName     string     `json:"first_name"`
	MiddleName    string     `json:"middle_name"`
	LastName      string     `json:"last_name"`
	FullName      string     `json:"full_name"`
	Age           int        `json:"age"`
	BirthDate     string     `json:"birth_date"`
	MaritalStatus string     `json:"marital_status"`
	Occupation    string     `json:"occupation"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	TimeZone      string     `json:"time_zone"`
	Employment    Employment `json:"employment"`
	Education     *Education `json:"education,omitempty"`
	Sensitive     *Sensitive `json:"sensitive,omitempty"`
	Phone         Phone      `json:"phone"`
	Address       Address    `json:"address"`
}

type Employment struct {
	Status        string `json:"status"`
	CompanyName   string `json:"company_name,omitempty"`
	CompanyDomain string `json:"company_domain,omitempty"`
	JobTitle      string `json:"job_title,omitempty"`
	Department    string `json:"department,omitempty"`
	WorkEmail     string `json:"work_email,omitempty"`
	Industry      string `json:"industry,omitempty"`
	Seniority     string `json:"seniority,omitempty"`
	StartDate     string `json:"start_date,omitempty"`
	WorkMode      string `json:"work_mode,omitempty"`
	ManagerName   string `json:"manager_name,omitempty"`
	OfficeCity    string `json:"office_city,omitempty"`
	OfficeState   string `json:"office_state,omitempty"`
}

type Phone struct {
	E164     string `json:"e164"`
	National string `json:"national"`
	AreaCode string `json:"area_code"`
}

type Address struct {
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state"`
	StateCode   string `json:"state_code"`
	PostalCode  string `json:"postal_code"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Generator struct {
	data *dataset.Data
	now  func() time.Time
}

var ageBands = [...]struct {
	Min    int
	Max    int
	Weight int
}{
	{Min: 0, Max: 12, Weight: 15},
	{Min: 13, Max: 17, Weight: 7},
	{Min: 18, Max: 24, Weight: 10},
	{Min: 25, Max: 34, Weight: 18},
	{Min: 35, Max: 44, Weight: 16},
	{Min: 45, Max: 54, Weight: 14},
	{Min: 55, Max: 64, Weight: 10},
	{Min: 65, Max: 79, Weight: 8},
	{Min: 80, Max: 95, Weight: 2},
}

func New(now func() time.Time) *Generator {
	if now == nil {
		now = time.Now
	}

	return &Generator{
		data: dataset.US(),
		now:  now,
	}
}

func (g *Generator) Generate(input GenerateInput) (Result, error) {
	return g.GenerateContext(context.Background(), input)
}

func (g *Generator) GenerateContext(ctx context.Context, input GenerateInput) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if strings.ToUpper(strings.TrimSpace(input.Country)) != "US" {
		return Result{}, &ValidationError{Message: "country must be US"}
	}
	if input.Count < 1 || input.Count > 100 {
		return Result{}, &ValidationError{Message: "count must be between 1 and 100"}
	}

	filters, err := normalizeFilters(input.Filters, g.data)
	if err != nil {
		return Result{}, err
	}
	selection, err := normalizeIncludeFields(input.IncludeFields)
	if err != nil {
		return Result{}, err
	}

	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	generatedAt := g.now().UTC().Truncate(time.Second)
	rng, err := newRandom()
	if err != nil {
		return Result{}, fmt.Errorf("create random source: %w", err)
	}

	profiles := make([]Profile, 0, input.Count)
	for i := 0; i < input.Count; i++ {
		if err := ctx.Err(); err != nil {
			return Result{}, err
		}

		profile, err := g.generateProfile(rng, generatedAt, filters, selection)
		if err != nil {
			return Result{}, err
		}
		profiles = append(profiles, profile)
	}

	return Result{
		Country:     "US",
		Count:       len(profiles),
		GeneratedAt: generatedAt,
		Profiles:    profiles,
	}, nil
}

func normalizeFilters(filters Filters, data *dataset.Data) (Filters, error) {
	normalized := filters

	switch normalized.Gender {
	case "", GenderRandom:
		normalized.Gender = GenderRandom
	case GenderMale, GenderFemale:
	default:
		return Filters{}, &ValidationError{Message: "filters.gender must be one of male, female, or random"}
	}

	if (normalized.AgeMin == nil) != (normalized.AgeMax == nil) {
		return Filters{}, &ValidationError{Message: "filters.age_min and filters.age_max must be provided together"}
	}
	if normalized.AgeMin != nil && normalized.AgeMax != nil {
		if *normalized.AgeMin < 0 || *normalized.AgeMax < 0 || *normalized.AgeMin > 100 || *normalized.AgeMax > 100 {
			return Filters{}, &ValidationError{Message: "filters.age_min and filters.age_max must be between 0 and 100"}
		}
		if *normalized.AgeMin > *normalized.AgeMax {
			return Filters{}, &ValidationError{Message: "filters.age_min must be less than or equal to filters.age_max"}
		}
	}

	if normalized.StateCode != "" {
		normalized.StateCode = strings.ToUpper(strings.TrimSpace(normalized.StateCode))
		if _, ok := data.StatesByCode[normalized.StateCode]; !ok {
			return Filters{}, &ValidationError{Message: "filters.state must be a supported US state code"}
		}
	}

	return normalized, nil
}

func (g *Generator) generateProfile(rng *mathrand.Rand, generatedAt time.Time, filters Filters, selection fieldSelection) (Profile, error) {
	state, err := g.selectState(rng, filters)
	if err != nil {
		return Profile{}, err
	}

	city, err := selectCity(rng, state.Cities)
	if err != nil {
		return Profile{}, err
	}

	actualGender := resolveGender(rng, filters.Gender)
	firstName := g.selectFirstName(rng, actualGender)
	middleName := g.selectDistinctFirstName(rng, actualGender, firstName)
	lastName := selectWeightedValue(rng, g.data.LastNames)
	firstHandle := normalizeHandlePart(firstName)
	lastHandle := normalizeHandlePart(lastName)
	age := selectAge(rng, filters)
	birthDate := randomBirthDate(rng, generatedAt, age)
	birthDateText := birthDate.Format("2006-01-02")
	areaCode := city.AreaCodes[rng.Intn(len(city.AreaCodes))]

	profileID := newUUID(rng)
	phone := generatePhone(rng, areaCode)
	address := Address{
		StreetLine1: generateStreetAddress(rng, g.data),
		StreetLine2: generateAddressLine2(rng, g.data),
		City:        city.Name,
		State:       state.Name,
		StateCode:   state.Code,
		PostalCode:  city.PostalCode,
	}
	username := generateUsername(rng, firstHandle, lastHandle, birthDateText)
	email := generateEmail(rng, username, firstHandle, lastHandle, g.data)
	occupation := selectOccupation(rng, age, g.data)
	var educationRecord educationRecord
	if selection.wantsEducation() || selection.EmploymentStartDate {
		educationRecord = generateEducationRecord(rng, generatedAt, age, occupation, firstHandle, lastHandle, state, city)
	}
	employment := g.generateEmployment(rng, age, occupation, firstName, lastName, firstHandle, lastHandle, state, city, selection, educationRecord, generatedAt)
	education := buildEducation(selection, educationRecord)
	sensitive := buildSensitive(selection, rng, generatedAt, state)

	return Profile{
		ID:            profileID,
		Country:       "US",
		Gender:        actualGender,
		FirstName:     firstName,
		MiddleName:    middleName,
		LastName:      lastName,
		FullName:      buildFullName(firstName, middleName, lastName),
		Age:           calculateAge(birthDate, generatedAt),
		BirthDate:     birthDateText,
		MaritalStatus: selectMaritalStatus(rng, age),
		Occupation:    occupation,
		Username:      username,
		Email:         email,
		TimeZone:      selectTimeZone(rng, state.Code),
		Employment:    employment,
		Education:     education,
		Sensitive:     sensitive,
		Phone:         phone,
		Address:       address,
	}, nil
}

func (g *Generator) selectState(rng *mathrand.Rand, filters Filters) (dataset.State, error) {
	if filters.StateCode != "" {
		state, ok := g.data.StatesByCode[filters.StateCode]
		if !ok {
			return dataset.State{}, &ValidationError{Message: "filters.state must be a supported US state code"}
		}
		return state, nil
	}

	if len(g.data.States) == 0 {
		return dataset.State{}, fmt.Errorf("state dataset is empty")
	}

	index, err := pickWeightedIndex(len(g.data.States), func(i int) int {
		return g.data.States[i].Weight
	}, rng)
	if err != nil {
		return dataset.State{}, err
	}

	return g.data.States[index], nil
}

func (g *Generator) selectFirstName(rng *mathrand.Rand, gender Gender) string {
	if gender == GenderFemale {
		return selectWeightedValue(rng, g.data.FemaleFirstName)
	}
	return selectWeightedValue(rng, g.data.MaleFirstNames)
}

func resolveGender(rng *mathrand.Rand, requested Gender) Gender {
	if requested == GenderMale || requested == GenderFemale {
		return requested
	}
	if rng.Intn(2) == 0 {
		return GenderMale
	}
	return GenderFemale
}

func selectAge(rng *mathrand.Rand, filters Filters) int {
	if filters.AgeMin != nil && filters.AgeMax != nil {
		if *filters.AgeMin == *filters.AgeMax {
			return *filters.AgeMin
		}
		return *filters.AgeMin + rng.Intn(*filters.AgeMax-*filters.AgeMin+1)
	}

	index, err := pickWeightedIndex(len(ageBands), func(i int) int {
		return ageBands[i].Weight
	}, rng)
	if err != nil {
		return 30
	}

	band := ageBands[index]
	if band.Min == band.Max {
		return band.Min
	}
	return band.Min + rng.Intn(band.Max-band.Min+1)
}

func selectCity(rng *mathrand.Rand, cities []dataset.City) (dataset.City, error) {
	if len(cities) == 0 {
		return dataset.City{}, fmt.Errorf("city dataset is empty")
	}

	index, err := pickWeightedIndex(len(cities), func(i int) int {
		return cities[i].Weight
	}, rng)
	if err != nil {
		return dataset.City{}, err
	}

	return cities[index], nil
}

func selectWeightedValue(rng *mathrand.Rand, values []dataset.WeightedValue) string {
	index, err := pickWeightedIndex(len(values), func(i int) int {
		return values[i].Weight
	}, rng)
	if err != nil {
		return ""
	}
	return values[index].Value
}

func pickWeightedIndex(length int, weightFn func(int) int, rng *mathrand.Rand) (int, error) {
	if length == 0 {
		return 0, fmt.Errorf("weighted set is empty")
	}

	total := 0
	for i := 0; i < length; i++ {
		weight := weightFn(i)
		if weight < 0 {
			return 0, fmt.Errorf("weight cannot be negative")
		}
		total += weight
	}
	if total == 0 {
		return 0, fmt.Errorf("weight total must be greater than zero")
	}

	target := rng.Intn(total)
	running := 0
	for i := 0; i < length; i++ {
		running += weightFn(i)
		if target < running {
			return i, nil
		}
	}

	return length - 1, nil
}

func randomBirthDate(rng *mathrand.Rand, generatedAt time.Time, age int) time.Time {
	referenceDate := time.Date(generatedAt.Year(), generatedAt.Month(), generatedAt.Day(), 0, 0, 0, 0, time.UTC)
	earliest := referenceDate.AddDate(-(age + 1), 0, 1)
	latest := referenceDate.AddDate(-age, 0, 0)

	daySpan := int(latest.Sub(earliest).Hours() / 24)
	if daySpan <= 0 {
		return latest
	}

	return earliest.AddDate(0, 0, rng.Intn(daySpan+1))
}

func calculateAge(birthDate time.Time, reference time.Time) int {
	birthDate = birthDate.UTC()
	reference = reference.UTC()

	age := reference.Year() - birthDate.Year()
	if reference.Month() < birthDate.Month() || (reference.Month() == birthDate.Month() && reference.Day() < birthDate.Day()) {
		age--
	}
	return age
}

func generateStreetAddress(rng *mathrand.Rand, data *dataset.Data) string {
	number := 100 + rng.Intn(9900)
	streetName := selectWeightedValue(rng, data.StreetNames)
	suffix := selectWeightedValue(rng, data.StreetSuffixes)
	return fmt.Sprintf("%d %s %s", number, streetName, suffix)
}

func generatePhone(rng *mathrand.Rand, areaCode string) Phone {
	exchange := 200 + rng.Intn(800)
	for exchange%100 == 11 {
		exchange = 200 + rng.Intn(800)
	}
	subscriber := rng.Intn(10000)

	return Phone{
		E164:     fmt.Sprintf("+1%s%03d%04d", areaCode, exchange, subscriber),
		National: fmt.Sprintf("(%s) %03d-%04d", areaCode, exchange, subscriber),
		AreaCode: areaCode,
	}
}

func newRandom() (*mathrand.Rand, error) {
	var seedBytes [8]byte
	if _, err := rand.Read(seedBytes[:]); err != nil {
		return nil, err
	}

	seed := int64(0)
	for _, b := range seedBytes {
		seed = (seed << 8) | int64(b)
	}
	return mathrand.New(mathrand.NewSource(seed)), nil
}

func newUUID(rng *mathrand.Rand) string {
	var raw [16]byte
	for i := range raw {
		raw[i] = byte(rng.Intn(256))
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	var encoded [36]byte
	hex.Encode(encoded[0:8], raw[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], raw[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], raw[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], raw[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], raw[10:16])
	return string(encoded[:])
}
