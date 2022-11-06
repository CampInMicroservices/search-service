CREATE TABLE "listings" (
    "id"            BIGSERIAL PRIMARY KEY,
    "name"          VARCHAR(40) NOT NULL,
    "created_at"    TIMESTAMP NOT NULL DEFAULT(now())
);