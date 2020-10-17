package main

import (
	"os"

	"github.com/abserari/pet/utils/fileserver"
)

func main() {
	wdir, _ := os.Getwd()
	fileserver.StartFileServer(":9573", wdir)
}
