package main

import (
	"mask_api_gin/src"
)

func main() {
	src.ConfigurationInit()
	if err := src.RunServer(); err != nil {
		src.ConfigurationClose()
	}
}
