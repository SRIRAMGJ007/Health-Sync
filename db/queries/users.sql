-- name: CreateUserWithEmail :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING id, email, name, created_at, updated_at;


-- name: CreateUserWithGoogle :one
INSERT INTO users (google_id, email, name)
VALUES ($1, $2, $3)
RETURNING id, email, google_id, name;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, google_id, name, created_at, updated_at
FROM users
WHERE email = $1;


-- name: GetUserByGoogleID :one
SELECT id, email, google_id, name, created_at, updated_at
FROM users
WHERE google_id = $1;


-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE email = $1;

-- name: UpdateUserProfile :exec
UPDATE users
SET 
    name = COALESCE($1, name),
    age = COALESCE($2, age),
    gender = COALESCE($3, gender),
    blood_group = COALESCE($4, blood_group),
    emergency_contact_number = COALESCE($5, emergency_contact_number),
    emergency_contact_relationship = COALESCE($6, emergency_contact_relationship),
    updated_at = NOW()
WHERE id = $7;


-- name: GetUserByID :one
select id, email, name, age, gender, blood_group, emergency_contact_number, emergency_contact_relationship, updated_at 
FROM users 
WHERE id = $1;

-- name: GetUserProfileByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateDoctor :one
INSERT INTO doctors (
    name,
    specialization,
    experience,
    qualification,
    hospital_name,
    consultation_fee,
    contact_number,
    email,
    password_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, name, specialization, experience, qualification, hospital_name, consultation_fee, contact_number, email, created_at, updated_at;

-- name: GetDoctorByEmail :one
SELECT id, name, specialization, experience, qualification, hospital_name, consultation_fee, contact_number, email, password_hash, created_at, updated_at
FROM doctors
WHERE email = $1;

-- name: GetDoctorByID :one
SELECT *
FROM doctors
WHERE id = $1;

-- name: ListDoctors :many
SELECT *
FROM doctors;

-- name: CreateDoctorAvailability :one
INSERT INTO doctor_availability (doctor_id, availability_date, start_time, end_time)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDoctorAvailabilityByID :one
SELECT *
FROM doctor_availability
WHERE id = $1;

-- name: GetDoctorAvailabilityByDoctorAndDate :many
SELECT *
FROM doctor_availability
WHERE doctor_id = $1 AND availability_date = $2;

-- name: GetDoctorAvailabilityByDoctor :many
SELECT *
FROM doctor_availability
WHERE doctor_id = $1;

-- name: UpdateDoctorAvailability :exec
UPDATE doctor_availability
SET
    start_time = COALESCE($1, start_time),
    end_time = COALESCE($2, end_time),
    is_booked = COALESCE($3, is_booked),
    updated_at = NOW()
WHERE id = $4 AND doctor_id = $5;

-- name: DeleteDoctorAvailability :exec
DELETE FROM doctor_availability
WHERE id = $1 AND doctor_id = $2;

-- name: UpdateAvailabilityBookedStatus :exec
UPDATE doctor_availability
SET
    is_booked = $1,
    updated_at = NOW()
WHERE id = $2;



-- name: CreateBooking :one
INSERT INTO bookings (user_id, doctor_id, availability_id, booking_date, booking_start_time, booking_end_time, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBookingByID :one
SELECT *
FROM bookings
WHERE id = $1;

-- name: GetBookingsByUserID :many
SELECT *
FROM bookings
WHERE user_id = $1;

-- name: GetBookingsByDoctorID :many
SELECT *
FROM bookings
WHERE doctor_id = $1;

-- name: GetBookingsByAvailabilityID :many
SELECT *
FROM bookings
WHERE availability_id = $1;

-- name: DeleteBooking :exec
DELETE FROM bookings
WHERE id = $1;

-- name: UpdateBookingStatus :exec
UPDATE bookings
SET
    status = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: CreateMedication :one
INSERT INTO medications (user_id, medication_name, dosage, time_to_notify, frequency)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMedicationsToNotify :many
SELECT *
FROM medications
WHERE time_to_notify = $1
  AND (frequency = 'daily' OR (frequency = 'weekly' AND EXTRACT(DOW FROM NOW()) = EXTRACT(DOW FROM updated_at)));

-- name: UpdateMedicationReadStatus :exec
UPDATE medications
SET is_readbyuser = TRUE
WHERE id = $1;

-- name: GetUserFCMToken :one
SELECT fcm_token
FROM users
WHERE id = $1;