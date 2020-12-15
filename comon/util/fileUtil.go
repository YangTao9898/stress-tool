package util

import (
	"errors"
	"io/ioutil"
	"os"
)

// 列出文件夹下面所有的文件名，不包含文件夹
func ListChildWholeFileName(folderPath string) ([]string, error) {
	finfo, e := os.Stat(folderPath)
	if e != nil && !os.IsExist(e) {
		return nil, errors.New("路径[" + folderPath + "]不存在")
	}
	if !finfo.IsDir() {
		return nil, errors.New("路径[" + folderPath + "]不是一个文件夹")
	}
	files, e := ioutil.ReadDir(folderPath)
	if e != nil {
		return nil, e
	}
	var nameArr []string
	var path string
	length := len(folderPath)
	if folderPath[length - 1] == '/' && length > 1  {
		path = folderPath[:len(folderPath) - 2]
	} else {
		path = folderPath
	}
	for _, f := range files {
		if !f.IsDir() {
			nameArr = append(nameArr, path + "/" + f.Name())
		}
	}
	return nameArr, nil
}
