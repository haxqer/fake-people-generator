package generator

import (
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"

	"fake-people-generator/internal/dataset"
)

type Education struct {
	Level          string `json:"level,omitempty"`
	DegreeType     string `json:"degree_type,omitempty"`
	SchoolName     string `json:"school_name,omitempty"`
	Major          string `json:"major,omitempty"`
	GraduationYear int    `json:"graduation_year,omitempty"`
	StudentStatus  string `json:"student_status,omitempty"`
	SchoolCity     string `json:"school_city,omitempty"`
	SchoolState    string `json:"school_state,omitempty"`
	SchoolEmail    string `json:"school_email,omitempty"`
}

type educationRecord struct {
	Level          string
	DegreeType     string
	SchoolName     string
	Major          string
	GraduationYear int
	StudentStatus  string
	SchoolCity     string
	SchoolState    string
	SchoolEmail    string
}

var collegeMajors = []string{
	"Computer Science", "Business Administration", "Marketing", "Accounting", "Finance",
	"Nursing", "Communications", "Graphic Design", "Psychology", "Economics",
	"Data Analytics", "Information Systems", "Education", "Biology", "Political Science",
}

var tradeMajors = []string{
	"Electrical Technology", "Industrial Maintenance", "Logistics Operations", "Office Administration",
	"Hospitality Management", "Applied Business", "Network Support", "Medical Assisting",
}

var universitySuffixes = []string{
	"University", "State University", "College", "Institute of Technology",
}

func generateEducationRecord(rng *mathrand.Rand, generatedAt time.Time, age int, occupation string, firstName string, lastName string, state dataset.State, homeCity dataset.City) educationRecord {
	schoolCity := selectSchoolCity(rng, state, homeCity)
	schoolState := state.Name

	switch {
	case age < 5:
		return educationRecord{
			Level:         "none",
			StudentStatus: "not_enrolled",
		}
	case age <= 10:
		return educationRecord{
			Level:         "elementary",
			SchoolName:    fmt.Sprintf("%s Elementary School", schoolCity),
			StudentStatus: "enrolled",
			SchoolCity:    schoolCity,
			SchoolState:   schoolState,
		}
	case age <= 13:
		return educationRecord{
			Level:         "middle_school",
			SchoolName:    fmt.Sprintf("%s Middle School", schoolCity),
			StudentStatus: "enrolled",
			SchoolCity:    schoolCity,
			SchoolState:   schoolState,
		}
	case age <= 17:
		return educationRecord{
			Level:          "high_school",
			DegreeType:     "High School Diploma",
			SchoolName:     fmt.Sprintf("%s High School", schoolCity),
			GraduationYear: generatedAt.Year() + max(0, 18-age),
			StudentStatus:  "enrolled",
			SchoolCity:     schoolCity,
			SchoolState:    schoolState,
		}
	default:
		record := educationRecord{
			SchoolCity:  schoolCity,
			SchoolState: schoolState,
		}

		assignAdultEducation(rng, generatedAt, age, occupation, &record)
		record.SchoolName = generateSchoolName(rng, record.Level, schoolCity, schoolState)
		record.SchoolEmail = generateSchoolEmail(firstName, lastName, record)
		return record
	}
}

func assignAdultEducation(rng *mathrand.Rand, generatedAt time.Time, age int, occupation string, record *educationRecord) {
	birthYearApprox := generatedAt.Year() - age

	switch occupation {
	case "Pharmacist":
		record.Level = "doctorate"
		record.DegreeType = "Doctor of Pharmacy"
		record.Major = "Pharmacy"
		record.StudentStatus = "alumni"
		record.GraduationYear = birthYearApprox + 26 + rng.Intn(3)
	case "Registered Nurse":
		if rng.Intn(100) < 65 {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Science"
			record.Major = "Nursing"
			record.GraduationYear = birthYearApprox + 22 + rng.Intn(2)
		} else {
			record.Level = "associate"
			record.DegreeType = "Associate of Science"
			record.Major = "Nursing"
			record.GraduationYear = birthYearApprox + 20 + rng.Intn(2)
		}
		record.StudentStatus = "alumni"
	case "Teacher":
		if age >= 30 && rng.Intn(100) < 45 {
			record.Level = "master"
			record.DegreeType = "Master of Education"
			record.Major = "Education"
			record.GraduationYear = birthYearApprox + 24 + rng.Intn(2)
		} else {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Arts"
			record.Major = "Education"
			record.GraduationYear = birthYearApprox + 22 + rng.Intn(2)
		}
		record.StudentStatus = "alumni"
	case "Software Engineer", "Financial Analyst", "Business Analyst", "Project Manager", "HR Generalist", "Marketing Specialist", "Graphic Designer", "Customer Success Manager", "Accountant":
		if age <= 24 && rng.Intn(100) < 40 {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Science"
			record.Major = pickFromSet(rng, []string{"Computer Science", "Business Administration", "Marketing", "Accounting", "Finance", "Data Analytics", "Information Systems", "Graphic Design"})
			record.StudentStatus = "enrolled"
			record.GraduationYear = generatedAt.Year() + rng.Intn(3)
		} else if age >= 32 && rng.Intn(100) < 20 {
			record.Level = "master"
			record.DegreeType = "Master of Science"
			record.Major = pickFromSet(rng, []string{"Computer Science", "Business Administration", "Finance", "Data Analytics", "Marketing"})
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 24 + rng.Intn(2)
		} else {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Science"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 22 + rng.Intn(2)
		}
	case "Electrician", "Warehouse Supervisor", "Delivery Driver", "Office Administrator", "Restaurant Manager":
		if age <= 22 && rng.Intn(100) < 35 {
			record.Level = "professional_certificate"
			record.DegreeType = "Technical Certificate"
			record.Major = pickFromSet(rng, tradeMajors)
			record.StudentStatus = "enrolled"
			record.GraduationYear = generatedAt.Year() + rng.Intn(2)
		} else if rng.Intn(100) < 55 {
			record.Level = "professional_certificate"
			record.DegreeType = "Technical Certificate"
			record.Major = pickFromSet(rng, tradeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 19 + rng.Intn(2)
		} else {
			record.Level = "high_school"
			record.DegreeType = "High School Diploma"
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 18
		}
	default:
		assignGenericAdultEducation(rng, generatedAt, age, birthYearApprox, record)
	}

	if record.GraduationYear > generatedAt.Year()+3 {
		record.GraduationYear = generatedAt.Year() + 3
	}
	if record.StudentStatus != "enrolled" && record.GraduationYear > generatedAt.Year() {
		record.GraduationYear = generatedAt.Year()
	}
}

func assignGenericAdultEducation(rng *mathrand.Rand, generatedAt time.Time, age int, birthYearApprox int, record *educationRecord) {
	switch {
	case age <= 20:
		if rng.Intn(100) < 55 {
			record.Level = "high_school"
			record.DegreeType = "High School Diploma"
			record.StudentStatus = "enrolled"
			record.GraduationYear = birthYearApprox + 18
		} else {
			record.Level = "associate"
			record.DegreeType = "Associate of Arts"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "enrolled"
			record.GraduationYear = generatedAt.Year() + rng.Intn(3)
		}
	case age <= 24:
		if rng.Intn(100) < 50 {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Arts"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "enrolled"
			record.GraduationYear = generatedAt.Year() + rng.Intn(2)
		} else {
			record.Level = "associate"
			record.DegreeType = "Associate of Arts"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 20 + rng.Intn(2)
		}
	case age <= 34:
		if rng.Intn(100) < 65 {
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Science"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 22 + rng.Intn(2)
		} else {
			record.Level = "associate"
			record.DegreeType = "Associate of Science"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 20 + rng.Intn(2)
		}
	default:
		pattern := rng.Intn(100)
		switch {
		case pattern < 50:
			record.Level = "bachelor"
			record.DegreeType = "Bachelor of Science"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 22 + rng.Intn(2)
		case pattern < 72:
			record.Level = "associate"
			record.DegreeType = "Associate of Arts"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 20 + rng.Intn(2)
		case pattern < 84:
			record.Level = "master"
			record.DegreeType = "Master of Science"
			record.Major = pickFromSet(rng, collegeMajors)
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 24 + rng.Intn(2)
		default:
			record.Level = "high_school"
			record.DegreeType = "High School Diploma"
			record.StudentStatus = "alumni"
			record.GraduationYear = birthYearApprox + 18
		}
	}
}

func selectSchoolCity(rng *mathrand.Rand, state dataset.State, homeCity dataset.City) string {
	if len(state.Cities) == 0 || rng.Intn(100) < 70 {
		return homeCity.Name
	}

	city, err := selectCity(rng, state.Cities)
	if err != nil {
		return homeCity.Name
	}
	return city.Name
}

func generateSchoolName(rng *mathrand.Rand, level string, schoolCity string, schoolState string) string {
	switch level {
	case "associate", "professional_certificate":
		return fmt.Sprintf("%s Community College", schoolCity)
	case "bachelor", "master", "doctorate":
		if rng.Intn(100) < 55 {
			return fmt.Sprintf("%s %s", schoolCity, pickFromSet(rng, universitySuffixes))
		}
		return fmt.Sprintf("%s %s", strings.TrimSpace(schoolState), pickFromSet(rng, []string{"University", "College"}))
	case "high_school":
		return fmt.Sprintf("%s High School", schoolCity)
	default:
		return ""
	}
}

func generateSchoolEmail(firstName string, lastName string, record educationRecord) string {
	if record.StudentStatus != "enrolled" || record.SchoolName == "" {
		return ""
	}

	if record.Level == "high_school" || record.Level == "middle_school" || record.Level == "elementary" || record.Level == "none" {
		return ""
	}

	localPart := normalizeHandlePart(firstName) + "." + normalizeHandlePart(lastName)
	localPart = strings.Trim(localPart, ".")
	if localPart == "" {
		localPart = "student"
	}

	return fmt.Sprintf("%s@%s", localPart, generateSchoolDomain(record.SchoolName))
}

func generateSchoolDomain(schoolName string) string {
	parts := strings.Fields(schoolName)
	var builder strings.Builder

	for _, part := range parts {
		normalized := normalizeHandlePart(part)
		switch normalized {
		case "", "university", "college", "community", "institute", "technology", "state":
			continue
		}
		builder.WriteString(normalized)
		if builder.Len() >= 16 {
			break
		}
	}

	domain := builder.String()
	if domain == "" {
		domain = "campus"
	}
	return domain + ".edu"
}

func buildEducation(selection fieldSelection, record educationRecord) *Education {
	if !selection.wantsEducation() {
		return nil
	}

	education := &Education{}
	if selection.EducationLevel {
		education.Level = record.Level
	}
	if selection.EducationDegreeType {
		education.DegreeType = record.DegreeType
	}
	if selection.EducationSchoolName {
		education.SchoolName = record.SchoolName
	}
	if selection.EducationMajor {
		education.Major = record.Major
	}
	if selection.EducationGraduationYear {
		education.GraduationYear = record.GraduationYear
	}
	if selection.EducationStudentStatus {
		education.StudentStatus = record.StudentStatus
	}
	if selection.EducationSchoolCity {
		education.SchoolCity = record.SchoolCity
	}
	if selection.EducationSchoolState {
		education.SchoolState = record.SchoolState
	}
	if selection.EducationSchoolEmail {
		education.SchoolEmail = record.SchoolEmail
	}

	return education
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
