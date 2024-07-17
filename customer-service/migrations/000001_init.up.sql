CREATE TABLE customer
(
    id    SERIAL PRIMARY KEY,
    email VARCHAR(320) UNIQUE NOT NULL,
    subscription_id INT NOT NULL
);