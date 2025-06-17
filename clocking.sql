CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    firstname VARCHAR(100) NOT NULL,
    middlename VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    latitude_dms VARCHAR(50) NOT NULL,
    longitude_dms VARCHAR(50) NOT NULL,
    latitude DECIMAL(9,6),     -- e.g., 14.599512
    longitude DECIMAL(9,6),    -- e.g., 120.984222
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE events (
    eventID SERIAL PRIMARY KEY,
    eventName VARCHAR(255) NOT NULL,
    releaseTime TIMESTAMP NOT NULL DEFAULT date_trunc('minute', now()),
    releaseLat DOUBLE PRECISION NOT NULL,         -- Decimal degrees (positive for N, negative for S)
    releaseLng DOUBLE PRECISION NOT NULL,         -- Decimal degrees (positive for E, negative for W)
    releaseLatDMS VARCHAR(50),                    -- Optional: store original DMS format (e.g. "12°36′15.47″ N")
    releaseLngDMS VARCHAR(50)                     -- Optional: store original DMS format (e.g. "123°45′30.25″ E")
);

CREATE TABLE arrivals (
    id SERIAL PRIMARY KEY,
    userID INT REFERENCES users(id),
    eventID INT REFERENCES events(eventID),
    arrivedAt TIMESTAMP NOT NULL,
    speed DECIMAL(10,3) -- m/min
);


select * from users
drop table users,events,arrivals cascade



