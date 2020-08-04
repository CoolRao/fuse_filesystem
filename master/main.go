package main

import (
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

func main() {

	mountDir := "/home/abel/tmp/fusesrc"
	copyDir := "/home/abel/tmp/fusedest"
	exitSig := make(chan os.Signal)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, config.SIGUSR1, config.SIGUSR2)

	err := log.InitLogger("./log", "master")
	if err != nil {
		lg.Fatalln("init log error ", err.Error())
	}
	fuseManage, err := fusecore.NewFuseManage(mountDir,copyDir)
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
}
