package bookings

import "time"

type Booking struct {
	ID          string    `json:"id"`
	ListingID   string    `json:"listing_id"`
	GuestID     string    `json:"guest_id"`
	CheckIn     string    `json:"check_in"`
	CheckOut    string    `json:"check_out"`
	TotalNights int       `json:"total_nights"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateBookingRequest struct {
	ListingID string `json:"listing_id"`
	CheckIn   string `json:"check_in"`
	CheckOut  string `json:"check_out"`
}

type BookingsResponse struct {
	Data   []Booking `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
