-- 創建 passengers 表
CREATE TABLE passengers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    date_of_birth DATE,
    nationality VARCHAR(50),
    passport_number VARCHAR(50),
    passport_expiry DATE,
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    frequent_flyer_number VARCHAR(50) UNIQUE,
    frequent_flyer_tier VARCHAR(20),
    frequent_flyer_points INTEGER DEFAULT 0,
    special_meal_preference VARCHAR(50),
    seat_preference VARCHAR(20),
    special_assistance BOOLEAN DEFAULT FALSE,
    total_flights INTEGER DEFAULT 0,
    total_spent DECIMAL(10, 2) DEFAULT 0,
    last_flight_date DATE,
    marketing_consent BOOLEAN DEFAULT FALSE,
    preferred_language VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建 flights 表
CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    origin VARCHAR(50) NOT NULL,
    destination VARCHAR(50) NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    economy_seats_total INTEGER NOT NULL,
    economy_seats_booked INTEGER DEFAULT 0,
    economy_seats_overbooking_ratio DECIMAL(3, 2) DEFAULT 0,
    business_seats_total INTEGER NOT NULL,
    business_seats_booked INTEGER DEFAULT 0,
    business_seats_overbooking_ratio DECIMAL(3, 2) DEFAULT 0,
    first_class_seats_total INTEGER NOT NULL,
    first_class_seats_booked INTEGER DEFAULT 0,
    first_class_seats_overbooking_ratio DECIMAL(3, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建 bookings 表
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    passenger_id INTEGER REFERENCES passengers(id),
    flight_id INTEGER REFERENCES flights(id),
    class VARCHAR(20) NOT NULL,
    seat_number VARCHAR(10),
    status VARCHAR(20) NOT NULL,
    booking_time TIMESTAMP WITH TIME ZONE NOT NULL,
    check_in_time TIMESTAMP WITH TIME ZONE,
    has_checked_in BOOLEAN DEFAULT FALSE,
    price_amount DECIMAL(10, 2) NOT NULL,
    price_currency VARCHAR(3) NOT NULL,
    compensation_amount DECIMAL(10, 2),
    compensation_currency VARCHAR(3),
    risk_score DECIMAL(3, 2),
    is_cheapest_fare BOOLEAN DEFAULT FALSE,
    special_requests TEXT,
    baggage_checked_bags INTEGER DEFAULT 0,
    baggage_carry_on_bags INTEGER DEFAULT 0,
    baggage_total_weight DECIMAL(5, 2),
    baggage_excess_weight DECIMAL(5, 2),
    baggage_excess_charge_amount DECIMAL(10, 2),
    baggage_excess_charge_currency VARCHAR(3),
    cancellation_time TIMESTAMP WITH TIME ZONE,
    refund_amount DECIMAL(10, 2),
    refund_currency VARCHAR(3),
    is_overbooked BOOLEAN DEFAULT FALSE,
    upgraded_from VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建 historical_data 表
CREATE TABLE historical_data (
    id SERIAL PRIMARY KEY,
    route VARCHAR(100) NOT NULL,
    day_of_week INTEGER NOT NULL,
    average_no_show_rate DECIMAL(4, 3) NOT NULL,
    average_booking_rate DECIMAL(4, 3) NOT NULL,
    average_load_factor DECIMAL(4, 3) NOT NULL,
    average_yield DECIMAL(10, 2) NOT NULL,
    peak_season_start DATE,
    peak_season_end DATE,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引
CREATE INDEX idx_passengers_email ON passengers(email);
CREATE INDEX idx_passengers_frequent_flyer_number ON passengers(frequent_flyer_number);
CREATE INDEX idx_passengers_last_name ON passengers(last_name);

CREATE INDEX idx_flights_origin_destination ON flights(origin, destination);
CREATE INDEX idx_flights_departure_time ON flights(departure_time);

CREATE INDEX idx_bookings_passenger_id ON bookings(passenger_id);
CREATE INDEX idx_bookings_flight_id ON bookings(flight_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_booking_time ON bookings(booking_time);

CREATE INDEX idx_historical_data_route_day ON historical_data(route, day_of_week);
