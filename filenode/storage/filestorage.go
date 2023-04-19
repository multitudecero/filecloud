package storage

import (
	"os"
	"path"
)

type FileStorage struct {
	filePartsFolder string
}

func NewFileStorage(filePartsFolder string) *FileStorage {
	return &FileStorage{
		filePartsFolder: filePartsFolder,
	}
}

func (fs *FileStorage) CreateFile(fileName string, data []byte) error {
	fullFileName := path.Join(fs.filePartsFolder, fileName)
	file, err := os.Create(fullFileName)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (fs *FileStorage) AppendFile(fileName string, data []byte) error {
	fullFileName := path.Join(fs.filePartsFolder, fileName)

	file, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (fs *FileStorage) ReadFile(fileName string) ([]byte, error) {
	fullFileName := path.Join(fs.filePartsFolder, fileName)
	return os.ReadFile(fullFileName)
}
