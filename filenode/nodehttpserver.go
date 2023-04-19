package main

import (
	"filecloud/common/comm"
	"filecloud/common/dl/models"
	"filecloud/filenode/storage"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NodeOrcHttpServer struct {
	config         *models.NodeAppCfg
	receiptStorage *storage.ReceiptStorage
	fileStorage    *storage.FileStorage
}

func NewNodeHttpServer(config *models.NodeAppCfg) *NodeOrcHttpServer {
	return &NodeOrcHttpServer{
		config:         config,
		receiptStorage: storage.NewReceiptStorage(config.ReceiptFolder),
		fileStorage:    storage.NewFileStorage(config.FilePartsFolder),
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
	step := models.UploadStep(request.Header.Get(models.UploadStepHeader))
	log.Printf("put: получение и отправка файла %+v (%+v)", fileName, step)
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	switch step {
	case models.UploadStepReceipt:
		return hs.receiptStorage.PutReceipt(fileName, data)
	case models.UploadStepStart:
		return hs.fileStorage.CreateFile(fileName, data)
	case models.UploadStepAppend:
		return hs.fileStorage.AppendFile(fileName, data)
	default:
		return fmt.Errorf("ошибка определения шага сохранения файла")
	}
}

func (hs *NodeOrcHttpServer) getHandler(response http.ResponseWriter, request *http.Request) error {
	fileName := request.Header.Get(models.FileNameHeader)
	step := models.UploadStep(request.Header.Get(models.UploadStepHeader))
	log.Printf("get: получение и отправка файла %+v (%+v)", fileName, step)

	switch step {
	case models.DownloadReceipt:
		data, exists, err := hs.receiptStorage.GetReceipt(fileName)
		if !exists {
			response.WriteHeader(http.StatusNotFound)
			return nil
		}
		if err != nil {
			return err
		}
		_, err = response.Write(data)
		if err != nil {
			return err
		}
	case models.DownloadFile:
		data, err := hs.fileStorage.ReadFile(fileName)
		if err != nil {
			return err
		}
		_, err = response.Write(data)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("ошибка определения шага сохранения файла")
	}

	return nil
}
