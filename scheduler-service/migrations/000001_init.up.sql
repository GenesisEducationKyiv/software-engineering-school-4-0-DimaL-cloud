CREATE SCHEMA "mail-service";

CREATE TABLE event
(
    id        SERIAL PRIMARY KEY,
    type      VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP    NOT NULL,
    body      TEXT
);