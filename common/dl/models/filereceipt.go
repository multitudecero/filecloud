package models

type FileReceipt struct {
	NodeID string
	// имя файла с данными
	FileName string
	// кол-во нод, которые участвовали в приеме файла
	NodeCount int
	// порядковый номер текущей ноды
	NodeOrder int
	// размер блока данных, на который файл был разбит
	BlockSize int
	// список узлов с частями файлов
	Nodes []NodeInfo
}
