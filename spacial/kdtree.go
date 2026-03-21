package spacial

import "go_surf/models"

type node struct {
	buoy  models.Buoy
	axis  int
	left  *node
	right *node
}

func buildKDTree(buoys []models.Buoy, depth int) *KDnode, error {
	if len(points) == 0 {
		return nil
	}

	// choose axis
	axis := depth % 2
	sort.Slice(buoys, func(i, j int) bool {
		if axis == 0 {
			return buoys[i].Latitude < buoys[i].Latitude
		}
		return buoys[i].Longitude < buoys[j].Longitude
	})

}
