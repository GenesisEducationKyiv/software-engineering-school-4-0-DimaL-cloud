CREATE SCHEMA "rate-service";

CREATE TABLE subscription (
    id SERIAL PRIMARY KEY,
    email VARCHAR(320) UNIQUE NOT NULL
);