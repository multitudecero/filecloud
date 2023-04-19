package client

import (
	"bytes"
	"encoding/json"
	"filecloud/common/dl/models"
	"fmt"
	"io"
	"net/http"
)

type FileNodeClient struct {
	config *models.NodeInfo
	url    string
}

func NewFileNodeClient(config *models.NodeInfo) *FileNodeClient {
	return &FileNodeClient{
		config: config,
		url:    fmt.Sprintf("http://%+v", config.HostAndPort),
	}
}

func (fnc *FileNodeClient) PutReceipt(receipt *models.FileReceipt) error {
	data, err := json.MarshalIndent(receipt, "", "	")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	return fnc.Put(receipt.FileName, reader, models.UploadStepReceipt)
}

func (fnc *FileNodeClient) Put(fileName string, reader io.Reader, uploadStep models.UploadStep) error {
	request, err := http.NewRequest(http.MethodPut, fnc.url, reader)
	if err != nil {
		return err
	}

	request.Header.Add(models.UploadStepHeader, string(uploadStep))
	request.Header.Add(models.FileNameHeader, fileName)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		if body, err := io.ReadAll(resp.Body); err == nil {
			return fmt.Errorf("ошибка отправки в хранилище: %+v", string(body))
		} else {
			return fmt.Errorf("ошибка отправки в хранилище, но при распаковке ошибки возникла другая ошибка: %+v", err.Error())
		}
	}
	return nil
}

func (fnc *FileNodeClient) GetReceipt(fileName string) (*models.FileReceipt, error) {
	request, err := http.NewRequest(http.MethodGet, fnc.url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add(models.UploadStepHeader, string(models.DownloadReceipt))
	request.Header.Add(models.FileNameHeader, fileName)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		if body, err := io.ReadAll(resp.Body); err == nil {
			return nil, fmt.Errorf("ошибка получения тикета: %+v", string(body))
		} else {
			return nil, fmt.Errorf("ошибка получения тикета, но при распаковке ошибки возникла другая ошибка: %+v", err.Error())
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var receipt models.FileReceipt
	err = json.Unmarshal(data, &receipt)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (fnc *FileNodeClient) GetReader(fileName string) (io.Reader, error) {
	request, err := http.NewRequest(http.MethodGet, fnc.url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add(models.UploadStepHeader, string(models.DownloadFile))
	request.Header.Add(models.FileNameHeader, fileName)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		if body, err := io.ReadAll(resp.Body); err == nil {
			return nil, fmt.Errorf("ошибка получения тикета: %+v", string(body))
		} else {
			return nil, fmt.Errorf("ошибка получения тикета, но при распаковке ошибки возникла другая ошибка: %+v", err.Error())
		}
	}

	return resp.Body, nil
}
