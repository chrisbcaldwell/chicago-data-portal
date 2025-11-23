package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

// connection details
const (
	Hostname      = "localhost"
	Port          = 5433
	Username      = "postgres"
	Password      = "pwpwpw"
	Database      = "chicago"
	Table         = "taxitrips"
	LastUpdateCol = "trip_start_timestamp"
	URL           = "https://data.cityofchicago.org/resource/wrvz-psew.json"
)

func main() {
	// connect to the db
	db, err := openConnection()
	if err != nil {
		fmt.Println("Error connecting to", Database)
		fmt.Println(err)
	}
	defer db.Close()

	// find the newest record already in the database
	updated := lastUpdate(db)

	// get the JSON of all records newer than what's already loaded in the database
	request := URL + "?$where=" + LastUpdateCol + ">%27" + updated + "%27" // "%27" = "'"

	err = readAndUpdate(request, db)
	if err != nil {
		log.Fatal(err)
	}
}

type taxiTrip struct {
	TripID      string `json:"trip_id"`
	TaxiID      string `json:"taxi_id"`
	TripStart   string `json:"trip_start"`
	TripEnd     string `json:"trip_end"`
	TrimSeconds string `json:"trip_seconds"`
	TripMiles   string `json:"trip_miles"`
	PTract      string `json:"pickup_census_tract"`
	DTract      string `json:"dropoff_census_tract"`
	PCA         string `json:"pickup_community_area"`
	DCA         string `json:"dropoff_community_area"`
	Fare        string `json:"fare"`
	Tips        string `json:"tips"`
	Tolls       string `json:"tolls"`
	Extras      string `json:"extras"`
	Total       string `json:"trip_total"`
	Payment     string `json:"payment_type"`
	Company     string `json:"company"`
	PLat        string `json:"pickup_centroid_latitude"`
	PLong       string `json:"pickup_centroid_longitude"`
	PLoc        string `json:"pickup_centroid_location"`
	DLat        string `json:"dropoff_centroid_latitude"`
	DLong       string `json:"dropoff_centroid_longitude"`
	DLoc        string `json:"dropoff_centroid_location"`
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

func lastUpdate(db *sql.DB) string {
	query := "SELECT MAX(" + LastUpdateCol + ") FROM " + Table + ";"
	defaultDate := "1978-04-09 06:00:00"
	var updated string
	err := db.QueryRow(query).Scan(&updated)
	if err != nil {
		fmt.Println(err)
		return "Error"
	}
	// this condition might end up being nil instead
	if updated == "NULL" {
		updated = defaultDate
	}
	// get into the SODA format
	updated = strings.Replace(updated, " ", "T", 1)
	return updated
}

func readAndUpdate(u string, db *sql.DB) error {
	var err error

	// read the JSON line by line

	// update the database with each line's record

	return err
}
