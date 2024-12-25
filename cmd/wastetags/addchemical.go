package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Chemical struct {
	Cas      string `json:"cas"`
	ChemName string `json:"chem_name"`
}

func AddChemical(c *gin.Context) (string, gin.H) {
	if c.Request.Method == http.MethodPost {
		return "add-chemical.html", nil
	}
	return "add-chemical.html", nil
}

func InsertNewChemical(chemicalEntry *Chemical) error {
	chemicalValues := [][]string{
		{chemicalEntry.Cas, chemicalEntry.ChemName},
	}

	aliasValues := [][]string{
		{chemicalEntry.ChemName, chemicalEntry.ChemName},
	}

	addChemicalSql := readEmbeddedFile("query/insert_chemical.sql")
	err := db.InsertData("chemicals", addChemicalSql, chemicalValues)
	if err != nil {
		return err
	}

	addAliasSql := readEmbeddedFile("query/insert_alias.sql")
	err = db.InsertData("alias", addAliasSql, aliasValues)
	return err
}

func checkIfChemicalExists(chemicalEntry *Chemical) (int, gin.H) {
	sqlStatement := readEmbeddedFile("query/get_rows_by_col_value.sql")

	results, err := db.GetRowsByColumnValue("chemicals", sqlStatement, "cas", chemicalEntry.Cas)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message":   err.Error(),
			"cas":       chemicalEntry.Cas,
			"chem_name": chemicalEntry.ChemName,
		}
	}

	if len(results) == 0 {
		err := InsertNewChemical(chemicalEntry)
		if err != nil {
			return http.StatusInternalServerError, gin.H{
				"message":   fmt.Sprintf("Failed to insert new chemical: %s", err.Error()),
				"cas":       chemicalEntry.Cas,
				"chem_name": chemicalEntry.ChemName,
			}
		}

		return http.StatusOK, gin.H{
			"message":   "New chemical successfully inserted",
			"cas":       chemicalEntry.Cas,
			"chem_name": chemicalEntry.ChemName,
		}
	}

	return http.StatusOK, gin.H{
		"message":   "Chemical Already Exists",
		"cas":       results[0]["cas"].(string),
		"chem_name": results[0]["chem_name"].(string),
	}
}

func GetEntryByCas(c *gin.Context) {
	var requestData *Chemical
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   "Invalid input",
			"cas":       "",
			"chem_name": "",
		})
	} else if requestData.Cas == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   "CAS number is required",
			"cas":       "",
			"chem_name": "",
		})
	} else {
		fmt.Printf("Received CAS: %s, ChemName: %s\n", requestData.Cas, requestData.ChemName)
		c.JSON(checkIfChemicalExists(requestData))
	}
}
