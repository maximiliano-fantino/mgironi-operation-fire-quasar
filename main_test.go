package main

import (
	"testing"

	"os"
)

func TestMain(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-distances=500,424.26,707.10", "-messages=this..the.complete.message,.is.the..message,.is...message"}
	main()
	os.Args = oldsArgs
}
