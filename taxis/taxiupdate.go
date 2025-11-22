package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// connection details
const (
	Hostname = "localhost"
	Port     = 5433
	Username = "postgres"
	Password = "pwpwpw"
	Database = "chicago"
	Table    = "taxitrips"
)

func main() {

}

type taxiTrip struct {
	trip_id                    string
	taxi_id                    string
	trip_start                 string
	trip_end                   string
	trip_seconds               string
	trip_miles                 string
	pickup_census_tract        string
	dropoff_census_tract       string
	pickup_community_area      string
	dropoff_community_area     string
	fare                       string
	tips                       string
	tolls                      string
	extras                     string
	trip_total                 string
	payment_type               string
	company                    string
	pickup_centroid_latitude   string
	pickup_centroid_longitude  string
	pickup_centroid_location   string
	dropoff_centroid_latitude  string
	dropoff_centroid_longitude string
	dropoff_centroid_location  string
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
