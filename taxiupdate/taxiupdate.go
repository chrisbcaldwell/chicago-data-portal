package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// connection details
// URL is for the taxi trips 2013-23 data set
const (
	Hostname      = "localhost"
	Port          = 5434
	Username      = "postgres"
	Password      = "italianbeef"
	Database      = "chicago"
	Table         = "taxitrips"
	LastUpdateCol = "trip_start_timestamp"
	URL           = "https://data.cityofchicago.org/resource/wrvz-psew.json"
)

var URLs = []string{
	"https://data.cityofchicago.org/resource/wrvz-psew.json", // 2013-2024
	"https://data.cityofchicago.org/resource/ajtu-isnz.json"} // 2024+
// URL that works for getting row count
// https://data.cityofchicago.org/resource/wrvz-psew.json?$select=count(*)

func main() {
	// connect to the db
	db, err := openConnection()
	if err != nil {
		fmt.Println("Error connecting to", Database)
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("openConnection() successful, connected to", Database)

	// find the newest record already in the database
	updated := lastUpdate(db)
	if updated == "Error" {
		fmt.Println("Error connecting to database", Database)
		panic("Error connecting to database")
	}

	fmt.Println("last update to", Table, updated)

	// get the JSON of all records newer than what's already loaded in the database
	var requests []string
	for _, u := range URLs {
		requests = append(requests, u+"?$where="+LastUpdateCol+">%27"+updated+"%27") // "%27" = "'"
	}

	err = readAndUpdate(requests, db)
	if err != nil {
		log.Fatal(err)
	}
}

type record struct {
	RowID       string           `json:"trip_id"`
	TaxiID      string           `json:"taxi_id"`
	TripStart   string           `json:"trip_start_timestamp"`
	TripEnd     string           `json:"trip_end_timestamp"`
	TripSeconds string           `json:"trip_seconds"`
	TripMiles   string           `json:"trip_miles"`
	PTract      string           `json:"pickup_census_tract"`
	DTract      string           `json:"dropoff_census_tract"`
	PCA         string           `json:"pickup_community_area"`
	DCA         string           `json:"dropoff_community_area"`
	Fare        string           `json:"fare"`
	Tips        string           `json:"tips"`
	Tolls       string           `json:"tolls"`
	Extras      string           `json:"extras"`
	Total       string           `json:"trip_total"`
	Payment     string           `json:"payment_type"`
	Company     string           `json:"company"`
	PLat        string           `json:"pickup_centroid_latitude"`
	PLong       string           `json:"pickup_centroid_longitude"`
	PLoc        *json.RawMessage `json:"pickup_centroid_location"`
	DLat        string           `json:"dropoff_centroid_latitude"`
	DLong       string           `json:"dropoff_centroid_longitude"`
	DLoc        *json.RawMessage `json:"dropoff_centroid_location"`
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
	var u sql.NullString
	err := db.QueryRow(query).Scan(&u)
	if err != nil {
		fmt.Println("Error in query", query)
		log.Fatal(err)
		return "Error"
	}
	var updated string
	if u.Valid {
		updated = u.String
	} else {
		updated = defaultDate
	}
	// get into the SODA format
	updated = strings.Replace(updated, " ", "T", 1)
	return updated
}

func readAndUpdate(urls []string, db *sql.DB) error {
	for _, u := range urls {
		fmt.Println("getting data from", u)
		data, err := getJSON(u)
		if err != nil {
			return err
		}
		fmt.Println("getJSON() function returned successfully for", u)
		err = updateDB(db, data)
		if err != nil {
			fmt.Println("error updating database from", u)
			return err
		}
	}
	return nil
}

func getNumRows(u string) int {
	q := u + "?$select=count(*)"
	r, err := http.Get(q)
	if err != nil {
		return -1
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return -1
	}
	dec := json.NewDecoder(r.Body)
	var m map[string]string
	err = dec.Decode(&m)
	if err != nil {
		return -1
	}
	i, _ := strconv.Atoi(m["count"])
	return i
}

func getJSON(u string) ([]record, error) {
	fmt.Println(getNumRows(u))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: http status code %d", resp.StatusCode)
	}
	fmt.Println("Data download successful")
	var records []record
	i := 0
	dec := json.NewDecoder(resp.Body)
	// read open bracket
	t, err := dec.Token()
	if err != nil {
		fmt.Println("Error reading opening JSON token")
		log.Fatal(err)
	}
	fmt.Println("Opening token read successfully:", t)
	// while dec contains values
	for dec.More() {
		i++
		var r record
		err := dec.Decode(&r)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
		fmt.Printf("\rJSON record #%d decoded", i)
	}
	fmt.Println(i, "\nJSON records returned")
	// saving JSON for logging purposes:
	data, _ := json.Marshal(records)
	os.WriteFile("data.json", data, 0666)
	return records, nil
}

func updateDB(db *sql.DB, d []record) error {
	var x record // nil record from which to get the field names
	fieldCount := reflect.ValueOf(x).NumField()
	placeholders := ""
	for i := 1; i < fieldCount; i++ {
		placeholders = placeholders + "$" + strconv.Itoa(i) + ", "
	}
	placeholders = placeholders + "$" + strconv.Itoa(fieldCount)
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", Table, placeholders)
	fmt.Println("Query:")
	fmt.Println(query)
	for i, r := range d {
		_, err := db.Exec(query,
			// some entries will be "" strings for numeric values, PostgreSQL hates that
			// they will be transformed to "0"
			r.RowID,
			r.TaxiID,
			r.TripStart,
			r.TripEnd,
			handleBlank(r.TripSeconds),
			handleBlank(r.TripMiles),
			r.PTract,
			r.DTract,
			handleBlank(r.PCA),
			handleBlank(r.DCA),
			handleBlank(r.Fare),
			handleBlank(r.Tips),
			handleBlank(r.Tolls),
			handleBlank(r.Extras),
			handleBlank(r.Total),
			r.Payment,
			r.Company,
			handleBlank(r.PLat),
			handleBlank(r.PLong),
			r.PLoc,
			handleBlank(r.DLat),
			handleBlank(r.DLong),
			r.DLoc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\rRecord #%d inserted into %s", i+1, Table)
	}
	return nil
}

func handleBlank(s string) string {
	if s == "" {
		return "0"
	}
	return s
}
