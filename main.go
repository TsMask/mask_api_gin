package main

import (
	"fmt"
	"mask_api_gin/src"
	"runtime"
)

func main() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)

	src.Configuration()
	src.RunServer()
}
