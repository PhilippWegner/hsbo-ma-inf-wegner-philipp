package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/restful-broker-service/docs"
)

const webPort = "80"

type Config struct{}

// @title Broker Service represented as RESTful API
// @version 1.0
// @description This is a RESTful API for the Broker Service as a single point of entry for all api caLLs.
// @host localhost:8080
// @BasePath /api
func main() {
	app := Config{}

	server := &http.Server{
		Addr:         ":" + webPort,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
