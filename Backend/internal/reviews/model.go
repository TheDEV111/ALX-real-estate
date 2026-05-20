package reviews

import "time"

type Review struct {
	ID        string    `json:"id"`
	BookingID string    `json:"booking_id"`
	ListingID string    `json:"listing_id"`
	GuestID   string    `json:"guest_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewsResponse struct {
	Data          []Review `json:"data"`
	Total         int      `json:"total"`
	AverageRating float64  `json:"average_rating"`
}
