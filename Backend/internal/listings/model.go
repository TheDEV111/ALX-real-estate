package listings

import "time"

type Listing struct {
	ID            string    `json:"id"`
	HostID        string    `json:"host_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	PropertyType  string    `json:"property_type"`
	Country       string    `json:"country"`
	City          string    `json:"city"`
	Address       string    `json:"address"`
	Latitude      *float64  `json:"latitude,omitempty"`
	Longitude     *float64  `json:"longitude,omitempty"`
	PricePerNight float64   `json:"price_per_night"`
	Currency      string    `json:"currency"`
	MaxGuests     int       `json:"max_guests"`
	Bedrooms      int       `json:"bedrooms"`
	Bathrooms     float64   `json:"bathrooms"`
	Amenities     []string  `json:"amenities"`
	Images        []string  `json:"images"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateListingRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	PropertyType  string   `json:"property_type"`
	Country       string   `json:"country"`
	City          string   `json:"city"`
	Address       string   `json:"address"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	PricePerNight float64  `json:"price_per_night"`
	Currency      string   `json:"currency"`
	MaxGuests     int      `json:"max_guests"`
	Bedrooms      int      `json:"bedrooms"`
	Bathrooms     float64  `json:"bathrooms"`
	Amenities     []string `json:"amenities"`
	Images        []string `json:"images"`
}

type UpdateListingRequest struct {
	Title         *string  `json:"title"`
	Description   *string  `json:"description"`
	PropertyType  *string  `json:"property_type"`
	Country       *string  `json:"country"`
	City          *string  `json:"city"`
	Address       *string  `json:"address"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	PricePerNight *float64 `json:"price_per_night"`
	Currency      *string  `json:"currency"`
	MaxGuests     *int     `json:"max_guests"`
	Bedrooms      *int     `json:"bedrooms"`
	Bathrooms     *float64 `json:"bathrooms"`
	Amenities     []string `json:"amenities"`
	Images        []string `json:"images"`
	IsActive      *bool    `json:"is_active"`
}

type SearchParams struct {
	Query     *string
	City      *string
	MinPrice  *float64
	MaxPrice  *float64
	Guests    *int
	Amenities []string
	CheckIn   *string
	CheckOut  *string
	Limit     int
	Offset    int
}

type ListingsResponse struct {
	Data   []Listing `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	Page   int       `json:"page"`
}
