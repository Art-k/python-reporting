package main

import (
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

import (
	inc "./include"
)

// TODO Add logo to email. Add logo processing, we need to know if email opened.

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
	inc.Sch = gocron.NewScheduler(time.UTC)
	go inc.ApplicationStartAllTasks()

	inc.OpenGmailCredentials(os.Getenv("SENDER_EMAIL"))
	inc.OAuthGmailService()

	inc.ApiProcessing()
}
