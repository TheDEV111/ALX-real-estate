CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS bookings (
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id   UUID           NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    guest_id     UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in     DATE           NOT NULL,
    check_out    DATE           NOT NULL,
    total_price  NUMERIC(10,2)  NOT NULL,
    status       TEXT           NOT NULL DEFAULT 'pending'
                     CHECK (status IN ('pending','confirmed','cancelled','completed')),
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_dates CHECK (check_out > check_in),
    CONSTRAINT no_overlap EXCLUDE USING GIST (
        listing_id WITH =,
        daterange(check_in, check_out, '[)') WITH &&
    ) WHERE (status NOT IN ('cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_bookings_listing_id ON bookings(listing_id);
CREATE INDEX IF NOT EXISTS idx_bookings_guest_id   ON bookings(guest_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status     ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_dates      ON bookings USING GIST(daterange(check_in, check_out, '[)'));
