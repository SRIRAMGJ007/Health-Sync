CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    fcm_token TEXT,
    google_id TEXT UNIQUE,
    name TEXT,
    age INT,
    gender TEXT,
    blood_group TEXT,
    emergency_contact_number TEXT,
    emergency_contact_relationship TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
    
CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    password_hash TEXT,
    specialization TEXT NOT NULL,
    experience INT NOT NULL,  -- In years
    qualification TEXT NOT NULL,
    hospital_name TEXT NOT NULL,
    consultation_fee DECIMAL(10, 2) NOT NULL,
    contact_number TEXT,
    email TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE doctor_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doctor_id UUID REFERENCES doctors(id),
    availability_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    doctor_id UUID REFERENCES doctors(id),
    availability_id UUID REFERENCES doctor_availability(id),
    booking_date DATE NOT NULL,
    booking_start_time TIME NOT NULL,
    booking_end_time TIME NOT NULL,
    status TEXT NOT NULL, -- e.g., 'pending', 'confirmed', 'canceled', 'completed'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    medication_name TEXT NOT NULL,
    dosage TEXT NOT NULL,
    time_to_notify TIME NOT NULL,
    frequency TEXT NOT NULL CHECK (frequency IN ('daily', 'weekly')),
    is_readbyuser BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);