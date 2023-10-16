package storage

const queryAdd = `INSERT INTO weather(location_name,latitude,longitude,timestamp,temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m)
					values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

const queryRecordExists = `SELECT 1
							FROM weather
							WHERE location_name = $1
							AND timestamp = $2`

const queryUpdateByLocationAndTime = `UPDATE weather
										SET latitude=$1,
										longitude=$2,
										temperature_2m=$3,
										relativehumidity_2m=$4,
										precipitation_probability=$5,
										visibility=$6,
										windspeed_10m=$7
										WHERE location_name = $1
										AND timestamp = $2`

const queryGetLatest = `SELECT location_name,latitude,longitude,timestamp,temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m
						FROM weather
						WHERE location_name = $1
						ORDER BY timestamp DESC
						LIMIT 1`

const queryGetPeriod = `SELECT location_name,latitude,longitude,timestamp,temperature_2m,relativehumidity_2m,precipitation_probability,visibility,windspeed_10m
						FROM weather
						WHERE location_name = $1
						AND timestamp BETWEEN $2 and $3
						ORDER BY timestamp DESC`