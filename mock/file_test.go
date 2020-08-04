package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)


func TestReadAll(t *testing.T){
	bytes, err := ioutil.ReadFile("/home/abel/tmp/test/xjrw-test1.txt")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(bytes))
}



func TestCreateFile(t *testing.T){
	file, err := os.Create("/home/abel/tmp/test/test1.txt")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	fmt.Println(file.Stat())
}

func TestWiriteFile(t *testing.T){
	//file, err := os.Open("/home/abel/tmp/test/test1.txt")
	//file, err := os.OpenFile("/home/abel/tmp/test/test1.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	file, err := os.OpenFile("/home/abel/tmp/test/test1.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	n, err := file.WriteString("write test1 ")
	fmt.Println(n,err)

}

func TestReadFile(t *testing.T) {
	file, err := os.Open("/home/abel/tmp/test/test1.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	bytes:=make([]byte,1024)
	n, err := file.Read(bytes)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(n,string(bytes))
}




