package model

import (
	"os"
	"time"
)

type DirInfo struct {
	FileList map[string]*FileStat
}

type FileStat struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
	Sys     interface{} `json:"sys"`
	Ip      string      `json:"ip"`
}

func (fs *FileStat) GetSize() int64        { return fs.Size }
func (fs *FileStat) GetMode() os.FileMode  { return fs.Mode }
func (fs *FileStat) GetModTime() time.Time { return fs.ModTime }
func (fs *FileStat) GetSys() interface{}   { return &fs.Sys }

func CovertFileState(fileInfo os.FileInfo) FileStat {
	return FileStat{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		Mode:    fileInfo.Mode(),
		ModTime: fileInfo.ModTime(),
		Sys:     fileInfo.Sys(),
	}
}
