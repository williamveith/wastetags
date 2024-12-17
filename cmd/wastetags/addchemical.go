package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func addChemical(c *gin.Context) (string, gin.H) {
	if c.Request.Method == http.MethodPost {
		sqlStatement := readSql("query/insert_chemical.sql")
		casNumber := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(c.PostForm("cas1")), strings.TrimSpace(c.PostForm("cas2")), strings.TrimSpace(c.PostForm("cas3")))
		values := [][]string{
			{casNumber, strings.TrimSpace(c.PostForm("chemical-name"))},
		}
		err := db.InsertData("chemicals", sqlStatement, values)
		if err != nil {
			errorMessage := fmt.Sprintln("Error adding chemical:", err)
			return "error.html", gin.H{
				"message": errorMessage,
			}
		}
	}
	return "addchemical.html", nil
}
