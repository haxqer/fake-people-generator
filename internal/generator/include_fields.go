package generator

import "strings"

type fieldSelection struct {
	EducationLevel          bool
	EducationDegreeType     bool
	EducationSchoolName     bool
	EducationMajor          bool
	EducationGraduationYear bool
	EducationStudentStatus  bool
	EducationSchoolCity     bool
	EducationSchoolState    bool
	EducationSchoolEmail    bool

	EmploymentIndustry    bool
	EmploymentSeniority   bool
	EmploymentStartDate   bool
	EmploymentWorkMode    bool
	EmploymentManager     bool
	EmploymentOfficeCity  bool
	EmploymentOfficeState bool

	SensitiveSSN                bool
	SensitiveDriverLicenseNo    bool
	SensitiveDriverLicenseState bool
	SensitivePassportNumber     bool
	SensitivePassportCountry    bool
	SensitiveBankCardBrand      bool
	SensitiveBankCardNumber     bool
	SensitiveBankCardMonth      bool
	SensitiveBankCardYear       bool
	SensitiveBankCardCVV        bool
}

var supportedIncludeFields = map[string]func(*fieldSelection){
	"education.level":                     func(s *fieldSelection) { s.EducationLevel = true },
	"education.degree_type":               func(s *fieldSelection) { s.EducationDegreeType = true },
	"education.school_name":               func(s *fieldSelection) { s.EducationSchoolName = true },
	"education.major":                     func(s *fieldSelection) { s.EducationMajor = true },
	"education.graduation_year":           func(s *fieldSelection) { s.EducationGraduationYear = true },
	"education.student_status":            func(s *fieldSelection) { s.EducationStudentStatus = true },
	"education.school_city":               func(s *fieldSelection) { s.EducationSchoolCity = true },
	"education.school_state":              func(s *fieldSelection) { s.EducationSchoolState = true },
	"education.school_email":              func(s *fieldSelection) { s.EducationSchoolEmail = true },
	"employment.industry":                 func(s *fieldSelection) { s.EmploymentIndustry = true },
	"employment.seniority":                func(s *fieldSelection) { s.EmploymentSeniority = true },
	"employment.start_date":               func(s *fieldSelection) { s.EmploymentStartDate = true },
	"employment.work_mode":                func(s *fieldSelection) { s.EmploymentWorkMode = true },
	"employment.manager_name":             func(s *fieldSelection) { s.EmploymentManager = true },
	"employment.office_city":              func(s *fieldSelection) { s.EmploymentOfficeCity = true },
	"employment.office_state":             func(s *fieldSelection) { s.EmploymentOfficeState = true },
	"sensitive.ssn":                       func(s *fieldSelection) { s.SensitiveSSN = true },
	"sensitive.driver_license.number":     func(s *fieldSelection) { s.SensitiveDriverLicenseNo = true },
	"sensitive.driver_license.state_code": func(s *fieldSelection) { s.SensitiveDriverLicenseState = true },
	"sensitive.passport.number":           func(s *fieldSelection) { s.SensitivePassportNumber = true },
	"sensitive.passport.country":          func(s *fieldSelection) { s.SensitivePassportCountry = true },
	"sensitive.bank_card.brand":           func(s *fieldSelection) { s.SensitiveBankCardBrand = true },
	"sensitive.bank_card.number":          func(s *fieldSelection) { s.SensitiveBankCardNumber = true },
	"sensitive.bank_card.expiry_month":    func(s *fieldSelection) { s.SensitiveBankCardMonth = true },
	"sensitive.bank_card.expiry_year":     func(s *fieldSelection) { s.SensitiveBankCardYear = true },
	"sensitive.bank_card.cvv":             func(s *fieldSelection) { s.SensitiveBankCardCVV = true },
}

func normalizeIncludeFields(fields []string) (fieldSelection, error) {
	var selection fieldSelection

	for _, field := range fields {
		normalized := strings.ToLower(strings.TrimSpace(field))
		if normalized == "" {
			continue
		}

		apply, ok := supportedIncludeFields[normalized]
		if !ok {
			return fieldSelection{}, &ValidationError{Message: "include_fields contains unsupported field: " + normalized}
		}

		apply(&selection)
	}

	return selection, nil
}

func (s fieldSelection) wantsEducation() bool {
	return s.EducationLevel || s.EducationDegreeType || s.EducationSchoolName || s.EducationMajor ||
		s.EducationGraduationYear || s.EducationStudentStatus || s.EducationSchoolCity ||
		s.EducationSchoolState || s.EducationSchoolEmail
}

func (s fieldSelection) wantsEmploymentEnhancements() bool {
	return s.EmploymentIndustry || s.EmploymentSeniority || s.EmploymentStartDate ||
		s.EmploymentWorkMode || s.EmploymentManager || s.EmploymentOfficeCity || s.EmploymentOfficeState
}

func (s fieldSelection) wantsSensitive() bool {
	return s.SensitiveSSN || s.SensitiveDriverLicenseNo || s.SensitiveDriverLicenseState ||
		s.SensitivePassportNumber || s.SensitivePassportCountry || s.SensitiveBankCardBrand ||
		s.SensitiveBankCardNumber || s.SensitiveBankCardMonth || s.SensitiveBankCardYear ||
		s.SensitiveBankCardCVV
}

func (s fieldSelection) wantsDriverLicense() bool {
	return s.SensitiveDriverLicenseNo || s.SensitiveDriverLicenseState
}

func (s fieldSelection) wantsPassport() bool {
	return s.SensitivePassportNumber || s.SensitivePassportCountry
}

func (s fieldSelection) wantsBankCard() bool {
	return s.SensitiveBankCardBrand || s.SensitiveBankCardNumber || s.SensitiveBankCardMonth ||
		s.SensitiveBankCardYear || s.SensitiveBankCardCVV
}
