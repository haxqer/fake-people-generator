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
		if profile.Age != 30 {
			t.Fatalf("expected age 30, got %d", profile.Age)
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
