package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/config"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/routes"
)

func main() {

	config.InitLog()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to load config")
		return
	}

	database, err := db.Connnect()
	if err != nil {
		log.Fatal("failed to connect db")
		return
	}

	r := mux.NewRouter()
	routes.Route(database, r)

	fmt.Printf("starting server on :%s\n", conf[config.APP_PORT])
	log.Println(fmt.Sprintf("starting server on :%s\n", conf[config.APP_PORT]))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf[config.APP_PORT]), r))

}
