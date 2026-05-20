CREATE TABLE IF NOT EXISTS reviews (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID         NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    listing_id UUID         NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    guest_id   UUID         NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    rating     INT          NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment    TEXT,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT one_review_per_booking UNIQUE (booking_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_listing_id ON reviews(listing_id);
CREATE INDEX IF NOT EXISTS idx_reviews_guest_id   ON reviews(guest_id);
