package main

import (
	"flag"
	"fmt"
	"fuse_file_system/config"
	"fuse_file_system/fusecore"
	"fuse_file_system/log"
	lg "log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)


var host string
var mountDir string
var copyDir = "/home/abel/tmp/fusedest"
func main() {


	flag.StringVar(&host,"host","0.0.0.0:8888","host addr ")
	flag.StringVar(&mountDir,"mountDir","/home/abel/tmp/fusesrc","mount dir path ")
	exitSig := make(chan os.Signal)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, config.SIGUSR1, config.SIGUSR2)
	log.InitLogger()
	fuseManage, err := fusecore.NewFuseManage(mountDir,copyDir,host)
	if err != nil {
		lg.Fatalln(err.Error())
	}

	go func() {
		for c := range exitSig {
			fmt.Println("exit sig :", c)
			fuseManage.Close()
			unmount(mountDir)
			os.Exit(0)
		}
	}()
	err = fuseManage.Run()
	if err != nil {
		lg.Fatalln(err.Error())
	}
}
func unmount(mountDir string) {
	cmd := exec.Command("fusermount", "-u", mountDir)
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("unMount dir error ", err.Error()))
	}
	os.RemoveAll(copyDir)

}
