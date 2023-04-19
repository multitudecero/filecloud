package main

import (
	"bytes"
	"errors"
	"filecloud/common/dl/models"
	"filecloud/filenode/client"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
)

type FileNode struct {
	ID     string
	Cfg    models.NodeInfo
	client *client.FileNodeClient
}

type NodeReceipt struct {
	Node    *FileNode
	Receipt *models.FileReceipt
}

type BlockReader struct {
	EOF    bool
	Reader io.Reader
}

type NodeOrchestrator struct {
	Config *models.OrcAppCfg
	Nodes  []FileNode
}

func NewNodeOrchestrator(config *models.OrcAppCfg) *NodeOrchestrator {
	no := &NodeOrchestrator{
		Config: config,
	}
	nodes := make([]FileNode, 0)
	for _, nodeCfg := range config.Nodes {
		cl := client.NewFileNodeClient(&nodeCfg)
		node := FileNode{
			ID:     nodeCfg.ID,
			Cfg:    nodeCfg,
			client: cl,
		}
		nodes = append(nodes, node)
	}
	no.Nodes = nodes
	return no
}

func (no *NodeOrchestrator) Send(nodeID int, fileName string, step models.UploadStep, buf []byte) error {
	if len(no.Nodes) == 0 {
		return fmt.Errorf("отсутствуют приемники для данных")
	}

	err := no.Nodes[nodeID].client.Put(fileName, bytes.NewReader(buf), step)
	if err != nil {
		return err
	}

	return nil
}

func (no *NodeOrchestrator) SendReceipts(fileName string) error {
	for i, node := range no.Nodes {
		receipt := models.FileReceipt{
			NodeID:    node.ID,
			FileName:  fileName,
			NodeCount: len(no.Nodes),
			NodeOrder: i,
			BlockSize: no.Config.BlockSize,
			Nodes:     no.Config.Nodes,
		}
		err := node.client.PutReceipt(&receipt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (no *NodeOrchestrator) SendBlocks(fileName string, request *http.Request) error {
	buf := make([]byte, no.Config.BlockSize)
	nodeID := 0
	step := models.UploadStepStart
	for {
		n, err := request.Body.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n > 0 {
			toSend := buf[:n]
			log.Println(string(toSend))
			err = no.Send(nodeID, fileName, step, toSend)
			if err != nil {
				return err
			}

			if nodeID < len(no.Nodes)-1 {
				nodeID++
			} else {
				nodeID = 0
				step = models.UploadStepAppend
			}
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
	}
}

func (no *NodeOrchestrator) GetFile(fileName string, response http.ResponseWriter) error {
	nodes, err := no.GetAndValidateReceipts(fileName)
	if err != nil {
		return err
	}

	readers := make([]BlockReader, 0)
	for _, node := range nodes {
		reader, err := node.Node.client.GetReader(fileName)
		if err != nil {
			return err
		}
		readers = append(readers, BlockReader{Reader: reader, EOF: false})
	}

	buf := make([]byte, nodes[0].Receipt.BlockSize)
	toSend := make([]byte, 0)
	nodeID := -1
	for {
		if nodeID+1 < len(nodes) {
			nodeID++
		} else {
			nodeID = 0
		}

		n, err := readers[nodeID].Reader.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n > 0 {
			toSend = append(toSend, buf[:n]...)
		}
		if errors.Is(err, io.EOF) {
			readers[nodeID].EOF = true
			found := false
			for _, r := range readers {
				if !r.EOF {
					found = true
					break
				}
			}
			if !found {
				_, err = response.Write(toSend)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}
}

func (no *NodeOrchestrator) GetAndValidateReceipts(fileName string) ([]NodeReceipt, error) {
	nodes := make([]NodeReceipt, 0)

	receipts := make([]models.FileReceipt, 0)
	for _, node := range no.Nodes {
		receipt, err := node.client.GetReceipt(fileName)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, *receipt)
	}

	// проверка корректности
	// for _, receipt1 := range receipts {
	// 	for _, receipt2 := range receipts {
	// 		for _, node1 := range receipt1.Nodes {
	// 			for _, node2 := range receipt2.Nodes {
	// 				_ = node1
	// 				_ = node2
	// TODO будем считать, что тикты везде равны
	// 			}
	// 		}
	// 	}
	// }

	// проверяем, на месте ли ноды
	for _, receipt := range receipts {
		found := false
		for _, node := range no.Nodes {
			if node.ID == receipt.NodeID {
				found = true
				nodes = append(nodes, NodeReceipt{Node: &node, Receipt: &receipt})
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("не найдена нода %+v. целостность файла нарушена", receipt.NodeID)
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Receipt.NodeOrder < nodes[j].Receipt.NodeOrder
	})

	return nodes, nil
}
