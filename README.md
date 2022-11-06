# Search service

```
git clone repo-name
cd repo-name

go mod init service-name
go get -u github.com/gin-gonic/gin

go get github.com/spf13/viper
go get github.com/jmoiron/sqlx
```

## Set up database running in docker

```
make network
make postgres
make createdb
```

## Database migrations

Install golang-migrate package (https://github.com/golang-migrate/migrate).

To upgrade current db version run:
```
make network
make postgres
make createdb
```

To create new db migration run:
```
migrate create -ext sql -dir db/migration -seq <your-schema-name>
```

## Start service

```
make server
```

Available at:

GET  localhost:8080/v1/listings/:id
GET  localhost:8080/v1/listings?offset=0&limit=10
POST localhost:8080/v1/listings
