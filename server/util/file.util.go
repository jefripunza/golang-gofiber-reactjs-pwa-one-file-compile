package util

import (
	"backend/server/env"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct{}

func (ref File) IsExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (ref File) Create(filePath string, data []byte, perm fs.FileMode) error {
	err := ioutil.WriteFile(filePath, data, perm)
	if err != nil {
		return err
	}
	return nil
}

func (ref File) CreateIfNotExist(filePath string) error {
	file := File{}
	if !file.IsExist(filePath) {
		err := file.Create(filePath, nil, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ref File) AddExtensionIfNotExist(filePath string, replaceExtension string) string {
	extension := filepath.Ext(filePath)
	if extension != replaceExtension {
		if !strings.HasPrefix(replaceExtension, ".") {
			replaceExtension = "." + replaceExtension
		}
		if extension == "" {
			filePath += replaceExtension
		} else {
			filePath = strings.Replace(filePath, extension, replaceExtension, 1)
		}
	}
	return filePath
}

func (ref File) ClearAllOnFolder(folderPath string) error {
	pwd := env.GetPwd()
	folderPath = filepath.Join(pwd, folderPath)
	if folderPath == pwd { // pengaman...
		return fmt.Errorf("cannot clear core, please :')")
	}
	dir, err := os.Open(folderPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Baca isi folder
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Hapus semua file dalam folder
	for _, fileInfo := range fileInfos {
		filePath := filepath.Join(folderPath, fileInfo.Name())
		err := os.RemoveAll(filePath)
		if err != nil {
			return err
		}
		fmt.Println("Deleted:", filePath)
	}

	return nil
}
