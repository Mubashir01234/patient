# Patient

This project backend APIs for patient website facilitate patient registration and appointment scheduling. Patients can register themselves by providing necessary information through a registration form. Additionally, they can submit appointment requests through a form.


## APIs
This patient project contains following APIs:
### 1. Register patient
This API is used for register new patient.
##### Method
```txt
POST
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/register
```
##### Body (raw json)
```json
{
    "email": "patient@gmail.com",
    "password": "patient@01234"
}
```
##### Response
```json
{
    "message": "registration successful"
}
```

### 2. Login patient
This API is used for login patient. When the patient login, a secure JWT token will be generated. This JWT token will be used to make calls to other APIs by setting it as the Authorization in the header.
Admin can also use this login endpoint for accessing this project.
##### Method
```txt
POST
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/login
```
##### Body (raw json)
```json
{
    "email": "patient@gmail.com",
    "password": "patient@01234"
}
```
##### Response
In response it will generates JWT token
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc"
}
```

### 3. Update patient
This API is used to update patient information. A patient can update their existing information using this API or add new information.
##### Method
```txt
PUT
```
##### Request URL
```link
http://localhost:8080/api/v1/patient
```
##### Header
Header is request for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`. The Bearer keyword in the Authorization header of an HTTP request indicates that the request is using the Bearer authentication scheme. This scheme is used to transmit an access token, such as a JSON Web Token (JWT), in the HTTP header to authorize the request.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Body (raw json)
In request body all the options are optional like patient don't want to update email and password, so just remove the these fields from request body. All fields are optional.
```json
{
    "email": "patient@gmail.com",
    "password": "patient@01234",
    "full_name": "patient 1",
    "dob": "2000-01-01",
    "mobile_number": "+61432423523",
    "home_address": "home address"
}
```
##### Response
In response, It will generate new token because patient information was updated.
```json
{
   {
    "message": "patient updated successfully",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjI1Mzd9.l1k1120h7TxaLmPrCh4KREdnfYhT-iMrWRiEwShkVV0"
    }
}
```

### 4. Get patient by email
This API is used to retrieve a patient's profile information based on their email. The patient can view their own profile information, but they do not have access to view other patient's information. However, the admin has the access to view all patient information using the patient's email.
##### Method
```txt
GET
```
##### Request URL
We can pass email in param like that:
```link
http://localhost:8080/api/v1/patient/<email>
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Response
```json
{
    "data": {
        "_id": "650761e8ffcbe9202b01ce29",
        "email": "patient@gmail.com",
        "role": "patient",
        "full_name": "patient 1",
        "dob": "2000-01-01",
        "mobile_number": "+61432423523",
        "home_address": "home address",
        "created_at": "2023-09-17T20:30:32.719Z",
        "updated_at": "2023-09-21T16:42:17.896Z"
    }
}
```

### 5. Delete patient
This API is used to delete a patient based on their email. A patient can only delete their own profile, meaning they cannot delete another patient's profile. However, the admin has the access to delete any patient's profile using the patient's email.
##### Method
```txt
DELETE
```
##### Request URL
We can pass email in param like that:
```link
http://localhost:8080/api/v1/patient/<email>
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Response
```json
{
    "message": "patient deleted successfully"
}
```

### 6. Get all patients
This API is used to retrieve all patients with pagination. Only the admin can access this API. The admin can use this API to retrieve all patient records along with their information. Patients do not have access to use this API.
***Benefits of Backend Pagination***
`Improved Performance`: By fetching and displaying data in smaller chunks or pages, pagination reduces the amount of data that needs to be loaded and processed at once. This results in faster loading times and a more responsive user interface.
##### Method
```txt
GET
```
##### Request URL
```link
http://localhost:8080/api/v1/patients
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Query params
We can add query parameters for pagination so that we only fetch the amount of data we need from the backend like that:
If we can get data for `first page`:
```json
json: {"page": 1, "limit":10}
```
If we can get data for `second page`:
```json
json: {"page": 2, "limit":10}
```

##### Response
The `metadata` will accompany each response, containing necessary information such as the current page(page), the number of items displayed per page(limit), the starting(first_page) and ending(last_page) page, and the total record count(total_records).
```json
{
    "data": [
        {
            "_id": "650761e8ffcbe9202b01ce29",
            "email": "patient@gmail.com",
            "role": "patient",
            "full_name": "patient 1",
            "dob": "2000-01-01",
            "mobile_number": "+61432423523",
            "home_address": "home address",
            "created_at": "2023-09-17T20:30:32.719Z",
            "updated_at": "2023-09-21T16:42:17.896Z"
        },
        {
            "_id": "65075406e5dedfb55efc84bd",
            "email": "user@gmail.com",
            "role": "patient",
            "full_name": "patient 2",
            "dob": "2001-01-01",
            "mobile_number": "+61433572467",
            "home_address": "home address",
            "created_at": "2023-09-17T19:31:18.318Z",
            "updated_at": "2023-09-17T19:33:26.036Z"
        }
    ],
    "metadata": {
        "page": 1,
        "limit": 10,
        "first_page": 1,
        "last_page": 1,
        "total_records": 2
    }
}
```

### 7. Upload file/documents
This API is used for uploading scanned documents for a form. The API will store the document in the database and send back the file details in the response. In the response, we will use the `id` to link the file when submitting the form.
##### Method
```txt
POST
```
##### Request URL
```link
http://localhost:8080/api/v1/file/upload
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Body (form-data)
In the `form-data`, we will have key-value pairs. For the key, we will set `file`, and for the value, we will `select the file to upload`.
```txt
file: /file/to/path
```
##### Response
We can use this `id` in uploading fields in submit form api.
```json
{
    "file": {
        "id": "650d305aca772c5b38cbe965",
        "name": "file.pdf",
        "size": 147930
    },
    "msg": "file uploaded successfully"
}
```

### 8. View uploaded documents
This API is used to view an uploaded document using the provided `id` as a parameter.
##### Method
```txt
GET
```
##### Request URL
```link
http://localhost:8080/api/v1/file/view/<file-id>
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```

##### Response
The response will display the uploaded document.


### 9. Form submit
This API is used to upload a patient's medical form with all the necessary details.
##### Method
```txt
POST
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/form
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Body (raw json)
In request body includes all the fields of medicial form.
```json
{
    "patient_detail": {
        "patient_full_name": "John Doe",
        "dob": "1980-01-01",
        "email": "johndoe@example.com",
        "mobile_number": "+1234567890",
        "home_address": "123 Main St"
    },
    "emergency_contact_information": {
        "emergency_contact": "Jane Smith",
        "mobile": "+9876543210",
        "relationship_to_you": "Spouse"
    },
    "work_information": {
        "occupation": "Software Engineer"
    },
    "health_information": {
        "health_objective": "Maintain overall health",
        "health_practitioners": true,
        "practice_name": "Healthy Clinic",
        "medications_list": [
            {
                "medications": "Medication A",
                "doze": "1 pill daily"
            },
            {
                "medications": "Medication B",
                "doze": "2 pills twice daily"
            }
        ],
        "allergies_list": [
            {
                "allergies": "Pollen",
                "hospitalisations_list": [
                    {
                        "name": "Hospital A",
                        "date": "2022-03-15"
                    },
                    {
                        "name": "Hospital B",
                        "date": "2019-07-20"
                    }
                ],
                "upload_relevant_scans": "6506d67e33abf3d5b7c8d007", // When a user adds a document, we will first upload the file and then link the uploaded document's ID in the respective field.
                "upload_reports": "6506d67e33abf3d5b7c8d012" // we all this ID after the uploading the document
            }
        ]
    },
    "consent": {
        "confirm_info": true,
        "signature": "6506d67e33abf3d5b7c8d002", // we all this ID after the uploading the document
        "date": "2023-09-17"
    }
}
```

##### Response
```json
{
    "form_id": "650d37c6ca772c5b38cbe969",
    "message": "form submitted successfully"
}
```

### 10. Get form by form id
This API is used to fetch a submitted form using the `form id`. With this API, a patient can only retrieve their own submitted forms, while they cannot access to forms submitted by other patients. However, the admin has the capability to access to fetch all the submitted forms of all patients.
##### Method
```txt
GET
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/form/<form-id>
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```

##### Response
```json
{
    "data": {
        "patient_id": "650761e8ffcbe9202b01ce29",
        "patient_email": "patient@gmail.com",
        "patient_detail": {
            "patient_full_name": "John Doe",
            "dob": "1980-01-01",
            "email": "patient@gmail.com",
            "mobile_number": "+1234567890",
            "home_address": "123 Main St"
        },
        "emergency_contact_information ": {
            "emergency_contact": "Jane Smith",
            "mobile": "+9876543210",
            "relationship_to_you": "Spouse"
        },
        "work_information": {
            "occupation": "Software Engineer"
        },
        "health_information": {
            "health_objective": "Maintain overall health",
            "health_practitioners": true,
            "practice_name": "Healthy Clinic",
            "medications_list": [
                {
                    "medications": "Medication A",
                    "doze": "1 pill daily"
                },
                {
                    "medications": "Medication B",
                    "doze": "2 pills twice daily"
                }
            ],
            "allergies_list": [
                {
                    "allergies": "Pollen",
                    "hospitalisations_list": [
                        {
                            "name": "Hospital A",
                            "date": "2022-03-15"
                        },
                        {
                            "name": "Hospital B",
                            "date": "2019-07-20"
                        }
                    ],
                    "upload_relevant_scans": "6506d67e33abf3d5b7c8d007",
                    "upload_reports": "6506d67e33abf3d5b7c8d012"
                }
            ]
        },
        "consent": {
            "confirm_info": true,
            "signature": "6506d67e33abf3d5b7c8d002",
            "date": "2023-09-17"
        },
        "created_at": "2023-09-22T06:44:22.897Z",
        "updated_at": "2023-09-22T06:44:22.897Z"
    }
}
```

### 11. Get patient all forms by email
This API is used to fetch all forms of a patient with pagination based on their email. Patients can only retrieve their own forms using this API. However, the admin has the access to retrieve all forms of any patient using their respective email.
***Benefits of Backend Pagination***
`Improved Performance`: By fetching and displaying data in smaller chunks or pages, pagination reduces the amount of data that needs to be loaded and processed at once. This results in faster loading times and a more responsive user interface.
##### Method
```txt
GET
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/forms/<email>
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Query params
We can add query parameters for pagination so that we only fetch the amount of data we need from the backend like that:
If we can get data for `first page`:
```json
json: {"page": 1, "limit":10}
```
If we can get data for `second page`:
```json
json: {"page": 2, "limit":10}
```
##### Response
The `metadata` will accompany each response, containing necessary information such as the current page(page), the number of items displayed per page(limit), the starting(first_page) and ending(last_page) page, and the total record count(total_records).
```json
{
    "data": [
        {
            "patient_id": "650761e8ffcbe9202b01ce29",
            "patient_email": "patient@gmail.com",
            "patient_detail": {
                "patient_full_name": "John Doe 2",
                "dob": "1980-01-01",
                "email": "patient@gmail.com",
                "mobile_number": "+1234567890",
                "home_address": "123 Main St"
            },
            "emergency_contact_information ": {
                "emergency_contact": "Jane Smith",
                "mobile": "+9876543210",
                "relationship_to_you": "Spouse"
            },
            "work_information": {
                "occupation": "Software Engineer"
            },
            "health_information": {
                "health_objective": "Maintain overall health",
                "health_practitioners": true,
                "practice_name": "Healthy Clinic",
                "medications_list": [
                    {
                        "medications": "Medication A",
                        "doze": "1 pill daily"
                    },
                    {
                        "medications": "Medication B",
                        "doze": "2 pills twice daily"
                    }
                ],
                "allergies_list": [
                    {
                        "allergies": "Pollen",
                        "hospitalisations_list": [
                            {
                                "name": "Hospital A",
                                "date": "2022-03-15"
                            },
                            {
                                "name": "Hospital B",
                                "date": "2019-07-20"
                            }
                        ],
                        "upload_relevant_scans": "6506d67e33abf3d5b7c8d007",
                        "upload_reports": "6506d67e33abf3d5b7c8d007"
                    }
                ]
            },
            "consent": {
                "confirm_info": true,
                "signature": "John Doe",
                "date": "2023-09-17"
            },
            "created_at": "2023-09-22T07:34:29.462Z",
            "updated_at": "2023-09-22T07:34:29.462Z"
        },
        {
            "patient_id": "650761e8ffcbe9202b01ce29",
            "patient_email": "patient@gmail.com",
            "patient_detail": {
                "patient_full_name": "John Doe",
                "dob": "1980-01-01",
                "email": "johndoe@example.com",
                "mobile_number": "+1234567890",
                "home_address": "123 Main St"
            },
            "emergency_contact_information ": {
                "emergency_contact": "Jane Smith",
                "mobile": "+9876543210",
                "relationship_to_you": "Spouse"
            },
            "work_information": {
                "occupation": "Software Engineer"
            },
            "health_information": {
                "health_objective": "Maintain overall health",
                "health_practitioners": true,
                "practice_name": "Healthy Clinic",
                "medications_list": [
                    {
                        "medications": "Medication A",
                        "doze": "1 pill daily"
                    },
                    {
                        "medications": "Medication B",
                        "doze": "2 pills twice daily"
                    }
                ],
                "allergies_list": [
                    {
                        "allergies": "Pollen",
                        "hospitalisations_list": [
                            {
                                "name": "Hospital A",
                                "date": "2022-03-15"
                            },
                            {
                                "name": "Hospital B",
                                "date": "2019-07-20"
                            }
                        ],
                        "upload_relevant_scans": "6506d67e33abf3d5b7c8d007",
                        "upload_reports": "6506d67e33abf3d5b7c8d012"
                    }
                ]
            },
            "consent": {
                "confirm_info": true,
                "signature": "6506d67e33abf3d5b7c8d023",
                "date": "2023-09-17"
            },
            "created_at": "2023-09-22T06:44:22.897Z",
            "updated_at": "2023-09-22T06:44:22.897Z"
        }
    ],
    "metadata": {
        "page": 1,
        "limit": 10,
        "first_page": 1,
        "last_page": 1,
        "total_records": 2
    }
}
```

### 12. Get all the forms
This API is used to fetch all forms submitted by all patients. Only the admin has access to this API and can retrieve all forms using pagination. The latest data will be displayed at the top.
##### Method
```txt
GET
```
##### Request URL
```link
http://localhost:8080/api/v1/patient/forms
```
##### Header
Header is required for this request. We need to set `Authorization` and `JWT token` as a value including `Bearer`.
```txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50X2lkIjoiNjUwNzYxZThmZmNiZTkyMDJiMDFjZTI5IiwiZW1haWwiOiJhZG1pbkBnbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2OTU0MjEzNDd9.jw0M8xDoum_uoEPzZ2vv1Y-9r7WDiyTqm1YkikEANoc
```
##### Query params
We can add query parameters for pagination so that we only fetch the amount of data we need from the backend like that:
If we can get data for `first page`:
```json
json: {"page": 1, "limit":10}
```
If we can get data for `second page`:
```json
json: {"page": 2, "limit":10}
```

##### Response
The `metadata` will accompany each response, containing necessary information such as the current page(page), the number of items displayed per page(limit), the starting(first_page) and ending(last_page) page, and the total record count(total_records).
```json
{
    "data": [
        {
            "patient_id": "650761e8ffcbe9202b01ce29",
            "patient_email": "admin@gmail.com",
            "patient_detail": {
                "patient_full_name": "John Doe 2",
                "dob": "1980-01-01",
                "email": "johndoe@example.com",
                "mobile_number": "+1234567890",
                "home_address": "123 Main St"
            },
            "emergency_contact_information ": {
                "emergency_contact": "Jane Smith",
                "mobile": "+9876543210",
                "relationship_to_you": "Spouse"
            },
            "work_information": {
                "occupation": "Software Engineer"
            },
            "health_information": {
                "health_objective": "Maintain overall health",
                "health_practitioners": true,
                "practice_name": "Healthy Clinic",
                "medications_list": [
                    {
                        "medications": "Medication A",
                        "doze": "1 pill daily"
                    },
                    {
                        "medications": "Medication B",
                        "doze": "2 pills twice daily"
                    }
                ],
                "allergies_list": [
                    {
                        "allergies": "Pollen",
                        "hospitalisations_list": [
                            {
                                "name": "Hospital A",
                                "date": "2022-03-15"
                            },
                            {
                                "name": "Hospital B",
                                "date": "2019-07-20"
                            }
                        ],
                        "upload_relevant_scans": "6506d67e33abf3d5b7c8d007",
                        "upload_reports": "6506d67e33abf3d5b7c8d007"
                    }
                ]
            },
            "consent": {
                "confirm_info": true,
                "signature": "John Doe",
                "date": "2023-09-17"
            },
            "created_at": "2023-09-22T07:34:29.462Z",
            "updated_at": "2023-09-22T07:34:29.462Z"
        },
        {
            "patient_id": "650761e8ffcbe9202b01ce29",
            "patient_email": "admin@gmail.com",
            "patient_detail": {
                "patient_full_name": "John Doe",
                "dob": "1980-01-01",
                "email": "johndoe@example.com",
                "mobile_number": "+1234567890",
                "home_address": "123 Main St"
            },
            "emergency_contact_information ": {
                "emergency_contact": "Jane Smith",
                "mobile": "+9876543210",
                "relationship_to_you": "Spouse"
            },
            "work_information": {
                "occupation": "Software Engineer"
            },
            "health_information": {
                "health_objective": "Maintain overall health",
                "health_practitioners": true,
                "practice_name": "Healthy Clinic",
                "medications_list": [
                    {
                        "medications": "Medication A",
                        "doze": "1 pill daily"
                    },
                    {
                        "medications": "Medication B",
                        "doze": "2 pills twice daily"
                    }
                ],
                "allergies_list": [
                    {
                        "allergies": "Pollen",
                        "hospitalisations_list": [
                            {
                                "name": "Hospital A",
                                "date": "2022-03-15"
                            },
                            {
                                "name": "Hospital B",
                                "date": "2019-07-20"
                            }
                        ],
                        "upload_relevant_scans": "6506d67e33abf3d5b7c8d007",
                        "upload_reports": "6506d67e33abf3d5b7c8d007"
                    }
                ]
            },
            "consent": {
                "confirm_info": true,
                "signature": "John Doe",
                "date": "2023-09-17"
            },
            "created_at": "2023-09-22T06:44:22.897Z",
            "updated_at": "2023-09-22T06:44:22.897Z"
        },
    ...
    ...
    ...
    ...
    ],
    "metadata": {
        "page": 1,
        "limit": 10,
        "first_page": 1,
        "last_page": 2,
        "total_records": 11
    }
}
```
