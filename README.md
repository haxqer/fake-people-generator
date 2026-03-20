[English](README.md) | [简体中文](README.zh-CN.md)

# Fake People Generator

A synthetic profile generation API written in Go. The current version supports US-only data generation and aims to produce "highly realistic but explicitly synthetic" personal information for testing, demos, integration, load testing, and fake data seeding.

## Features

- Supports a `REST JSON` API
- Includes a utility frontend page that automatically generates a sample profile on the home page
- Supports single and batch generation, up to `100` profiles per request
- Supports filtering by `gender`, `age_min`, `age_max`, and `state`
- Supports requesting enhanced fields and sensitive fields precisely through `include_fields`
- Keeps generated name, age, birth date, state/city/ZIP code, area code, and phone number consistent with each other
- Returns a more complete synthetic profile schema:
  - Name: `first_name`, `middle_name`, `last_name`, `full_name`
  - Basic info: `gender`, `age`, `birth_date`, `marital_status`, `occupation`
  - Employment info: `employment.status`, `employment.company_name`, `employment.job_title`, `employment.company_domain`, `employment.work_email`
  - Education info: `education.level`, `education.degree_type`, `education.school_name`, `education.major`, `education.graduation_year`
  - Sensitive info: `sensitive.ssn`, `sensitive.driver_license.*`, `sensitive.passport.*`, `sensitive.bank_card.*`
  - Contact info: `phone`, `username`, `email`
  - Address info: `street_line1`, `street_line2`, `city`, `state`, `state_code`, `postal_code`
  - Region-derived field: `time_zone`
- Provides OpenAPI documentation and Swagger UI

## Notes

- Currently only supports `country = "US"`
- Returned data is synthetic and does not guarantee a match to a real person, address, or phone number
- By default only base fields are returned; education enhancements and sensitive fields are returned only when explicitly requested in `include_fields`

## Running

Requirements:

- Go `1.25+`

Start the service:

```bash
go run ./cmd/server
```

Default listen address:

```text
http://127.0.0.1:8080
```

Home page tool:

```text
http://127.0.0.1:8080/
```

You can also change the port with an environment variable:

```bash
ADDR=:9090 go run ./cmd/server
```

Or use the more container-friendly `PORT`:

```bash
PORT=9090 go run ./cmd/server
```

## Docker Deployment

Build the image:

```bash
docker build -t fake-people-generator:local .
```

Run the container:

```bash
docker run --rm -p 8080:8080 fake-people-generator:local
```

If you need a custom listen port inside the container:

```bash
docker run --rm -e PORT=9090 -p 9090:9090 fake-people-generator:local
```

The current `Dockerfile` uses a multi-stage build:

- Build stage: `golang:1.25-bookworm`
- Runtime stage: `distroless static-debian12:nonroot`
- Default container port: `8080`
- Supports multi-arch builds for `linux/amd64` and `linux/arm64`

## GitHub Actions: Build and Push to Docker Hub

The repository already includes this workflow:

```text
.github/workflows/docker.yml
```

Triggers:

- Push to `main`
- Push a version tag such as `v0.1.0`
- Manual `workflow_dispatch`

Configure the following GitHub repository secrets first:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

You can optionally configure this repository variable:

- `DOCKERHUB_REPOSITORY`

Notes:

- If `DOCKERHUB_REPOSITORY` is not set, the current GitHub repository name is used as the Docker Hub image name by default
- The final image name defaults to `docker.io/<DOCKERHUB_USERNAME>/fake-people-generator`
- The `main` branch pushes `latest`, `main`, and `sha-<commit>` tags
- `v*` tags push the matching version tag, for example `v0.1.0`
- The workflow runs `go test ./...` before building and pushing multi-arch images

Release example:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Pull example:

```bash
docker pull docker.io/<DOCKERHUB_USERNAME>/fake-people-generator:latest
```

## API

### Health Check

```bash
curl http://127.0.0.1:8080/healthz
```

### Generate Profiles

Request:

```bash
curl -X POST http://127.0.0.1:8080/v1/profiles/generate \
  -H 'Content-Type: application/json' \
  -d '{
    "country": "US",
    "count": 2,
    "filters": {
      "gender": "female",
      "age_min": 28,
      "age_max": 35,
      "state": "CA"
    },
    "include_fields": [
      "education.level",
      "employment.industry",
      "sensitive.ssn"
    ]
  }'
```

Example response:

```json
{
  "meta": {
    "country": "US",
    "count": 1,
    "generated_at": "2026-03-20T05:35:07Z"
  },
  "profiles": [
    {
      "id": "b88418b7-3fc4-4405-bfe3-9dbe558a75ea",
      "country": "US",
      "gender": "female",
      "first_name": "Emily",
      "middle_name": "Grace",
      "last_name": "Smith",
      "full_name": "Emily Grace Smith",
      "age": 29,
      "birth_date": "1996-04-20",
      "marital_status": "married",
      "occupation": "Software Engineer",
      "username": "emily.smith42",
      "email": "emily.smith@gmail.com",
      "time_zone": "America/Los_Angeles",
      "employment": {
        "status": "employed",
        "company_name": "Pacific Systems Inc.",
        "company_domain": "pacificsystems.com",
        "job_title": "Senior Software Engineer",
        "department": "Engineering",
        "work_email": "emily.smith@pacificsystems.com",
        "industry": "Technology",
        "seniority": "senior",
        "start_date": "2020-09-14",
        "work_mode": "hybrid",
        "manager_name": "Nancy Wilson",
        "office_city": "San Jose",
        "office_state": "California"
      },
      "education": {
        "level": "bachelor",
        "school_name": "California University",
        "major": "Computer Science",
        "graduation_year": 2018,
        "student_status": "alumni",
        "school_email": "emily.smith@california.edu"
      },
      "sensitive": {
        "ssn": "428-15-9012"
      },
      "phone": {
        "e164": "+14086827484",
        "national": "(408) 682-7484",
        "area_code": "408"
      },
      "address": {
        "street_line1": "9535 Highland Boulevard",
        "street_line2": "Apt 217",
        "city": "San Jose",
        "state": "California",
        "state_code": "CA",
        "postal_code": "95113"
      }
    }
  ]
}
```

### Request Parameters

Request body:

```json
{
  "country": "US",
  "count": 1,
  "filters": {
    "gender": "male|female|random",
    "age_min": 18,
    "age_max": 65,
    "state": "CA"
  },
  "include_fields": [
    "education.level",
    "employment.industry",
    "sensitive.ssn"
  ]
}
```

Rules:

- `country` is required and currently can only be `US`
- `count` is optional, defaults to `1`, and must be within `1..100`
- `filters.gender` is optional and supports `male`, `female`, and `random`
- `filters.age_min` and `filters.age_max` must be provided together
- `filters.state` is a two-letter US state code such as `CA`, `NY`, or `TX`
- `include_fields` is optional; undeclared enhanced fields are not returned

Common `include_fields` values:

- Education: `education.level`, `education.degree_type`, `education.school_name`, `education.major`, `education.graduation_year`, `education.student_status`, `education.school_city`, `education.school_state`, `education.school_email`
- Employment enhancements: `employment.industry`, `employment.seniority`, `employment.start_date`, `employment.work_mode`, `employment.manager_name`, `employment.office_city`, `employment.office_state`
- Sensitive: `sensitive.ssn`, `sensitive.driver_license.number`, `sensitive.driver_license.state_code`, `sensitive.passport.number`, `sensitive.passport.country`, `sensitive.bank_card.brand`, `sensitive.bank_card.number`, `sensitive.bank_card.expiry_month`, `sensitive.bank_card.expiry_year`, `sensitive.bank_card.cvv`

### Error Response

```json
{
  "error": {
    "code": "validation_error",
    "message": "count must be between 1 and 100"
  }
}
```

## Documentation

- Home page tool: `http://127.0.0.1:8080/`
- OpenAPI JSON: `http://127.0.0.1:8080/openapi.json`
- Swagger UI: `http://127.0.0.1:8080/docs`

Notes:

- `/docs` uses Swagger UI assets from a CDN, so loading the page requires network access

## Testing

Run all tests:

```bash
go test ./...
```
