package repositories

import (
	"context"
	"database/sql"
	"time"

	"airline-booking/models"
)

type PassengerRepository interface {
	CreatePassenger(ctx context.Context, passenger *models.Passenger) error
	GetPassengerByID(ctx context.Context, passengerID int) (*models.Passenger, error)
	UpdatePassenger(ctx context.Context, passenger *models.Passenger) error
	DeletePassenger(ctx context.Context, passengerID int) error
	ListPassengers(ctx context.Context, filter models.PassengerFilter) ([]*models.Passenger, error)
	GetPassengerHistory(ctx context.Context, passengerID int) (*models.PassengerHistory, error)
	UpdateFrequentFlyerPoints(ctx context.Context, passengerID int, pointsToAdd int) error
}

type passengerRepository struct {
	db *sql.DB
}

func NewPassengerRepository(db *sql.DB) PassengerRepository {
	return &passengerRepository{db: db}
}

func (r *passengerRepository) CreatePassenger(ctx context.Context, passenger *models.Passenger) error {
	query := `
        INSERT INTO passengers (
            first_name, last_name, email, phone_number, date_of_birth, nationality,
            passport_number, passport_expiry, address, city, country, postal_code,
            frequent_flyer_number, frequent_flyer_tier, frequent_flyer_points,
            special_meal_preference, seat_preference, special_assistance,
            total_flights, total_spent, last_flight_date,
            marketing_consent, preferred_language, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
            $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
        ) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		passenger.FirstName, passenger.LastName, passenger.Email, passenger.PhoneNumber,
		passenger.DateOfBirth, passenger.Nationality, passenger.PassportNumber,
		passenger.PassportExpiry, passenger.Address, passenger.City, passenger.Country,
		passenger.PostalCode, passenger.FrequentFlyerNumber, passenger.FrequentFlyerTier,
		passenger.FrequentFlyerPoints, passenger.SpecialMealPreference, passenger.SeatPreference,
		passenger.SpecialAssistance, passenger.TotalFlights, passenger.TotalSpent,
		passenger.LastFlightDate, passenger.MarketingConsent, passenger.PreferredLanguage,
		time.Now(), time.Now(),
	).Scan(&passenger.ID)

	return err
}

func (r *passengerRepository) GetPassengerByID(ctx context.Context, passengerID int) (*models.Passenger, error) {
	query := `
        SELECT id, first_name, last_name, email, phone_number, date_of_birth, nationality,
               passport_number, passport_expiry, address, city, country, postal_code,
               frequent_flyer_number, frequent_flyer_tier, frequent_flyer_points,
               special_meal_preference, seat_preference, special_assistance,
               total_flights, total_spent, last_flight_date,
               marketing_consent, preferred_language, created_at, updated_at
        FROM passengers
        WHERE id = $1`

	var passenger models.Passenger
	err := r.db.QueryRowContext(ctx, query, passengerID).Scan(
		&passenger.ID, &passenger.FirstName, &passenger.LastName, &passenger.Email,
		&passenger.PhoneNumber, &passenger.DateOfBirth, &passenger.Nationality,
		&passenger.PassportNumber, &passenger.PassportExpiry, &passenger.Address,
		&passenger.City, &passenger.Country, &passenger.PostalCode,
		&passenger.FrequentFlyerNumber, &passenger.FrequentFlyerTier, &passenger.FrequentFlyerPoints,
		&passenger.SpecialMealPreference, &passenger.SeatPreference, &passenger.SpecialAssistance,
		&passenger.TotalFlights, &passenger.TotalSpent, &passenger.LastFlightDate,
		&passenger.MarketingConsent, &passenger.PreferredLanguage,
		&passenger.CreatedAt, &passenger.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &passenger, nil
}

func (r *passengerRepository) UpdatePassenger(ctx context.Context, passenger *models.Passenger) error {
	query := `
        UPDATE passengers SET
            first_name = $2, last_name = $3, email = $4, phone_number = $5,
            date_of_birth = $6, nationality = $7, passport_number = $8,
            passport_expiry = $9, address = $10, city = $11, country = $12,
            postal_code = $13, frequent_flyer_number = $14, frequent_flyer_tier = $15,
            frequent_flyer_points = $16, special_meal_preference = $17,
            seat_preference = $18, special_assistance = $19, total_flights = $20,
            total_spent = $21, last_flight_date = $22, marketing_consent = $23,
            preferred_language = $24, updated_at = $25
        WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		passenger.ID, passenger.FirstName, passenger.LastName, passenger.Email,
		passenger.PhoneNumber, passenger.DateOfBirth, passenger.Nationality,
		passenger.PassportNumber, passenger.PassportExpiry, passenger.Address,
		passenger.City, passenger.Country, passenger.PostalCode,
		passenger.FrequentFlyerNumber, passenger.FrequentFlyerTier, passenger.FrequentFlyerPoints,
		passenger.SpecialMealPreference, passenger.SeatPreference, passenger.SpecialAssistance,
		passenger.TotalFlights, passenger.TotalSpent, passenger.LastFlightDate,
		passenger.MarketingConsent, passenger.PreferredLanguage, time.Now(),
	)

	return err
}

func (r *passengerRepository) DeletePassenger(ctx context.Context, passengerID int) error {
	query := `DELETE FROM passengers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, passengerID)
	return err
}

func (r *passengerRepository) ListPassengers(ctx context.Context, filter models.PassengerFilter) ([]*models.Passenger, error) {
	query := `
        SELECT id, first_name, last_name, email, phone_number, date_of_birth, nationality,
               passport_number, passport_expiry, address, city, country, postal_code,
               frequent_flyer_number, frequent_flyer_tier, frequent_flyer_points,
               special_meal_preference, seat_preference, special_assistance,
               total_flights, total_spent, last_flight_date,
               marketing_consent, preferred_language, created_at, updated_at
        FROM passengers
        WHERE ($1 = '' OR first_name ILIKE $1 || '%')
          AND ($2 = '' OR last_name ILIKE $2 || '%')
          AND ($3 = '' OR email = $3)
          AND ($4 = '' OR frequent_flyer_number = $4)
          AND ($5 = '' OR nationality = $5)
          AND ($6 = '' OR frequent_flyer_tier = $6)
          AND (total_flights >= $7)
          AND (total_spent >= $8)
          AND (last_flight_date >= $9 OR $9 IS NULL)
        ORDER BY last_name, first_name`

	rows, err := r.db.QueryContext(ctx, query,
		filter.FirstName, filter.LastName, filter.Email,
		filter.FrequentFlyerNumber, filter.Nationality, filter.FrequentFlyerTier,
		filter.MinTotalFlights, filter.MinTotalSpent, filter.LastFlightAfter,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []*models.Passenger
	for rows.Next() {
		var p models.Passenger
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.PhoneNumber,
			&p.DateOfBirth, &p.Nationality, &p.PassportNumber, &p.PassportExpiry,
			&p.Address, &p.City, &p.Country, &p.PostalCode,
			&p.FrequentFlyerNumber, &p.FrequentFlyerTier, &p.FrequentFlyerPoints,
			&p.SpecialMealPreference, &p.SeatPreference, &p.SpecialAssistance,
			&p.TotalFlights, &p.TotalSpent, &p.LastFlightDate,
			&p.MarketingConsent, &p.PreferredLanguage, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		passengers = append(passengers, &p)
	}

	return passengers, nil
}

func (r *passengerRepository) GetPassengerHistory(ctx context.Context, passengerID int) (*models.PassengerHistory, error) {
	query := `
        SELECT 
            CASE WHEN frequent_flyer_tier != '' THEN true ELSE false END as is_frequent_flyer,
            total_flights,
            total_spent,
            last_flight_date
        FROM passengers
        WHERE id = $1`

	var history models.PassengerHistory
	err := r.db.QueryRowContext(ctx, query, passengerID).Scan(
		&history.IsFrequentFlyer,
		&history.TotalFlights,
		&history.TotalSpent,
		&history.LastFlightDate,
	)

	if err != nil {
		return nil, err
	}

	return &history, nil
}

func (r *passengerRepository) UpdateFrequentFlyerPoints(ctx context.Context, passengerID int, pointsToAdd int) error {
	query := `
        UPDATE passengers
        SET frequent_flyer_points = frequent_flyer_points + $2
        WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, passengerID, pointsToAdd)
	return err
}
