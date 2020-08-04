package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMkdir(t *testing.T){
	err := os.Mkdir("/home/abel/tmp/test/demo",os.ModeDir)
	fmt.Println(err)
}

func TestMkdirAll(t *testing.T){
	err := os.MkdirAll("/home/abel/tmp/test/mkidr/mkdir",os.ModeDir)
	fmt.Println(err)
}

func TestStatefs(t *testing.T){
	// 查看磁盘分区使用情况
	s := syscall.Statfs_t{}
	err := syscall.Statfs("/home/abel/tmp/copy", &s)
	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println(s)

}


func TestState(t *testing.T){
	// 查看磁盘分区使用情况
	//s := syscall.Statfs_t{}
	//err := syscall.Statfs("/dev", &s)

	s := syscall.Stat_t{}
	err := syscall.Stat("/home/abel/tmp/test", &s)

	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println(s)

}