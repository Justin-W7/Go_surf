package api

import (
	"database/sql"
	"go_surf/backend/src/models"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// The handler sctruct is needed to provide the get functions with access
// to the data base.
type Handler struct {
	DB *sql.DB
}

type apiCity struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type apiSurfConditions struct {
	ID int `json:"id"`
}

// getCities - returns a json list of all city names and their IDs.
func (h *Handler) getCities(c *gin.Context) {
	cities := []apiCity{}

	rows, err := h.DB.Query("SELECT id, name, state FROM cities ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch cities",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var city apiCity
		if err := rows.Scan(&city.ID, &city.Name, &city.State); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse city data",
			})
			return
		}
		cities = append(cities, city)
	}

	c.JSON(http.StatusOK, cities)
}

// getSurfSpots - takes cityID, returns json of all surf spots and their
// ids for that cityID.
func (h *Handler) getSurfSpots(c *gin.Context) {
	cityIDParam := c.Param("cityID")

	cityID, err := strconv.Atoi(cityIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid cityID",
		})
		return
	}

	surfSpots := []models.StaticSurfSpot{}

	rows, err := h.DB.Query(`
		SELECT id, name, latitude, longitude, city_id, nearest_buoy 
		FROM surfspot 
		WHERE city_id = $1
		`, cityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch static surf spots",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var spot models.StaticSurfSpot
		if err := rows.Scan(
			&spot.ID,
			&spot.Name,
			&spot.Latitude,
			&spot.Longitude,
			&spot.CityID,
			&spot.NearestBuoy,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to parse surfspot data",
			})
			return
		}
		surfSpots = append(surfSpots, spot)
	}

	c.JSON(http.StatusOK, surfSpots)
}

// getSpotConditionsCurrent recieves a surfSpotID and retuns a json response
// of current conditions for that surfSpotID
func (h *Handler) getSpotConditionsCurrent(c *gin.Context) {
	spotIDParam := c.Param("spotID")

	surfSpotID, err := strconv.Atoi(spotIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid spotID",
		})
		return
	}

	var conditions models.CurrentSurfSpotConditions
	err = h.DB.QueryRow(`
		SELECT
			id,
			spot_id, 
			recorded_at, 
			dom_swell_height_m, 
			dom_swell_dir,
			wind_speed_mph,
			wind_direction,
			air_temp_deg_c,
			water_temp_deg_c,
			precipitation,
			cloud_coverage,
			domwp_sec
		FROM current_surf_spot_conditions
		WHERE spot_id = $1
	`, surfSpotID).Scan(
		&conditions.ID,
		&conditions.SpotId,
		&conditions.RecordedAt,
		&conditions.DomSwellHeightM,
		&conditions.DomSwellDir,
		&conditions.WindSpeedMph,
		&conditions.WindDirection,
		&conditions.AirTempDegC,
		&conditions.WaterTempDegC,
		&conditions.Precipitation,
		&conditions.CloudCoverage,
		&conditions.DominantWavePeriodSec,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "no surf conditions found for spot",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, conditions)
}

// StartRouter - creates gin router with default middleware.
// By default it serves on :8080 unless PORT variable is defined.
func StartRouter(db *sql.DB) {
	h := &Handler{DB: db}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/cities", h.getCities)
	router.GET("/surfspots/:cityID", h.getSurfSpots)
	router.GET("/surfforecast/current/:spotID", h.getSpotConditionsCurrent)

	router.Static("/gosurf", "./frontend/src_2")

	router.Run("0.0.0.0:8080")
}
