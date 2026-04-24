package models

type ObservationStationCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Geometry 		FeatureGeometry		`json:"geometry"`
	Properties 		FeatureProperties	`json:"properties"`
}

type FeatureGeometry struct {
	Coordinates []float64 `json:"coordinates"` // [ longitude, latitude ]
}

type FeatureProperties struct {
	StationIdentifier string `json:"stationIdentifier"`
}