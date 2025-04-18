// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Booking struct {
	ID               pgtype.UUID
	UserID           pgtype.UUID
	DoctorID         pgtype.UUID
	AvailabilityID   pgtype.UUID
	BookingDate      pgtype.Date
	BookingStartTime pgtype.Time
	BookingEndTime   pgtype.Time
	Status           string
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

type Doctor struct {
	ID              pgtype.UUID
	Name            string
	PasswordHash    *string
	Specialization  string
	Experience      int32
	Qualification   string
	HospitalName    string
	ConsultationFee pgtype.Numeric
	ContactNumber   *string
	Email           *string
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}

type DoctorAvailability struct {
	ID               pgtype.UUID
	DoctorID         pgtype.UUID
	AvailabilityDate pgtype.Date
	StartTime        pgtype.Time
	EndTime          pgtype.Time
	IsBooked         *bool
	CreatedAt        pgtype.Timestamp
	UpdatedAt        pgtype.Timestamp
}

type EncryptedFile struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	FileName  string
	FileData  []byte
	CreatedAt pgtype.Timestamp
}

type Medication struct {
	ID             pgtype.UUID
	UserID         pgtype.UUID
	MedicationName string
	Dosage         string
	TimeToNotify   pgtype.Time
	Frequency      string
	IsReadbyuser   *bool
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

type User struct {
	ID                           pgtype.UUID
	Email                        string
	PasswordHash                 *string
	FcmToken                     *string
	GoogleID                     *string
	Name                         *string
	Age                          *int32
	Gender                       *string
	BloodGroup                   *string
	EmergencyContactNumber       *string
	EmergencyContactRelationship *string
	CreatedAt                    pgtype.Timestamp
	UpdatedAt                    pgtype.Timestamp
}
