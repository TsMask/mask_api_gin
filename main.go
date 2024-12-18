package main

import (
	"mask_api_gin/src"

	"embed"
	"log"
)

//go:embed src/assets/**
var assetsDir embed.FS

//go:embed src/config/*.yaml
var configDir embed.FS

func main() {
	src.ConfigurationInit(&assetsDir, &configDir)
	defer src.ConfigurationClose()

	if err := src.RunServer(); err != nil {
		log.Fatalf("run server error: %v", err)
	}
}
