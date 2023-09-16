package models

type PatientRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePatientRequest struct {
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	FullName    string `json:"full_name" bson:"full_name"`
	DOB         string `json:"dob" bson:"dob"`
	Phone       string `json:"phone" bson:"phone"`
	HomeAddress string `json:"home_address" bson:"home_address"`
}
