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
