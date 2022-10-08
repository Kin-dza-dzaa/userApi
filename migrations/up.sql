CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    id                UUID                                                                  NOT NULL,
    user_name         TEXT                                                                  NOT NULL CHECK(user_name != ''),
    email             TEXT                                                                  NOT NULL CHECK(email != ''),
    password          TEXT                                                                  NOT NULL CHECK(password != ''),
    registration_date TIMESTAMP                                                             NOT NULL, 
    refresh_token     TEXT,
    expiration_time   TIMESTAMP,  
    verification_code TEXT                                                                  NOT NULL,
    verified          BOOL                                                                  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE(user_name),
    UNIQUE (email)
);
