package main

import (
	"fmt"
	"net/http"
	"strings"

	"sort"

	"github.com/gin-gonic/gin"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type ChemicalSearch struct {
	Name  string
	Score int
}

func fuzzySearch(query string, dataset []string, topN int) []ChemicalSearch {
	results := []ChemicalSearch{}

	for _, name := range dataset {
		if name != "" {
			score := levenshtein.DistanceForStrings([]rune(query), []rune(name), levenshtein.DefaultOptions)
			results = append(results, ChemicalSearch{Name: name, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	if len(results) > topN {
		results = results[:topN]
	}

	return results
}

func SearchChemicalNames(query string) []ChemicalSearch {
	if query == "" {
		initialDataSet := []string{
			"Water",
			"Tetramethylammonium Hydroxide",
			"N-Methyl-2-pyrrolidone",
			"Isopropanol",
			"1-Methoxy-2-propyl acetate",
			"Sodium Hydroxide",
			"Proprietary",
			"Methyl Isobutyl Ketone (MIBK)",
			"Methanol",
			"Poly[(o-cresyl glycidyl ether)-co-formaldehyde]",
			"Phosphoric Acid",
			"Nitric acid",
			"Acetone",
			"Acetic acid",
		}
		results := []ChemicalSearch{}
		for _, chemical := range initialDataSet {
			results = append(results, ChemicalSearch{Name: chemical, Score: 0})
		}
		return results
	}
	topN := 20
	sqlStatement := readEmbeddedFile("query/get_all.sql")
	chemicals, _ := db.GetAll("chemicals", sqlStatement)
	dataSet := make([]string, len(chemicals))
	for i, chemical := range chemicals {
		dataSet[i] = fmt.Sprint(chemical["chem_name"])
	}

	results := fuzzySearch(query, dataSet, topN)
	return results
}

func SearchForChemical(c *gin.Context) {
	type Query struct {
		Query string `json:"query"`
	}

	var query Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Name":  "",
			"Score": 0,
		})
	} else if query.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Name":  "",
			"Score": 0,
		})
	} else {
		c.JSON(http.StatusOK, SearchChemicalNames(query.Query))
	}
}

func AddMixture(c *gin.Context) (string, gin.H) {
	if c.Request.Method == http.MethodPost {
		sqlStatement := readEmbeddedFile("query/insert_mixture.sql")
		casNumber := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(c.PostForm("cas1")), strings.TrimSpace(c.PostForm("cas2")), strings.TrimSpace(c.PostForm("cas3")))
		values := [][]string{
			{casNumber, strings.TrimSpace(c.PostForm("chemical-name"))},
		}
		err := db.InsertData("mixtures", sqlStatement, values)
		if err != nil {
			errorMessage := fmt.Sprintln("Error adding chemical:", err)
			return "error.html", gin.H{
				"message": errorMessage,
			}
		}
	}

	initialData := SearchChemicalNames("")
	return "add-mixture.html", gin.H{"Components": initialData}
}
