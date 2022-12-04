CREATE TABLE users
(
    id            SERIAL UNIQUE,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP    NOT NULL
);