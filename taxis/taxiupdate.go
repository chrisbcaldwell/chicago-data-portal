package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// connection details
// URL is for the taxi trips 2013-23 data set
const (
	Hostname      = "localhost"
	Port          = 5433
	Username      = "postgres"
	Password      = "root"
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
		log.Fatal(err)
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

type record struct {
	RowID       string `json:"trip_id"`
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
	data, err := getJSON(u)
	if err != nil {
		return err
	}
	err = updateDB(db, data)
	return err
}

func getJSON(u string) ([]record, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: http status code %d", resp.StatusCode)
	}
	var records []record
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, fmt.Errorf("unable to decode JSON: %w", err)
	}
	return records, nil
}

func updateDB(db *sql.DB, d []record) error {
	var ctx context.Context
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var x record // nil record from which to get the field names
	fieldCount := reflect.ValueOf(x).NumField()
	placeholders := ""
	for i := 1; i < fieldCount; i++ {
		placeholders = placeholders + "$" + strconv.Itoa(i) + ", "
	}
	placeholders = placeholders + "$" + strconv.Itoa(fieldCount)
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", Table, placeholders)
	for _, r := range d {
		_, err = tx.Exec(query,
			r.RowID,
			r.TaxiID,
			r.TripStart,
			r.TripEnd,
			r.TrimSeconds,
			r.TripMiles,
			r.PTract,
			r.DTract,
			r.PCA,
			r.DCA,
			r.Fare,
			r.Tips,
			r.Tolls,
			r.Extras,
			r.Total,
			r.Payment,
			r.Company,
			r.PLat,
			r.PLong,
			r.PLoc,
			r.DLat,
			r.DLong,
			r.DLoc)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated records committed. ", len(d), "records added.")
	return nil
}
