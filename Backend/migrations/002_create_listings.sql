CREATE TABLE IF NOT EXISTS listings (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id         UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT           NOT NULL,
    description     TEXT           NOT NULL,
    property_type   TEXT           NOT NULL CHECK (property_type IN ('apartment','house','villa','cabin','studio')),

    country         TEXT           NOT NULL,
    city            TEXT           NOT NULL,
    address         TEXT           NOT NULL,
    latitude        NUMERIC(9,6),
    longitude       NUMERIC(9,6),

    price_per_night NUMERIC(10,2)  NOT NULL CHECK (price_per_night > 0),
    currency        TEXT           NOT NULL DEFAULT 'USD',

    max_guests      INT            NOT NULL DEFAULT 1 CHECK (max_guests > 0),
    bedrooms        INT            NOT NULL DEFAULT 1,
    bathrooms       NUMERIC(3,1)   NOT NULL DEFAULT 1,

    amenities       TEXT[]         NOT NULL DEFAULT '{}',
    images          TEXT[]         NOT NULL DEFAULT '{}',

    is_active       BOOLEAN        NOT NULL DEFAULT TRUE,

    search_vector   TSVECTOR,

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_listings_host_id  ON listings(host_id);
CREATE INDEX IF NOT EXISTS idx_listings_city     ON listings(city);
CREATE INDEX IF NOT EXISTS idx_listings_price    ON listings(price_per_night);
CREATE INDEX IF NOT EXISTS idx_listings_active   ON listings(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_listings_search   ON listings USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_listings_amenities ON listings USING GIN(amenities);

CREATE OR REPLACE FUNCTION listings_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.city, '')), 'A')  ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(NEW.country, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER listings_search_vector_trigger
BEFORE INSERT OR UPDATE ON listings
FOR EACH ROW EXECUTE FUNCTION listings_search_vector_update();
