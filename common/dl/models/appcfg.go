package models

type OrcAppCfg struct {
	HostAndPort string
	// BlockSize ограничение разбиваемого блока по размеру для разбивки файла между нодами
	// Например, 1024
	BlockSize int
	Nodes     []NodeInfo
}

type NodeAppCfg struct {
	ID              string
	HostAndPort     string
	ReceiptFolder   string
	FilePartsFolder string
}
