CREATE TABLE weather(
    id SERIAL PRIMARY KEY,
    location_name VARCHAR(128) NOT NULL,
    latitude NUMERIC(6, 4) NOT NULL,
    longitude NUMERIC(7, 4) NOT NULL,
    timestamp VARCHAR(64) NOT NULL,
    temperature_2m NUMERIC(3, 1),
    relativehumidity_2m INTEGER,
    precipitation_probability INTEGER,
    visibility INTEGER,
    windspeed_10m NUMERIC(3, 1)
);

