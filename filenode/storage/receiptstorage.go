package storage

import (
	"fmt"
	"os"
	"path"
)

type ReceiptStorage struct {
	Folder string
}

func NewReceiptStorage(folder string) *ReceiptStorage {
	return &ReceiptStorage{
		Folder: folder,
	}
}

func (rs *ReceiptStorage) PutReceipt(fileName string, data []byte) error {
	receiptFileName := path.Join(rs.Folder, rs.genReceiptFileName(fileName))
	return os.WriteFile(receiptFileName, data, 0775)
}

// GetReceipt возвращает квитанцию о файле, признак наличия файла в папке хранения и/или ошибку
func (rs *ReceiptStorage) GetReceipt(dataFileName string) ([]byte, bool, error) {
	receiptFileName := path.Join(rs.Folder, rs.genReceiptFileName(dataFileName))
	if stat, err := os.Stat(receiptFileName); err == nil && !stat.IsDir() {
		data, err := os.ReadFile(receiptFileName)
		if err != nil {
			return nil, true, err
		}
		return data, true, nil
	} else {
		if err != nil {
			return nil, false, err
		}
		return nil, false, fmt.Errorf("it isn't file")
	}
}

func (rs *ReceiptStorage) genReceiptFileName(dataFileName string) string {
	return fmt.Sprintf("%+v.receipt", dataFileName)
}
