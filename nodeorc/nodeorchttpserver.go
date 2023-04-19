package main

import (
	"filecloud/common/comm"
	"filecloud/common/dl/models"
	"log"
	"net/http"
)

type NodeOrcHttpServer struct {
	config *models.OrcAppCfg
	no     *NodeOrchestrator
}

func NewNodeOrcHttpServer(config *models.OrcAppCfg) *NodeOrcHttpServer {
	return &NodeOrcHttpServer{
		config: config,
		no:     NewNodeOrchestrator(config),
	}
}

func (hs *NodeOrcHttpServer) ListenAndServe() error {
	var s = &http.Server{
		Addr: hs.config.HostAndPort,
	}

	http.HandleFunc("/", hs.handler)

	err := s.ListenAndServe()
	return err
}

func (hs *NodeOrcHttpServer) handler(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPut {
		err := hs.putHandler(response, request)
		if err != nil {
			comm.SetResponse(response, http.StatusInternalServerError, nil, err)
			return
		}
	}

	if request.Method == http.MethodGet {
		err := hs.getHandler(response, request)
		if err != nil {
			comm.SetResponse(response, http.StatusInternalServerError, nil, err)
			return
		}
	}
}

func (hs *NodeOrcHttpServer) putHandler(response http.ResponseWriter, request *http.Request) error {
	fileName := request.Header.Get(models.FileNameHeader)
	log.Printf("put: получение и отправка файла %+v", fileName)
	err := hs.no.SendReceipts(fileName)
	if err != nil {
		return err
	}

	err = hs.no.SendBlocks(fileName, request)
	if err != nil {
		return err
	}

	return nil
}

func (hs *NodeOrcHttpServer) getHandler(response http.ResponseWriter, request *http.Request) error {
	fileName := request.Header.Get(models.FileNameHeader)
	log.Printf("get: получение и отправка файла %+v", fileName)

	return hs.no.GetFile(fileName, response)
}
