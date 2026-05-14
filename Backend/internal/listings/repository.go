package listings

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, l *Listing) error {
	query := `
		INSERT INTO listings
			(host_id, title, description, property_type, country, city, address,
			 latitude, longitude, price_per_night, currency, max_guests, bedrooms,
			 bathrooms, amenities, images)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		l.HostID, l.Title, l.Description, l.PropertyType,
		l.Country, l.City, l.Address, l.Latitude, l.Longitude,
		l.PricePerNight, l.Currency, l.MaxGuests, l.Bedrooms,
		l.Bathrooms, l.Amenities, l.Images,
	).Scan(&l.ID, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Listing, error) {
	query := `
		SELECT id, host_id, title, description, property_type, country, city, address,
		       latitude, longitude, price_per_night, currency, max_guests, bedrooms,
		       bathrooms, amenities, images, is_active, created_at, updated_at
		FROM listings
		WHERE id = $1
	`
	l := &Listing{}
	var amenities, images pgtype.TextArray

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.HostID, &l.Title, &l.Description, &l.PropertyType,
		&l.Country, &l.City, &l.Address, &l.Latitude, &l.Longitude,
		&l.PricePerNight, &l.Currency, &l.MaxGuests, &l.Bedrooms,
		&l.Bathrooms, &amenities, &images, &l.IsActive, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}

	l.Amenities = textArrayToSlice(amenities)
	l.Images = textArrayToSlice(images)
	return l, nil
}

func (r *Repository) Search(ctx context.Context, p SearchParams) ([]Listing, int, error) {
	var amenities interface{}
	if len(p.Amenities) > 0 {
		amenities = p.Amenities
	}

	query := `
		SELECT id, host_id, title, description, property_type, country, city, address,
		       latitude, longitude, price_per_night, currency, max_guests, bedrooms,
		       bathrooms, amenities, images, is_active, created_at, updated_at,
		       COUNT(*) OVER() AS total_count
		FROM listings
		WHERE is_active = TRUE
		  AND ($1::text IS NULL OR search_vector @@ plainto_tsquery('english', $1))
		  AND ($2::text IS NULL OR city ILIKE '%' || $2 || '%')
		  AND ($3::numeric IS NULL OR price_per_night >= $3)
		  AND ($4::numeric IS NULL OR price_per_night <= $4)
		  AND ($5::int IS NULL OR max_guests >= $5)
		  AND ($6::text[] IS NULL OR amenities @> $6)
		  AND (
		      $7::date IS NULL OR $8::date IS NULL
		      OR id NOT IN (
		          SELECT listing_id FROM bookings
		          WHERE status != 'cancelled'
		            AND daterange(check_in, check_out, '[)') && daterange($7::date, $8::date, '[)')
		      )
		  )
		ORDER BY
		    CASE WHEN $1::text IS NOT NULL
		         THEN ts_rank(search_vector, plainto_tsquery('english', $1))
		         ELSE 0 END DESC,
		    created_at DESC
		LIMIT $9 OFFSET $10
	`

	rows, err := r.pool.Query(ctx, query,
		p.Query, p.City, p.MinPrice, p.MaxPrice, p.Guests,
		amenities, p.CheckIn, p.CheckOut, p.Limit, p.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("Search: %w", err)
	}
	defer rows.Close()

	var results []Listing
	var total int

	for rows.Next() {
		l := Listing{}
		var am, im pgtype.TextArray
		err := rows.Scan(
			&l.ID, &l.HostID, &l.Title, &l.Description, &l.PropertyType,
			&l.Country, &l.City, &l.Address, &l.Latitude, &l.Longitude,
			&l.PricePerNight, &l.Currency, &l.MaxGuests, &l.Bedrooms,
			&l.Bathrooms, &am, &im, &l.IsActive, &l.CreatedAt, &l.UpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("Search scan: %w", err)
		}
		l.Amenities = textArrayToSlice(am)
		l.Images = textArrayToSlice(im)
		results = append(results, l)
	}

	if results == nil {
		results = []Listing{}
	}
	return results, total, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateListingRequest) (*Listing, error) {
	setClauses := []string{}
	args := []interface{}{}
	i := 1

	addField := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	if req.Title != nil {
		addField("title", *req.Title)
	}
	if req.Description != nil {
		addField("description", *req.Description)
	}
	if req.PropertyType != nil {
		addField("property_type", *req.PropertyType)
	}
	if req.Country != nil {
		addField("country", *req.Country)
	}
	if req.City != nil {
		addField("city", *req.City)
	}
	if req.Address != nil {
		addField("address", *req.Address)
	}
	if req.Latitude != nil {
		addField("latitude", *req.Latitude)
	}
	if req.Longitude != nil {
		addField("longitude", *req.Longitude)
	}
	if req.PricePerNight != nil {
		addField("price_per_night", *req.PricePerNight)
	}
	if req.Currency != nil {
		addField("currency", *req.Currency)
	}
	if req.MaxGuests != nil {
		addField("max_guests", *req.MaxGuests)
	}
	if req.Bedrooms != nil {
		addField("bedrooms", *req.Bedrooms)
	}
	if req.Bathrooms != nil {
		addField("bathrooms", *req.Bathrooms)
	}
	if req.Amenities != nil {
		addField("amenities", req.Amenities)
	}
	if req.Images != nil {
		addField("images", req.Images)
	}
	if req.IsActive != nil {
		addField("is_active", *req.IsActive)
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = NOW()"))
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE listings SET %s
		WHERE id = $%d
		RETURNING id, host_id, title, description, property_type, country, city, address,
		          latitude, longitude, price_per_night, currency, max_guests, bedrooms,
		          bathrooms, amenities, images, is_active, created_at, updated_at
	`, strings.Join(setClauses, ", "), i)

	l := &Listing{}
	var amenities, images pgtype.TextArray

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&l.ID, &l.HostID, &l.Title, &l.Description, &l.PropertyType,
		&l.Country, &l.City, &l.Address, &l.Latitude, &l.Longitude,
		&l.PricePerNight, &l.Currency, &l.MaxGuests, &l.Bedrooms,
		&l.Bathrooms, &amenities, &images, &l.IsActive, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	l.Amenities = textArrayToSlice(amenities)
	l.Images = textArrayToSlice(images)
	return l, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM listings WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func textArrayToSlice(arr pgtype.TextArray) []string {
	result := make([]string, 0, len(arr.Elements))
	for _, e := range arr.Elements {
		if e.Status == pgtype.Present {
			result = append(result, e.String)
		}
	}
	return result
}
