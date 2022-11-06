package db

import (
	"context"
	"time"
)

type Listing struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateListingParam struct {
	Name string
}

type UpdateListingParam struct {
	Name string
}

type ListListingParam struct {
	Offset int32
	Limit  int32
}

func (store *Store) GetListingByID(ctx context.Context, id int64) (listing Listing, err error) {

	const query = `SELECT * FROM "listings" WHERE "id" = $1`
	err = store.db.GetContext(ctx, &listing, query, id)

	return
}

func (store *Store) GetAllListings(ctx context.Context, arg ListListingParam) (users []Listing, err error) {

	const query = `SELECT * FROM "listings" OFFSET $1 LIMIT $2`
	users = []Listing{}
	err = store.db.SelectContext(ctx, &users, query, arg.Offset, arg.Limit)

	return
}

func (store *Store) CreateListing(ctx context.Context, arg CreateListingParam) (Listing, error) {

	const query = `
	INSERT INTO "listings" ("name") 
	VALUES ($1)
	RETURNING "id", "name", "created_at"
	`
	row := store.db.QueryRowContext(ctx, query, arg.Name)

	var user Listing
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
	)

	return user, err
}
