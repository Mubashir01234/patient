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
