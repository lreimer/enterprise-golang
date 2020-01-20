package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// configuration for static files and templates
	engine.LoadHTMLFiles("templates/index.html")
	engine.StaticFile("/favicon.ico", "favicon.ico")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "A Hitchhiker's Guide to Enterprise Microservices with Go",
		})
	})

	engine.GET("/api/gins", allSpirits)            // get list of gins
	engine.POST("/api/gins", createSpirit)         // create new gin
	engine.GET("/api/gins/:asin", getSpirit)       // get gin by ASIN
	engine.PUT("/api/gins/:asin", updateSpirit)    // update existing gin
	engine.DELETE("/api/gins/:asin", deleteSpirit) // delete book

	// run server on PORT
	engine.Run(port())
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func allSpirits(c *gin.Context) {
	c.JSON(http.StatusOK, AllSpirits())
}

func getSpirit(c *gin.Context) {
	asin := c.Params.ByName("asin")
	gin, found := GetSpirit(asin)
	if found {
		c.JSON(http.StatusOK, gin)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createSpirit(c *gin.Context) {
	var gin Spirit
	if c.BindJSON(&gin) == nil {
		asin, created := CreateSpirit(gin)
		if created {
			c.Header("Location", "/api/gins/"+asin)
			c.Status(http.StatusCreated)
		} else {
			c.Status(http.StatusConflict)
		}
	}
}

func updateSpirit(c *gin.Context) {
	asin := c.Params.ByName("asin")

	var gin Spirit
	if c.BindJSON(&gin) == nil {
		exists := UpdateSpirit(asin, gin)
		if exists {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}

func deleteSpirit(c *gin.Context) {
	asin := c.Params.ByName("asin")
	DeleteSpirit(asin)
	c.Status(http.StatusOK)
}

// Spirit type with ASIN, Name, Country and Alcohol
type Spirit struct {
	ASIN    string `json:"asin"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Alcohol int    `json:"alcohol"`
}

var spirits = map[string]Spirit{
	"B00N3P44OM": Spirit{ASIN: "B00N3P44OM", Name: "Windspiel Premium Dry Gin", Country: "Germany", Alcohol: 47},
	"B00IE97JCQ": Spirit{ASIN: "B00IE97JCQ", Name: "Granit Bavarian Gin", Country: "Germany", Alcohol: 42},
	"B00A0DF494": Spirit{ASIN: "B00A0DF494", Name: "Gin Mare", Country: "Spain", Alcohol: 43},
}

// AllSpirits returns a slice of all Spirits
func AllSpirits() []Spirit {
	values := make([]Spirit, len(spirits))
	idx := 0
	for _, spirit := range spirits {
		values[idx] = spirit
		idx++
	}
	return values
}

// GetSpirit returns the spirit for a given ASIN
func GetSpirit(asin string) (Spirit, bool) {
	spirit, found := spirits[asin]
	return spirit, found
}

// CreateSpirit creates a new Spirit if it does not exist
func CreateSpirit(spirit Spirit) (string, bool) {
	_, exists := spirits[spirit.ASIN]
	if exists {
		return "", false
	}
	spirits[spirit.ASIN] = spirit
	return spirit.ASIN, true
}

// UpdateSpirit updates an existing spirit
func UpdateSpirit(asin string, spirit Spirit) bool {
	_, exists := spirits[asin]
	if exists {
		spirits[asin] = spirit
	}
	return exists
}

// DeleteSpirit removes a spirit from the map by ASIN key
func DeleteSpirit(asin string) {
	delete(spirits, asin)
}
