package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginPatient struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Patient struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Role         string             `json:"role" bson:"role"`
	FullName     string             `json:"full_name" bson:"full_name"`
	DOB          string             `json:"dob" bson:"dob"`
	MobileNumber string             `json:"mobile_number" bson:"mobile_number"`
	HomeAddress  string             `json:"home_address" bson:"home_address"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type GetPatient struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Form struct {
	Patient      NewPatientForm       `json:"patient_detail" bson:"patient_detail"`
	EContactInfo EmergencyContactInfo `json:"emergency_contact_information " bson:"emergency_contact_information"`
	WorkInfo     WorkInfo             `json:"work_information" bson:"work_information"`
	HealthInfo   HealthInfo           `json:"health_information" bson:"health_information"`
	Consent      Consent              `json:"consent" bson:"consent"`
}

type NewPatientForm struct {
	PatientFullName string `json:"patient_full_name" bson:"patient_full_name"`
	DOB             string `json:"dob" bson:"dob"`
	Email           string `json:"email" bson:"email"`
	MobileNumber    string `json:"mobile_number" bson:"mobile_number"`
	HomeAddress     string `json:"home_address" bson:"home_address"`
}

type EmergencyContactInfo struct {
	EContact             string `json:"emergency_contact" bson:"emergency_contact"`
	Mobile               string `json:"mobile" bson:"mobile"`
	EContactRelationship string `json:"relationship_to_you" bson:"relationship_to_you"`
}

type WorkInfo struct {
	Occupation string `json:"occupation" bson:"occupation"`
}

type HealthInfo struct {
	HealthObj           string            `json:"health_objective" bson:"health_objective"`
	HealthPractitioners bool              `json:"health_practitioners" bson:"health_practitioner"`
	PracticeName        string            `json:"practice_name" bson:"practice_name"`
	MedicationsList     []MedicationsList `json:"medications_list" bson:"medications_list"`
	AllergiesList       []AllergiesList   `json:"allergies_list" bson:"allergies_list"`
}

type MedicationsList struct {
	Medications string `json:"medications" bson:"medications"`
	Doze        string `json:"doze" bson:"doze"`
}

type AllergiesList struct {
	Allergies            string                 `json:"allergies" bson:"allergies"`
	HospitalisationsList []HospitalisationsList `json:"hospitalisations_list" bson:"hospitalisations_list"`
	UploadRelevantScans  string                 `json:"upload_relevant_scans" bson:"upload_relevant_scans"`
	UploadReports        string                 `json:"upload_reports" bson:"upload_reports"`
}

type HospitalisationsList struct {
	Name string `json:"name" bson:"name"`
	Date string `json:"date" bson:"date"`
}

type Consent struct {
	ConfirmInfo bool   `json:"confirm_info" bson:"confirm_info"`
	Signature   string `json:"signature" bson:"signature"`
	Date        string `json:"date" bson:"date"`
}
