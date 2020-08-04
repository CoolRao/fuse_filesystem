package main

import (
	"flag"
	"fuse_file_system/log"
	lg "log"
)

var host string
var workDir string

func main() {
	flag.StringVar(&host, "host", "0.0.0.0:9999", "work host")
	flag.StringVar(&workDir, "workDir", "/home/abel/tmp/fusedest", "work path")
	flag.Parse()
	log.InitLogger()
	client := NewFileHandler(host,workDir)
	err := client.run()
	if err != nil {
		lg.Fatal(err.Error())
	}
}
