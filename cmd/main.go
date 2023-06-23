package main

import (
	"log"
	"yafa/internal/pkg/db"
	"yafa/internal/pkg/router"
)

func main() {
	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	router := router.New(db)
	log.Fatal(router.Start(":5000"))
}
