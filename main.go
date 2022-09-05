package main

import (
	"log"

	"github.com/namsral/flag"

	"github.com/rehacktive/qaservice/service"
)

const (
	envHost           = "DB_HOST"
	defaultMongoLocal = "mongodb://localhost:27017"
)

func main() {
	var host string

	log.Println("service starting...")

	flag.StringVar(&host, envHost, defaultMongoLocal, "host:port for postgres")

	flag.Parse()

	dbInstance, err := service.InitDb(host)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database.")

	srv := service.InitService(dbInstance)

	log.Println("service started")
	srv.Start()
}
