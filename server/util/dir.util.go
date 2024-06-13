package util

import (
	"os"
)

type Dir struct{}

func (ref Dir) Make(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
