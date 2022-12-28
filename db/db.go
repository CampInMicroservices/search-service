package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Store struct {
	db *sqlx.DB
}

func Connect(driver, source string) (*Store, error) {

	db, err := sqlx.Connect(driver, source)
	if err != nil {
		return nil, err
	}

	store := &Store{
		db: db,
	}

	log.Println("Connected to database!")
	return store, err
}

func (store *Store) PingDB() error {
	return store.db.Ping()
}

func (store *Store) Close() error {
	return store.db.Close()
}

func RunDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("Cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrate up", err)
	}

	log.Println("Db migrated successfully")
}
