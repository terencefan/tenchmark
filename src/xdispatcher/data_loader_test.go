package xdispatcher

import (
	"fmt"
	"testing"
)

func TestDirLoad(t *testing.T) {
	loader, err := NewDirDataLoader("../../../archimedes/cases")
	if err != nil {
		panic(err)
	}

	err = loader.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(loader.GetApi("fetchCtrByLR"))
}
