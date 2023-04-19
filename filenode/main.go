package main

import (
	"filecloud/common/dl"
	"filecloud/common/dl/models"
	"log"
)

func main() {
	log.Println("filenode started")
	var cfg models.NodeAppCfg
	err := dl.ReadConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := NewNodeHttpServer(&cfg)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
