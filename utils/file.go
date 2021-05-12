package utils

import (
	"errors"
	"os"
)

type FileType int32

const (
	Unknown  FileType = 0
	NotExist FileType = 1
	File     FileType = 2
	Dir      FileType = 3
)

func CheckPath(path string) (FileType, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NotExist, nil
		} else {
			return Unknown, err
		}
	}
	if fi.IsDir() {
		return Dir, nil
	} else {
		return File, nil
	}
}

func CheckAndCreateDir(path string) error {
	ft, err := CheckPath(path)
	if err != nil {
		return err
	}
	switch ft {
	case Dir:
		return nil
	case NotExist:
		err = os.MkdirAll(path, os.ModePerm)
		return err
	case File:
		return errors.New("目标路径是一个文件")
	}
	return errors.New("确认目录时遇到了一个未知的错误")
}
