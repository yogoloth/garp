package main

import (
	"log"
	"os"
)

var (
	LogDebug *log.Logger
	LogInfo  *log.Logger
	LogError *log.Logger
)

func init() {
	log.Println("init ...")
	LogDebug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	LogInfo = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
