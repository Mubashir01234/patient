package models

type PatientRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePatientRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	DOB          string `json:"dob"`
	MobileNumber string `json:"mobile_number"`
	HomeAddress  string `json:"home_address"`
}

type FormRequest struct {
	Patient      NewPatientFormRequest       `json:"patient_detail"`
	EContactInfo EmergencyContactInfoRequest `json:"emergency_contact_information"`
	WorkInfo     WorkInfoRequest             `json:"work_information"`
	HealthInfo   HealthInfoRequest           `json:"health_information"`
	Consent      ConsentRequest              `json:"consent"`
}

type NewPatientFormRequest struct {
	PatientFullName string `json:"patient_full_name" binding:"required"`
	DOB             string `json:"dob"`
	Email           string `json:"email" binding:"required"`
	MobileNumber    string `json:"mobile_number" binding:"required"`
	HomeAddress     string `json:"home_address"`
}

type EmergencyContactInfoRequest struct {
	EContact             string `json:"emergency_contact"`
	Mobile               string `json:"mobile"`
	EContactRelationship string `json:"relationship_to_you"`
}

type WorkInfoRequest struct {
	Occupation string `json:"occupation"`
}

type HealthInfoRequest struct {
	HealthObj           string                   `json:"health_objective" binding:"required"`
	HealthPractitioners bool                     `json:"health_practitioners"`
	PracticeName        string                   `json:"practice_name" binding:"required"`
	MedicationsList     []MedicationsListRequest `json:"medications_list"`
	AllergiesList       []AllergiesListRequest   `json:"allergies_list"`
}

type MedicationsListRequest struct {
	Medications string `json:"medications"`
	Doze        string `json:"doze"`
}

type AllergiesListRequest struct {
	Allergies            string                        `json:"allergies"`
	HospitalisationsList []HospitalisationsListRequest `json:"hospitalisations_list"`
	UploadRelevantScans  string                        `json:"upload_relevant_scans"`
	UploadReports        string                        `json:"upload_reports"`
}

type HospitalisationsListRequest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type ConsentRequest struct {
	ConfirmInfo bool   `json:"confirm_info" binding:"required"`
	Signature   string `json:"signature"`
	Date        string `json:"date"`
}
