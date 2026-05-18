package meteo

import (
	"context"
	"testing"
)

// takes an input file with bouy ids.
// If pass: returns true
// If fail: returns false and an error
func TestRTBouyGetObservation(t *testing.T) {
	ctx := context.Background()
	bouyId := "46086"

	client := NewClient()
	_, err := client.RTBouy.GetObservation(ctx, bouyId)
	if err != nil {
		t.Errorf("meteo.RTBouy.GetObservation failed: %s", err)
	}
}
