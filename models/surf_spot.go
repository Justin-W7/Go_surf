package models

type SurfSpot struct {
	SpotID        int       `json:"_id"`
	SpotName      string    `json:"spot_name"`
	CoastOrder    string    `json:"coast_order"`
	Coordinates   []float64 `json:"coordinates"`
	CountyID      int       `json:"county_id"`
	SpotIDChar    string    `json:"spot_id_char"`
	StreetAddress string    `json:"street_address"`
}
