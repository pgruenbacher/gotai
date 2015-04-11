package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/pgruenbacher/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFile(path string, ptr interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	err = readTOML(string(data), ptr)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

func ReadDir(dir string, ptr interface{}) error {
	entries, e := ioutil.ReadDir(dir)
	if e != nil {
		log.Error("%v", e)
	}

	e = parseFiles(entries, dir, ptr)

	return nil
}

func parseFiles(entries []os.FileInfo, dirName string, ptr interface{}) (err error) {
	for _, entry := range entries {
		var data []byte
		if entry.IsDir() {
			break
		}
		if validName(entry.Name(), dirName) {
			var path string
			path = filepath.Join(dirName, entry.Name())
			data, err = ioutil.ReadFile(path)
			if err != nil {
				log.Error("%v", err)
				return err
			}
			err = readTOML(string(data), ptr)
			if err != nil {
				log.Error("%v", err)
				return err
			}
		}
	}
	return nil
}

func readTOML(data string, t interface{}) (err error) {
	_, err = toml.Decode(data, t)
	return err
}
func validName(str1, str2 string) (valid bool) {
	valid = true
	ln1 := len(str1)
	ln2 := len(str2)
	if str1[ln1-4:ln1] != "toml" {
		valid = false
	}
	if str1[0:ln2] != str2 {
		valid = false
	}
	return valid
}
