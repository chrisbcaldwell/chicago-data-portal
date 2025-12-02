const Pool = require('pg').Pool
const pool = new Pool({
  user: 'Postgres',
  host: 'localhost',
  database: 'chicago',
  password: 'italianbeef',
  port: 5434,
})

// get all taxi trips
const getTaxiTrips = (request, response) => {
    let q = `
    SELECT
        trip_id,
        trip_start_timestamp,
        trip_end_timestamp,
        pickup_community_area,
        dropoff_community_area
    FROM taxitrips
    ORDER BY trip_start_timestamp ASC`
    pool.query(q, (error, results) => {
        if (error) {
            throw error
        }
        response.status(200).json(results.rows)
    })
}

// get all taxi trips originating in a community area
const getTaxiTripsByPickup = (request, response) => {
    const pickup = parseInt(request.params.pickup)
    let q = `
    SELECT
        trip_id,
        trip_start_timestamp,
        trip_end_timestamp,
        pickup_community_area,
        dropoff_community_area
    FROM taxitrips
    WHERE pickup_community_area = $1
    ORDER BY trip_start_timestamp ASC`
    pool.query(q, [pickup], (error, results) => {
        if (error) {
            throw error
        }
        response.status(200).json(results.rows)
    })
}

// get all taxi trips dropping off in a community area
const getTaxiTripsByDropoff = (request, response) => {
    const dropoff = parseInt(request.params.dropoff)
    let q = `
    SELECT
        trip_id,
        trip_start_timestamp,
        trip_end_timestamp,
        pickup_community_area,
        dropoff_community_area
    FROM taxitrips
    WHERE dropoff_community_area = $1
    ORDER BY trip_start_timestamp ASC`
    pool.query(q, [dropoff], (error, results) => {
        if (error) {
            throw error
        }
        response.status(200).json(results.rows)
    })
}

module.exports = {
    getTaxiTrips,
    getTaxiTripsByPickup,
    getTaxiTripsByDropoff
}