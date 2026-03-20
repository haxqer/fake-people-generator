package openapi

import "encoding/json"

func JSON() ([]byte, error) {
	spec := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Fake People Generator API",
			"version":     "1.0.0",
			"description": "Generate high-fidelity synthetic US personal profiles.",
		},
		"servers": []map[string]any{
			{"url": "/"},
		},
		"paths": map[string]any{
			"/healthz": map[string]any{
				"get": map[string]any{
					"summary": "Health check",
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Service is healthy",
						},
					},
				},
			},
			"/v1/profiles/generate": map[string]any{
				"post": map[string]any{
					"summary":     "Generate synthetic US profiles",
					"description": "Generate one or more synthetic but internally consistent US personal profiles.",
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/GenerateProfilesRequest",
								},
								"example": map[string]any{
									"country": "US",
									"count":   2,
									"filters": map[string]any{
										"gender":  "female",
										"age_min": 25,
										"age_max": 35,
										"state":   "CA",
									},
									"include_fields": []string{
										"education.level",
										"employment.industry",
										"sensitive.ssn",
									},
								},
							},
						},
					},
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Generated profiles",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"$ref": "#/components/schemas/GenerateProfilesResponse",
									},
								},
							},
						},
						"400": map[string]any{
							"description": "Validation error",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{
										"$ref": "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
		},
		"components": map[string]any{
			"schemas": map[string]any{
				"GenerateProfilesRequest": map[string]any{
					"type":     "object",
					"required": []string{"country"},
					"properties": map[string]any{
						"country": map[string]any{
							"type":    "string",
							"example": "US",
						},
						"count": map[string]any{
							"type":        "integer",
							"minimum":     1,
							"maximum":     100,
							"default":     1,
							"description": "Number of profiles to generate.",
						},
						"filters": map[string]any{
							"$ref": "#/components/schemas/GenerateFilters",
						},
						"include_fields": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
							"description": "Optional field-level expansions to include in the response.",
						},
					},
				},
				"GenerateFilters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"gender": map[string]any{
							"type": "string",
							"enum": []string{"male", "female", "random"},
						},
						"age_min": map[string]any{
							"type":    "integer",
							"minimum": 0,
							"maximum": 100,
						},
						"age_max": map[string]any{
							"type":    "integer",
							"minimum": 0,
							"maximum": 100,
						},
						"state": map[string]any{
							"type":        "string",
							"description": "Two-letter US state code.",
							"example":     "CA",
						},
					},
				},
				"GenerateProfilesResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"meta": map[string]any{
							"$ref": "#/components/schemas/ResponseMeta",
						},
						"profiles": map[string]any{
							"type": "array",
							"items": map[string]any{
								"$ref": "#/components/schemas/Profile",
							},
						},
					},
				},
				"ResponseMeta": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"country": map[string]any{"type": "string"},
						"count": map[string]any{
							"type": "integer",
						},
						"generated_at": map[string]any{
							"type":   "string",
							"format": "date-time",
						},
					},
				},
				"Profile": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":             map[string]any{"type": "string", "format": "uuid"},
						"country":        map[string]any{"type": "string"},
						"gender":         map[string]any{"type": "string"},
						"first_name":     map[string]any{"type": "string"},
						"middle_name":    map[string]any{"type": "string"},
						"last_name":      map[string]any{"type": "string"},
						"full_name":      map[string]any{"type": "string"},
						"age":            map[string]any{"type": "integer"},
						"birth_date":     map[string]any{"type": "string", "format": "date"},
						"marital_status": map[string]any{"type": "string"},
						"occupation":     map[string]any{"type": "string"},
						"username":       map[string]any{"type": "string"},
						"email":          map[string]any{"type": "string", "format": "email"},
						"time_zone":      map[string]any{"type": "string"},
						"employment": map[string]any{
							"$ref": "#/components/schemas/Employment",
						},
						"education": map[string]any{
							"$ref": "#/components/schemas/Education",
						},
						"sensitive": map[string]any{
							"$ref": "#/components/schemas/Sensitive",
						},
						"phone": map[string]any{
							"$ref": "#/components/schemas/Phone",
						},
						"address": map[string]any{
							"$ref": "#/components/schemas/Address",
						},
					},
				},
				"Phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"e164":      map[string]any{"type": "string"},
						"national":  map[string]any{"type": "string"},
						"area_code": map[string]any{"type": "string"},
					},
				},
				"Employment": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"status":         map[string]any{"type": "string"},
						"company_name":   map[string]any{"type": "string"},
						"company_domain": map[string]any{"type": "string"},
						"job_title":      map[string]any{"type": "string"},
						"department":     map[string]any{"type": "string"},
						"work_email":     map[string]any{"type": "string", "format": "email"},
						"industry":       map[string]any{"type": "string"},
						"seniority":      map[string]any{"type": "string"},
						"start_date":     map[string]any{"type": "string", "format": "date"},
						"work_mode":      map[string]any{"type": "string"},
						"manager_name":   map[string]any{"type": "string"},
						"office_city":    map[string]any{"type": "string"},
						"office_state":   map[string]any{"type": "string"},
					},
				},
				"Education": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"level":           map[string]any{"type": "string"},
						"degree_type":     map[string]any{"type": "string"},
						"school_name":     map[string]any{"type": "string"},
						"major":           map[string]any{"type": "string"},
						"graduation_year": map[string]any{"type": "integer"},
						"student_status":  map[string]any{"type": "string"},
						"school_city":     map[string]any{"type": "string"},
						"school_state":    map[string]any{"type": "string"},
						"school_email":    map[string]any{"type": "string", "format": "email"},
					},
				},
				"Sensitive": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"ssn": map[string]any{"type": "string"},
						"driver_license": map[string]any{
							"$ref": "#/components/schemas/DriverLicense",
						},
						"passport": map[string]any{
							"$ref": "#/components/schemas/Passport",
						},
						"bank_card": map[string]any{
							"$ref": "#/components/schemas/BankCard",
						},
					},
				},
				"DriverLicense": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"number":     map[string]any{"type": "string"},
						"state_code": map[string]any{"type": "string"},
					},
				},
				"Passport": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"number":  map[string]any{"type": "string"},
						"country": map[string]any{"type": "string"},
					},
				},
				"BankCard": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"brand":        map[string]any{"type": "string"},
						"number":       map[string]any{"type": "string"},
						"expiry_month": map[string]any{"type": "integer"},
						"expiry_year":  map[string]any{"type": "integer"},
						"cvv":          map[string]any{"type": "string"},
					},
				},
				"Address": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"street_line1": map[string]any{"type": "string"},
						"street_line2": map[string]any{"type": "string"},
						"city":         map[string]any{"type": "string"},
						"state":        map[string]any{"type": "string"},
						"state_code":   map[string]any{"type": "string"},
						"postal_code":  map[string]any{"type": "string"},
					},
				},
				"ErrorResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"error": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"code":    map[string]any{"type": "string"},
								"message": map[string]any{"type": "string"},
							},
						},
					},
				},
			},
		},
	}

	return json.MarshalIndent(spec, "", "  ")
}

func DocsHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Fake People Generator API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`
}
