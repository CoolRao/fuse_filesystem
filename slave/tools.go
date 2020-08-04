package main

import (
	"fmt"
	"fuse_file_system/model"
	"os"
	"path/filepath"
	"strings"
)

func TraverseDir(workDir string, dirType string) (map[string]*model.FileStat, error) {
	res := make(map[string]*model.FileStat)
	err := filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		split := strings.Split(path, dirType)
		if len(split) > 1 && split[1] != "" {
			stat := &model.FileStat{}
			stat.Size = info.Size()
			stat.Name = info.Name()
			stat.Mode = info.Mode()
			stat.IsDir = info.IsDir()
			stat.ModTime = info.ModTime()
			stat.Sys = info.Sys()
			res[fmt.Sprintf("%s%s", dirType, split[1])] = stat
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetAbsolutePath(workDir, fileName string) (string, bool, error) {
	var absolutePath string
	err := filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if strings.HasSuffix(path, fileName) {
			absolutePath = path
			return nil
		}
		return nil
	})
	if err != nil {
		return absolutePath, false, err
	}
	if absolutePath == "" {
		return absolutePath, false, nil
	}
	return absolutePath, true, nil
}

