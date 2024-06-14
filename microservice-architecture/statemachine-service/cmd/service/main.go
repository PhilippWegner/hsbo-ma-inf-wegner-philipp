package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/statemachine-service/data"
)

var (
	restful_api = "http://restful-api:80/api/handle"
	machines    = []string{"presse_11"}
)

type Config struct {
	ApiRepository data.Repository
	Machines      []string
}

func main() {
	log.Println("Starting application")
	app := Config{}
	app.ApiRepository = data.NewApiRepository(restful_api)
	app.Machines = machines

	log.Println("Starting looprunner")
	for {
		start := time.Now()
		app.looprunner()
		duration := time.Since(start)
		app.ApiRepository.CreateLog(data.Log{Name: "looprunner", Data: fmt.Sprintf("looprunner took %v", duration)})
	}
}

func (c *Config) looprunner() {
	var wg sync.WaitGroup
	wg.Add(len(c.Machines))
	for _, machine := range c.Machines {
		go func(machine string) {
			defer wg.Done()
			err := c.calculate(machine)
			if err != nil {
				c.ApiRepository.CreateLog(data.Log{Name: "error", Data: fmt.Sprintf("calculate(%v) failed: %v", machine, err)})
			}
		}(machine)
	}
	// Wait for all goroutines to finish
	wg.Wait()
}
