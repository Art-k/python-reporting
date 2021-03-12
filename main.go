package main

import (
	inc "./include"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR opening file: %v", err)
	}
	defer f.Close()

	err = godotenv.Load("p.env")
	if err != nil {
		log.Fatal("ERROR loading .env file")
	}

	inc.InitApplication(f)
	inc.ApiProcessing()
}
