package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"lightning-poll/db"
	"lightning-poll/lnd"
)

var router *gin.Engine

func main() {
	flag.Parse()

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("/Users/carla/personal/src/lightning-poll/templates/*")

	dbc, err := db.Connect()
	if err != nil {
		log.Fatalf("could not connect to DB: %v", err)
	}

	lndCl, err := lnd.New()
	if err != nil {
		log.Fatalf("could not connect to LND: %v", err)
	}

	// Initialize the routes
	initializeRoutes(Env{db: dbc, lnd: lndCl})

	// Start serving the application
	router.Run()

}
