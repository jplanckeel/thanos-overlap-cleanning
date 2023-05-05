package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var err error

func writeCache(meta string, key string, cacheDir string) error {

	filePath := fmt.Sprintf("%s/%s", cacheDir, key)
	dirPath := filepath.Dir(filePath)

	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		//TODO: add debug fmt.Println(dirPath)
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filePath, []byte(meta), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func checkCache(key string, cacheDir string) bool {

	filePath := fmt.Sprintf("%s/%s", cacheDir, key)
	_, err := os.Stat(filePath)

	return os.IsNotExist(err)
}

func readCache(key string, cacheDir string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", cacheDir, key)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", content), nil
}

func purgeCache(cacheDir string) error {
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return err
	}

	return nil
}
