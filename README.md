# Search service

## Installation (locally)


1. Install make (`choco install make` for Win)
2. Install golang-migrate package (https://github.com/golang-migrate/migrate)
3. ...

### Run the following commands

```
# Create docker network
make network

# Start database container and create main db
make postgres
make createdb

# Run db migrations
make migrateup

# Start the service
make server
```

Service is avalilable at http://localhost:8080

```
GET  localhost:8080/v1/listings/:id
GET  localhost:8080/v1/listings?offset=0&limit=10
POST localhost:8080/v1/listings
```
## Database migrations

To update schema of the database, run the following command:

```
migrate create -ext sql -dir db/migration -seq <your-schema-name>
```

New file should be created in db/migrations.

## How the service was created

```
git clone repo-name
cd repo-name

go mod init service-name
go get -u github.com/gin-gonic/gin

go get github.com/spf13/viper
go get github.com/jmoiron/sqlx
```
