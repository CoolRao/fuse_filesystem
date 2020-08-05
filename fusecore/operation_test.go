package fusecore

import (
	"fmt"
	"testing"
)

func TestFileState(t *testing.T) {
	state, err := FileState("test.txt")
	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println(state)
}

func TestFileRead(t *testing.T) {

}
