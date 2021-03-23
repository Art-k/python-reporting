package main

import (
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"log"
	"os"
	"python-reporter/pkg/include"
	"time"
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

	include.InitApplication(f)
	include.Sch = gocron.NewScheduler(time.UTC)
	go include.ApplicationStartAllTasks()

	include.OpenGmailCredentials(os.Getenv("SENDER_EMAIL"))
	include.OAuthGmailService()

	include.ApiProcessing()
}
