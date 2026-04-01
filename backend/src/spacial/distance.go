package spacial

import (
	"database/sql"
	"fmt"
	"log"
	"math"
)

// NearestBuoy finds the nearest buoy to a surf spot.
// Function runs once on database build and surfspot/buoy updates.
func NearestBuoy(lat, lon float64, db *sql.DB) int {
	// fetch db data
	rows, err := db.Query("SELECT id, latitude, longitude FROM buoys")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id int
	var bLat float64
	var bLon float64
	var nearestBuoyId int
	var current float64
	nearest := math.MaxFloat64

	// for each row in buoy table
	for rows.Next() {

		err = rows.Scan(&id, &bLat, &bLon)
		if err != nil {
			log.Fatal(err)
		}

		// find distance between buoy and spot
		d := haversine(lat, lon, bLat, bLon)
		current = d

		if current < nearest {
			nearest = current
			nearestBuoyId = id
		}
		fmt.Println("buoy_id: ", id, "   ", nearest)
	}

	return nearestBuoyId
}

// haversine function finds the distance bewteen two
// geo coordinates on earth.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earths radius, km

	const degToRad = math.Pi / 180.0

	phi1 := lat1 * degToRad
	phi2 := lat2 + degToRad
	dphi := (lat2 - lat1) * degToRad
	dlambda := (lon2 - lon1) * degToRad

	a := math.Sin(dphi/2)*math.Sin(dphi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(dlambda/2)*math.Sin(dlambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
