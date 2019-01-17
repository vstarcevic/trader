package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/vstarcevic/trader/api"
	"github.com/vstarcevic/trader/trader/postgres"
)

const (
	dbhost = "localhost"
	dbport = "5432"
	dbuser = "trader1"
	dbpass = "password"
	dbname = "Trader"

	httpport = "8080"
	cvsFile  = "export.csv"
)

func main() {

	conf := loadConfig()
	db, err := postgres.Connect(conf)
	if err != nil {
		log.Fatal("Error connecting to database. Will exit now.")
		os.Exit(1)
	}
	defer db.Close()

	importer := postgres.NewImport(db, cvsFile)
	importer.ImportContactData()

	contactRepo := postgres.NewContactRepository(db)
	substRepo := postgres.NewSubscriptionRepository(db)
	svc := api.Service{
		ContactRepo:      contactRepo,
		SubscriptionRepo: substRepo,
	}

	startHTTPServer(httpport, svc)

}

func loadConfig() postgres.Config {
	dbConfig := postgres.Config{
		Host: dbhost,
		Port: dbport,
		User: dbuser,
		Pass: dbpass,
		Name: dbname,
	}
	return dbConfig
}

func startHTTPServer(port string, svc api.Service) {
	p := fmt.Sprintf(":%s", port)
	http.ListenAndServe(p, api.MakeHTTPHandler(svc))
}
