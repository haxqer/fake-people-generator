package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fake-people-generator/internal/api"
	"fake-people-generator/internal/generator"
)

func TestGenerateProfilesSuccess(t *testing.T) {
	handler := api.NewServer(generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	}))

	body := `{"country":"US","count":3,"filters":{"gender":"female","age_min":30,"age_max":30,"state":"CA"}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/profiles/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Meta struct {
			Country     string `json:"country"`
			Count       int    `json:"count"`
			GeneratedAt string `json:"generated_at"`
		} `json:"meta"`
		Profiles []generator.Profile `json:"profiles"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Meta.Country != "US" {
		t.Fatalf("expected country US, got %s", resp.Meta.Country)
	}
	if resp.Meta.Count != 3 {
		t.Fatalf("expected count 3, got %d", resp.Meta.Count)
	}
	if len(resp.Profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(resp.Profiles))
	}

	generatedAt, err := time.Parse(time.RFC3339, resp.Meta.GeneratedAt)
	if err != nil {
		t.Fatalf("parse generated_at: %v", err)
	}

	for _, profile := range resp.Profiles {
		if profile.Gender != generator.GenderFemale {
			t.Fatalf("expected female profile, got %s", profile.Gender)
		}
		if profile.MiddleName == "" {
			t.Fatal("expected middle_name")
		}
		if profile.FullName == "" || !strings.Contains(profile.FullName, profile.MiddleName) {
			t.Fatalf("expected full_name to include middle name, got %s", profile.FullName)
		}
		if profile.Age != 30 {
			t.Fatalf("expected age 30, got %d", profile.Age)
		}
		if profile.MaritalStatus == "" {
			t.Fatal("expected marital_status")
		}
		if profile.Occupation == "" {
			t.Fatal("expected occupation")
		}
		if profile.Username == "" {
			t.Fatal("expected username")
		}
		if !strings.Contains(profile.Email, "@") {
			t.Fatalf("expected email, got %s", profile.Email)
		}
		if profile.TimeZone == "" {
			t.Fatal("expected time_zone")
		}
		if profile.Employment.Status != "employed" {
			t.Fatalf("expected employed status, got %s", profile.Employment.Status)
		}
		if profile.Employment.CompanyName == "" {
			t.Fatal("expected company_name")
		}
		if profile.Employment.JobTitle == "" {
			t.Fatal("expected job_title")
		}
		if profile.Employment.Industry != "" {
			t.Fatalf("expected industry to be omitted by default, got %s", profile.Employment.Industry)
		}
		if profile.Education != nil {
			t.Fatal("expected education to be omitted by default")
		}
		if profile.Sensitive != nil {
			t.Fatal("expected sensitive to be omitted by default")
		}
		if profile.Address.StateCode != "CA" {
			t.Fatalf("expected CA state code, got %s", profile.Address.StateCode)
		}
		birthDate, err := time.Parse("2006-01-02", profile.BirthDate)
		if err != nil {
			t.Fatalf("parse birth_date: %v", err)
		}
		if age := calcAge(birthDate, generatedAt); age != profile.Age {
			t.Fatalf("expected consistent age, got %d and %d", age, profile.Age)
		}
		if !strings.HasPrefix(profile.Phone.E164, "+1") {
			t.Fatalf("expected E164 phone, got %s", profile.Phone.E164)
		}
	}
}

func TestGenerateProfilesIncludeFields(t *testing.T) {
	handler := api.NewServer(generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	}))

	body := `{
		"country":"US",
		"count":1,
		"filters":{"gender":"female","age_min":31,"age_max":31,"state":"CA"},
		"include_fields":[
			"education.level",
			"education.school_name",
			"education.major",
			"education.graduation_year",
			"education.student_status",
			"education.school_email",
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
			"sensitive.bank_card.cvv"
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/profiles/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Profiles []generator.Profile `json:"profiles"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(resp.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(resp.Profiles))
	}

	profile := resp.Profiles[0]
	if profile.Education == nil {
		t.Fatal("expected education")
	}
	if profile.Education.Level == "" || profile.Education.SchoolName == "" {
		t.Fatalf("expected populated education, got %+v", profile.Education)
	}
	if profile.Employment.Industry == "" || profile.Employment.Seniority == "" || profile.Employment.StartDate == "" {
		t.Fatalf("expected employment enhancements, got %+v", profile.Employment)
	}
	if profile.Employment.ManagerName == "" || profile.Employment.OfficeCity == "" || profile.Employment.OfficeState == "" {
		t.Fatalf("expected office fields, got %+v", profile.Employment)
	}
	if profile.Sensitive == nil || profile.Sensitive.SSN == "" {
		t.Fatal("expected sensitive.ssn")
	}
	if profile.Sensitive.DriverLicense == nil || profile.Sensitive.DriverLicense.Number == "" || profile.Sensitive.DriverLicense.StateCode != "CA" {
		t.Fatalf("expected state-linked driver license, got %+v", profile.Sensitive.DriverLicense)
	}
	if profile.Sensitive.Passport == nil || profile.Sensitive.Passport.Country != "US" {
		t.Fatalf("expected passport country US, got %+v", profile.Sensitive.Passport)
	}
	if profile.Sensitive.BankCard == nil || profile.Sensitive.BankCard.Number == "" || profile.Sensitive.BankCard.CVV == "" {
		t.Fatalf("expected bank card, got %+v", profile.Sensitive.BankCard)
	}
}

func TestGenerateProfilesValidationErrors(t *testing.T) {
	handler := api.NewServer(generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	}))

	tests := []struct {
		name string
		body string
	}{
		{name: "unsupported country", body: `{"country":"CA"}`},
		{name: "count too high", body: `{"country":"US","count":101}`},
		{name: "invalid gender", body: `{"country":"US","filters":{"gender":"x"}}`},
		{name: "invalid age bounds", body: `{"country":"US","filters":{"age_min":50,"age_max":40}}`},
		{name: "invalid state", body: `{"country":"US","filters":{"state":"ZZ"}}`},
		{name: "invalid include field", body: `{"country":"US","include_fields":["education.unknown"]}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/profiles/generate", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), `"code":"validation_error"`) {
				t.Fatalf("expected validation_error body, got %s", rec.Body.String())
			}
		})
	}
}

func TestOpenAPIAndDocs(t *testing.T) {
	handler := api.NewServer(generator.New(func() time.Time {
		return time.Date(2026, time.March, 20, 15, 4, 5, 0, time.UTC)
	}))

	openAPIReq := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	openAPIRec := httptest.NewRecorder()
	handler.ServeHTTP(openAPIRec, openAPIReq)

	if openAPIRec.Code != http.StatusOK {
		t.Fatalf("expected /openapi.json 200, got %d", openAPIRec.Code)
	}
	if !strings.Contains(openAPIRec.Body.String(), `/v1/profiles/generate`) {
		t.Fatalf("expected spec to contain generate path")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"middle_name"`) {
		t.Fatalf("expected spec to contain middle_name field")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"email"`) {
		t.Fatalf("expected spec to contain email field")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"company_name"`) {
		t.Fatalf("expected spec to contain company_name field")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"job_title"`) {
		t.Fatalf("expected spec to contain job_title field")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"include_fields"`) {
		t.Fatalf("expected spec to contain include_fields")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"education"`) {
		t.Fatalf("expected spec to contain education schema")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"sensitive"`) {
		t.Fatalf("expected spec to contain sensitive schema")
	}
	if !strings.Contains(openAPIRec.Body.String(), `"expiry_month"`) {
		t.Fatalf("expected spec to contain bank card fields")
	}

	docsReq := httptest.NewRequest(http.MethodGet, "/docs", nil)
	docsRec := httptest.NewRecorder()
	handler.ServeHTTP(docsRec, docsReq)

	if docsRec.Code != http.StatusOK {
		t.Fatalf("expected /docs 200, got %d", docsRec.Code)
	}
	if !strings.Contains(docsRec.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("expected docs HTML to include SwaggerUIBundle")
	}
}

func calcAge(birthDate time.Time, reference time.Time) int {
	age := reference.Year() - birthDate.Year()
	if reference.Month() < birthDate.Month() || (reference.Month() == birthDate.Month() && reference.Day() < birthDate.Day()) {
		age--
	}
	return age
}
