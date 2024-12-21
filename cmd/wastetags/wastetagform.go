package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeWasteTagForm(c *gin.Context) (string, gin.H) {
	genericErrorMessage := gin.H{"message": "Internal Server Error"}

	if c.Request.Method != http.MethodGet {
		return "error.html", genericErrorMessage
	}

	sqlStatement := readEmbeddedFile("query/get_distinct_col_values.sql")
	components, err := db.GetColumnValues("mixtures", sqlStatement, "chem_name")
	if err != nil {
		fmt.Println("Error retrieving components:", err)
		return "error.html", genericErrorMessage
	}
	cleanComponents := make([]map[string]string, len(components))
	for i, comp := range components {
		if chemName, ok := comp["chem_name"].(string); ok {
			cleanComponents[i] = map[string]string{"chem_name": chemName}
		}
	}

	locations, err := db.GetColumnValues("locations", sqlStatement, "location")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		return "error.html", genericErrorMessage
	}

	cleanLocations := make([]map[string]string, len(locations))
	for i, comp := range locations {
		if location, ok := comp["location"].(string); ok {
			cleanLocations[i] = map[string]string{"location": location}
		}
	}

	units, err := db.GetColumnValues("units", sqlStatement, "full_name")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		return "error.html", genericErrorMessage
	}

	cleanUnits := make([]map[string]string, len(units))
	for i, comp := range units {
		if unit, ok := comp["full_name"].(string); ok {
			cleanUnits[i] = map[string]string{"full_name": unit}
		}
	}

	containers, err := db.GetColumnValues("containers", sqlStatement, "full_name")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		return "error.html", genericErrorMessage
	}

	cleanContainers := make([]map[string]string, len(containers))
	for i, comp := range containers {
		if container, ok := comp["full_name"].(string); ok {
			cleanContainers[i] = map[string]string{"full_name": container}
		}
	}

	states, err := db.GetColumnValues("states", sqlStatement, "state")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		return "error.html", genericErrorMessage
	}

	cleanStates := make([]map[string]string, len(states))
	for i, comp := range states {
		if state, ok := comp["state"].(string); ok {
			cleanStates[i] = map[string]string{"state": state}
		}
	}

	pageData := gin.H{
		"Components": cleanComponents,
		"Locations":  cleanLocations,
		"Units":      cleanUnits,
		"Containers": cleanContainers,
		"States":     cleanStates,
	}

	return "waste-tag-form.html", pageData
}
