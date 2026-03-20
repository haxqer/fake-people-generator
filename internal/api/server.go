package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fake-people-generator/internal/generator"
	"fake-people-generator/internal/openapi"
	"fake-people-generator/internal/web"
)

type generateProfilesRequest struct {
	Country       string                 `json:"country"`
	Count         *int                   `json:"count,omitempty"`
	Filters       generateFiltersRequest `json:"filters,omitempty"`
	IncludeFields []string               `json:"include_fields,omitempty"`
}

type generateFiltersRequest struct {
	Gender string `json:"gender,omitempty"`
	AgeMin *int   `json:"age_min,omitempty"`
	AgeMax *int   `json:"age_max,omitempty"`
	State  string `json:"state,omitempty"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type metaResponse struct {
	Country     string `json:"country"`
	Count       int    `json:"count"`
	GeneratedAt string `json:"generated_at"`
}

type generateProfilesResponse struct {
	Meta     metaResponse        `json:"meta"`
	Profiles []generator.Profile `json:"profiles"`
}

func NewServer(gen *generator.Generator) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", web.IndexHandler())
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", web.AssetsHandler()))
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("POST /v1/profiles/generate", generateProfiles(gen))
	mux.HandleFunc("GET /openapi.json", handleOpenAPI)
	mux.HandleFunc("GET /docs", handleDocs)
	return mux
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func generateProfiles(gen *generator.Generator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := validateJSONRequest(r); err != nil {
			writeAPIError(w, err)
			return
		}

		var req generateProfilesRequest
		if err := decodeJSON(r.Body, &req); err != nil {
			writeAPIError(w, &generator.ValidationError{Message: err.Error()})
			return
		}

		count := 1
		if req.Count != nil {
			count = *req.Count
		}

		result, err := gen.Generate(generator.GenerateInput{
			Country:       req.Country,
			Count:         count,
			IncludeFields: req.IncludeFields,
			Filters: generator.Filters{
				Gender:    generator.Gender(strings.ToLower(strings.TrimSpace(req.Filters.Gender))),
				AgeMin:    req.Filters.AgeMin,
				AgeMax:    req.Filters.AgeMax,
				StateCode: req.Filters.State,
			},
		})
		if err != nil {
			writeAPIError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, generateProfilesResponse{
			Meta: metaResponse{
				Country:     result.Country,
				Count:       result.Count,
				GeneratedAt: result.GeneratedAt.Format(time.RFC3339),
			},
			Profiles: result.Profiles,
		})
	}
}

func handleOpenAPI(w http.ResponseWriter, _ *http.Request) {
	spec, err := openapi.JSON()
	if err != nil {
		writeAPIError(w, fmt.Errorf("generate openapi spec: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(spec)
}

func handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, openapi.DocsHTML())
}

func validateJSONRequest(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return nil
	}
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		return &generator.ValidationError{Message: "Content-Type must be application/json"}
	}
	return nil
}

func decodeJSON(body io.ReadCloser, dst any) error {
	defer body.Close()

	decoder := json.NewDecoder(io.LimitReader(body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("request body is required")
		}
		return fmt.Errorf("invalid JSON body: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("request body must contain a single JSON object")
	}

	return nil
}

func writeAPIError(w http.ResponseWriter, err error) {
	var validationErr *generator.ValidationError
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "internal server error"

	if errors.As(err, &validationErr) {
		status = http.StatusBadRequest
		code = "validation_error"
		message = validationErr.Message
	}

	writeJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	_ = encoder.Encode(payload)
}
