package main

import (
	"filecloud/common/dl"
	"filecloud/common/dl/models"
	"log"
)

func main() {
	log.Println("nodeorc started")

	var cfg models.OrcAppCfg
	err := dl.ReadConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := NewNodeOrcHttpServer(&cfg)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
