CREATE TABLE movie
(
    id              SERIAL UNIQUE,
    name            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    production_year INTEGER,
    genre           VARCHAR(20)  NOT NULL,
    actors          VARCHAR(255),
    poster          VARCHAR(255)
);