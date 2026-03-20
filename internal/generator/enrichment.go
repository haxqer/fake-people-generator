package generator

import (
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"
	"unicode"

	"fake-people-generator/internal/dataset"
)

func (g *Generator) selectDistinctFirstName(rng *mathrand.Rand, gender Gender, firstName string) string {
	for i := 0; i < 6; i++ {
		candidate := g.selectFirstName(rng, gender)
		if candidate != "" && candidate != firstName {
			return candidate
		}
	}
	return g.selectFirstName(rng, gender)
}

func buildFullName(firstName string, middleName string, lastName string) string {
	if middleName == "" {
		return firstName + " " + lastName
	}
	return firstName + " " + middleName + " " + lastName
}

func generateAddressLine2(rng *mathrand.Rand, data *dataset.Data) string {
	if rng.Intn(100) >= 28 {
		return ""
	}

	label := selectWeightedValue(rng, data.UnitDesignators)
	number := 1 + rng.Intn(999)
	if label == "#" {
		return fmt.Sprintf("%s %03d", label, number)
	}
	return fmt.Sprintf("%s %d", label, number)
}

func generateUsername(rng *mathrand.Rand, first string, last string, birthDate string) string {
	if first == "" {
		first = "user"
	}
	if last == "" {
		last = "profile"
	}

	pattern := rng.Intn(4)
	switch pattern {
	case 0:
		return fmt.Sprintf("%s.%s%02d", first, last, rng.Intn(100))
	case 1:
		return fmt.Sprintf("%s_%s", string(first[0]), last)
	case 2:
		yearSuffix := ""
		if len(birthDate) >= 4 {
			yearSuffix = birthDate[2:4]
		}
		return fmt.Sprintf("%s%s%s", first, last, yearSuffix)
	default:
		return fmt.Sprintf("%s%s%03d", first, string(last[0]), 100+rng.Intn(900))
	}
}

func generateEmail(rng *mathrand.Rand, username string, first string, last string, data *dataset.Data) string {
	localPart := username
	switch rng.Intn(3) {
	case 0:
		localPart = username
	case 1:
		if first != "" && last != "" {
			localPart = fmt.Sprintf("%s.%s", first, last)
		}
	case 2:
		if first != "" && last != "" {
			localPart = fmt.Sprintf("%s%s%02d", string(first[0]), last, rng.Intn(100))
		}
	}

	return fmt.Sprintf("%s@%s", localPart, selectWeightedValue(rng, data.EmailDomains))
}

func selectOccupation(rng *mathrand.Rand, age int, data *dataset.Data) string {
	switch {
	case age < 18:
		return "Student"
	case age >= 67 && rng.Intn(100) < 60:
		return "Retired"
	default:
		return selectWeightedValue(rng, data.Occupations)
	}
}

func (g *Generator) generateEmployment(rng *mathrand.Rand, age int, occupation string, firstName string, lastName string, firstHandle string, lastHandle string, state dataset.State, city dataset.City, selection fieldSelection, education educationRecord, generatedAt time.Time) Employment {
	switch {
	case age < 18:
		return Employment{Status: "student"}
	case occupation == "Retired":
		return Employment{Status: "retired"}
	default:
		companyName := generateCompanyName(rng, city.Name, state.Name, g.data)
		companyDomain := generateCompanyDomain(rng, companyName, g.data)
		jobTitle := generateJobTitle(rng, age, occupation)
		department := selectDepartment(rng, occupation, g.data)

		employment := Employment{
			Status:        "employed",
			CompanyName:   companyName,
			CompanyDomain: companyDomain,
			JobTitle:      jobTitle,
			Department:    department,
			WorkEmail:     generateWorkEmail(rng, firstHandle, lastHandle, companyDomain),
		}

		if !selection.wantsEmploymentEnhancements() {
			return employment
		}

		industry := selectIndustry(occupation)
		seniority := selectSeniority(jobTitle, age)
		workMode := selectWorkMode(rng, occupation)
		managerName := generateManagerName(rng, firstName, lastName, g.data)
		officeCity, officeState := selectOfficeLocation(rng, state, city)
		startDate := generateEmploymentStartDate(rng, generatedAt, age, seniority, education)

		if selection.EmploymentIndustry {
			employment.Industry = industry
		}
		if selection.EmploymentSeniority {
			employment.Seniority = seniority
		}
		if selection.EmploymentStartDate {
			employment.StartDate = startDate
		}
		if selection.EmploymentWorkMode {
			employment.WorkMode = workMode
		}
		if selection.EmploymentManager {
			employment.ManagerName = managerName
		}
		if selection.EmploymentOfficeCity {
			employment.OfficeCity = officeCity
		}
		if selection.EmploymentOfficeState {
			employment.OfficeState = officeState
		}

		return employment
	}
}

func generateCompanyName(rng *mathrand.Rand, cityName string, stateName string, data *dataset.Data) string {
	cityToken := strings.TrimSpace(cityName)
	stateToken := strings.TrimSpace(stateName)
	prefix := selectWeightedValue(rng, data.CompanyPrefixes)
	noun := selectWeightedValue(rng, data.CompanyNouns)
	suffix := selectWeightedValue(rng, data.CompanySuffixes)

	switch rng.Intn(4) {
	case 0:
		return fmt.Sprintf("%s %s %s", prefix, noun, suffix)
	case 1:
		if cityToken != "" {
			return fmt.Sprintf("%s %s %s", cityToken, noun, suffix)
		}
	case 2:
		if stateToken != "" {
			return fmt.Sprintf("%s %s %s", stateToken, noun, suffix)
		}
	}

	return fmt.Sprintf("%s %s", prefix, noun)
}

func generateCompanyDomain(rng *mathrand.Rand, companyName string, data *dataset.Data) string {
	parts := strings.Fields(companyName)
	if len(parts) == 0 {
		return "example.com"
	}

	var builder strings.Builder
	for _, part := range parts {
		normalized := normalizeHandlePart(part)
		switch normalized {
		case "", "inc", "llc", "co", "corporation", "group", "partners", "brands":
			continue
		}
		builder.WriteString(normalized)
		if builder.Len() >= 18 {
			break
		}
	}

	domain := builder.String()
	if domain == "" {
		domain = "company"
	}

	return domain + selectWeightedValue(rng, data.DomainSuffixes)
}

func generateJobTitle(rng *mathrand.Rand, age int, occupation string) string {
	switch occupation {
	case "Software Engineer":
		return pickFromSet(rng, senioritySet(age, []string{"Junior Software Engineer", "Software Engineer", "Senior Software Engineer", "Lead Software Engineer"}))
	case "Project Manager":
		return pickFromSet(rng, senioritySet(age, []string{"Assistant Project Manager", "Project Manager", "Senior Project Manager", "Program Manager"}))
	case "Sales Associate":
		return pickFromSet(rng, senioritySet(age, []string{"Sales Associate", "Account Executive", "Senior Sales Associate", "Sales Manager"}))
	case "Registered Nurse":
		return pickFromSet(rng, senioritySet(age, []string{"Staff Nurse", "Registered Nurse", "Charge Nurse", "Clinical Nurse Lead"}))
	case "Teacher":
		return pickFromSet(rng, senioritySet(age, []string{"Teacher", "Lead Teacher", "Instructional Coach", "Academic Program Lead"}))
	case "Marketing Specialist":
		return pickFromSet(rng, senioritySet(age, []string{"Marketing Coordinator", "Marketing Specialist", "Senior Marketing Specialist", "Marketing Manager"}))
	case "Financial Analyst":
		return pickFromSet(rng, senioritySet(age, []string{"Junior Financial Analyst", "Financial Analyst", "Senior Financial Analyst", "Finance Manager"}))
	case "Customer Success Manager":
		return pickFromSet(rng, senioritySet(age, []string{"Customer Success Associate", "Customer Success Manager", "Senior Customer Success Manager", "Director of Customer Success"}))
	case "Graphic Designer":
		return pickFromSet(rng, senioritySet(age, []string{"Graphic Designer", "Senior Graphic Designer", "Art Director", "Creative Lead"}))
	case "Operations Coordinator":
		return pickFromSet(rng, senioritySet(age, []string{"Operations Assistant", "Operations Coordinator", "Operations Manager", "Senior Operations Manager"}))
	case "Accountant":
		return pickFromSet(rng, senioritySet(age, []string{"Staff Accountant", "Accountant", "Senior Accountant", "Accounting Manager"}))
	case "Warehouse Supervisor":
		return pickFromSet(rng, senioritySet(age, []string{"Warehouse Lead", "Warehouse Supervisor", "Distribution Manager", "Logistics Manager"}))
	case "Electrician":
		return pickFromSet(rng, senioritySet(age, []string{"Apprentice Electrician", "Electrician", "Senior Electrician", "Master Electrician"}))
	case "Pharmacist":
		return pickFromSet(rng, senioritySet(age, []string{"Pharmacist", "Clinical Pharmacist", "Lead Pharmacist", "Pharmacy Manager"}))
	case "Real Estate Agent":
		return pickFromSet(rng, senioritySet(age, []string{"Real Estate Agent", "Senior Real Estate Agent", "Broker Associate", "Managing Broker"}))
	case "HR Generalist":
		return pickFromSet(rng, senioritySet(age, []string{"HR Coordinator", "HR Generalist", "Senior HR Generalist", "HR Manager"}))
	case "Office Administrator":
		return pickFromSet(rng, senioritySet(age, []string{"Administrative Assistant", "Office Administrator", "Executive Assistant", "Office Manager"}))
	case "Restaurant Manager":
		return pickFromSet(rng, senioritySet(age, []string{"Shift Lead", "Restaurant Manager", "General Manager", "Operations Director"}))
	case "Delivery Driver":
		return pickFromSet(rng, senioritySet(age, []string{"Delivery Associate", "Delivery Driver", "Route Supervisor", "Fleet Coordinator"}))
	case "Business Analyst":
		return pickFromSet(rng, senioritySet(age, []string{"Junior Business Analyst", "Business Analyst", "Senior Business Analyst", "Strategy Manager"}))
	default:
		return occupation
	}
}

func senioritySet(age int, titles []string) []string {
	switch {
	case len(titles) == 0:
		return nil
	case age < 24 && len(titles) >= 2:
		return titles[:2]
	case age < 35 && len(titles) >= 3:
		return titles[:3]
	default:
		return titles
	}
}

func selectDepartment(rng *mathrand.Rand, occupation string, data *dataset.Data) string {
	switch occupation {
	case "Software Engineer":
		return "Engineering"
	case "Project Manager", "Operations Coordinator", "Warehouse Supervisor", "Restaurant Manager", "Delivery Driver":
		return "Operations"
	case "Electrician":
		return "Facilities"
	case "Sales Associate", "Real Estate Agent":
		return "Sales"
	case "Marketing Specialist", "Graphic Designer":
		return pickFromSet(rng, []string{"Marketing", "Design"})
	case "Financial Analyst", "Accountant":
		return "Finance"
	case "Registered Nurse", "Pharmacist":
		return "Clinical"
	case "Teacher", "Office Administrator":
		return "Administration"
	case "HR Generalist":
		return "Human Resources"
	case "Business Analyst":
		return pickFromSet(rng, []string{"Product", "Operations", "Finance"})
	default:
		return selectWeightedValue(rng, data.Departments)
	}
}

func selectIndustry(occupation string) string {
	switch occupation {
	case "Software Engineer", "Customer Success Manager", "Business Analyst":
		return "Technology"
	case "Project Manager", "HR Generalist", "Office Administrator":
		return "Professional Services"
	case "Sales Associate":
		return "Retail"
	case "Registered Nurse", "Pharmacist":
		return "Healthcare"
	case "Teacher":
		return "Education"
	case "Marketing Specialist", "Graphic Designer":
		return "Media and Marketing"
	case "Financial Analyst", "Accountant":
		return "Financial Services"
	case "Operations Coordinator", "Warehouse Supervisor", "Delivery Driver":
		return "Transportation and Logistics"
	case "Electrician":
		return "Construction and Facilities"
	case "Real Estate Agent":
		return "Real Estate"
	case "Restaurant Manager":
		return "Hospitality"
	default:
		return "Business Services"
	}
}

func selectSeniority(jobTitle string, age int) string {
	switch {
	case strings.Contains(jobTitle, "Director"), strings.Contains(jobTitle, "Lead"), strings.Contains(jobTitle, "Manager"):
		return "lead"
	case strings.Contains(jobTitle, "Senior"), strings.Contains(jobTitle, "Master"), strings.Contains(jobTitle, "Charge"):
		return "senior"
	case strings.Contains(jobTitle, "Junior"), strings.Contains(jobTitle, "Assistant"), strings.Contains(jobTitle, "Associate"), strings.Contains(jobTitle, "Apprentice"):
		return "entry"
	case age < 26:
		return "entry"
	case age < 36:
		return "mid"
	case age < 50:
		return "senior"
	default:
		return "lead"
	}
}

func selectWorkMode(rng *mathrand.Rand, occupation string) string {
	switch occupation {
	case "Registered Nurse", "Pharmacist", "Electrician", "Restaurant Manager", "Delivery Driver", "Warehouse Supervisor":
		return "onsite"
	case "Software Engineer", "Graphic Designer", "Marketing Specialist", "Financial Analyst", "Business Analyst":
		return pickFromSet(rng, []string{"remote", "hybrid", "hybrid"})
	case "Real Estate Agent", "Project Manager", "Customer Success Manager", "HR Generalist":
		return pickFromSet(rng, []string{"hybrid", "hybrid", "onsite", "remote"})
	default:
		return pickFromSet(rng, []string{"onsite", "hybrid", "hybrid"})
	}
}

func generateManagerName(rng *mathrand.Rand, firstName string, lastName string, data *dataset.Data) string {
	for i := 0; i < 6; i++ {
		var managerFirst string
		if rng.Intn(2) == 0 {
			managerFirst = selectWeightedValue(rng, data.MaleFirstNames)
		} else {
			managerFirst = selectWeightedValue(rng, data.FemaleFirstName)
		}
		managerLast := selectWeightedValue(rng, data.LastNames)
		if managerFirst != firstName || managerLast != lastName {
			return managerFirst + " " + managerLast
		}
	}
	return "Alex Taylor"
}

func selectOfficeLocation(rng *mathrand.Rand, state dataset.State, city dataset.City) (string, string) {
	if len(state.Cities) == 0 || rng.Intn(100) < 70 {
		return city.Name, state.Name
	}

	officeCity, err := selectCity(rng, state.Cities)
	if err != nil {
		return city.Name, state.Name
	}

	return officeCity.Name, state.Name
}

func generateEmploymentStartDate(rng *mathrand.Rand, generatedAt time.Time, age int, seniority string, education educationRecord) string {
	birthYearApprox := generatedAt.Year() - age
	minStartYear := birthYearApprox + 16
	if education.GraduationYear > 0 {
		if education.StudentStatus == "enrolled" {
			minStartYear = max(minStartYear, education.GraduationYear-2)
		} else {
			minStartYear = max(minStartYear, education.GraduationYear)
		}
	}

	switch seniority {
	case "entry":
		minStartYear = max(minStartYear, birthYearApprox+18)
	case "mid":
		minStartYear = max(minStartYear, birthYearApprox+21)
	case "senior":
		minStartYear = max(minStartYear, birthYearApprox+24)
	case "lead":
		minStartYear = max(minStartYear, birthYearApprox+26)
	}

	maxStartYear := generatedAt.Year()
	if seniority == "entry" {
		maxStartYear = generatedAt.Year()
	} else if seniority == "mid" {
		maxStartYear = generatedAt.Year() - 1
	} else {
		maxStartYear = generatedAt.Year() - 2
	}

	if maxStartYear < minStartYear {
		maxStartYear = minStartYear
	}

	year := minStartYear
	if maxStartYear > minStartYear {
		year = minStartYear + rng.Intn(maxStartYear-minStartYear+1)
	}

	month := time.Month(1 + rng.Intn(12))
	day := 1 + rng.Intn(28)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

func generateWorkEmail(rng *mathrand.Rand, first string, last string, companyDomain string) string {
	if first == "" || last == "" {
		return ""
	}

	switch rng.Intn(4) {
	case 0:
		return fmt.Sprintf("%s.%s@%s", first, last, companyDomain)
	case 1:
		return fmt.Sprintf("%s%s@%s", string(first[0]), last, companyDomain)
	case 2:
		return fmt.Sprintf("%s@%s", first, companyDomain)
	default:
		return fmt.Sprintf("%s.%s%02d@%s", first, last, rng.Intn(100), companyDomain)
	}
}

func selectMaritalStatus(rng *mathrand.Rand, age int) string {
	switch {
	case age < 18:
		return "single"
	case age <= 24:
		return pickFromSet(rng, []string{"single", "single", "single", "single", "married"})
	case age <= 34:
		return pickFromSet(rng, []string{"single", "single", "married", "married", "divorced"})
	case age <= 54:
		return pickFromSet(rng, []string{"married", "married", "married", "single", "divorced"})
	case age <= 69:
		return pickFromSet(rng, []string{"married", "married", "divorced", "widowed", "single"})
	default:
		return pickFromSet(rng, []string{"married", "widowed", "divorced", "single"})
	}
}

func selectTimeZone(rng *mathrand.Rand, stateCode string) string {
	switch stateCode {
	case "AK":
		return "America/Anchorage"
	case "HI":
		return "Pacific/Honolulu"
	case "AZ":
		return "America/Phoenix"
	case "CA", "NV", "OR", "WA":
		return "America/Los_Angeles"
	case "CO", "ID", "MT", "NM", "UT", "WY":
		return "America/Denver"
	case "ND", "SD", "NE", "KS", "MN", "IA", "MO", "WI", "IL", "TN", "MS", "AL", "OK", "TX", "LA", "AR":
		return "America/Chicago"
	case "FL", "IN", "KY", "MI":
		if rng.Intn(2) == 0 {
			return "America/New_York"
		}
		return "America/Chicago"
	default:
		return "America/New_York"
	}
}

func pickFromSet(rng *mathrand.Rand, values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[rng.Intn(len(values))]
}

func normalizeHandlePart(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))

	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
