package include

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var GmailService *gmail.Service

type DBOutgoingMails struct {
	Model
	ToEmail    string
	Subject    string
	Message    string
	OutMessage string
	Status     string
}

type GmailCredentialsType struct {
	AccessToken  string
	RefreshToken string
	ClientId     string
	ClientSecret string
}

var GmailCredentials GmailCredentialsType

func OpenGmailCredentials(sender string) {
	// Open our jsonFile
	jsonFile, err := os.Open(sender)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	Log.Trace("Successfully Opened gmail.api.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &GmailCredentials)

	//Log.Trace("ClientId ", GmailCredentials.ClientId)
	//Log.Trace("ClientSecret ", GmailCredentials.ClientSecret)
	//Log.Trace("AccessToken ", GmailCredentials.AccessToken)
	//Log.Trace("RefreshToken ", GmailCredentials.RefreshToken)
}

func OAuthGmailService() {

	config := oauth2.Config{
		ClientID:     GmailCredentials.ClientId,
		ClientSecret: GmailCredentials.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  GmailCredentials.AccessToken,
		RefreshToken: GmailCredentials.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		Log.Info("Email service is initialized")
	}

}

func SendEmailOAUTH2(to string, subj string, body string) (status string, messageHash string, errMsg string) {

	//emailBody, err := parseTemplate(template, data)
	//if err != nil {
	//	return false, errors.New("unable to parse email template")
	//}

	var dbOutMessage DBOutgoingMails

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	dbOutMessage.ToEmail = to
	Log.Trace("MESSAGE SEND : ", to)

	subject := "Subject: " + subj + "\n"
	dbOutMessage.Subject = subj
	Log.Trace("MESSAGE SEND : ", subj)
	dbOutMessage.Message = body
	db.Create(&dbOutMessage)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msgString := emailTo +
		subject +
		mime +
		"\n<html><table><tr><td>" +
		body +
		"</td></tr><tr><td><img src='" + os.Getenv("SERVICE_HOST") + "/logo/" + dbOutMessage.ID + "' alt=''></td></tr></table>"
	dbOutMessage.OutMessage = msgString

	msg := []byte(msgString)

	db.Save(&dbOutMessage)

	Log.Trace("MESSAGE SEND : ", body)
	Log.Trace("MESSAGE SEND : ", messageHash)

	message.Raw = base64.URLEncoding.EncodeToString(msg)
	Log.Trace("MESSAGE SEND : ", string(msg))
	// Send the message

	if os.Getenv("DO_NOT_SEND_SAVE_INSTEAD") != "1" {
		a, err := GmailService.Users.Messages.Send("me", &message).Do()

		fmt.Println(a)
		if err != nil {
			Log.Error("MESSAGE SEND : ", err)
			status = "failed"
			dbOutMessage.Status = status
			db.Save(&dbOutMessage)

			//var msgHistory DBOutgoingMailHistory
			//msgHistory.DBOutgoingMailsID = dbOutMessage.ID
			//msgHistory.RecType = "error"
			//msgHistory.HistoryMessage = err.Error()
			//Db.Create(&msgHistory)

			return status, dbOutMessage.ID, err.Error()
		}

		Log.Trace("MESSAGE SEND : Done")

		status = "sent"
		dbOutMessage.Status = status
		db.Save(&dbOutMessage)
		return status, dbOutMessage.ID, ""

	} else {

		f, err := os.Create("temp/email_" + dbOutMessage.ID + ".html")
		if err != nil {
			Log.Error(err)
		}
		defer f.Close()

		k, err := f.Write(msg)
		Log.Tracef("wrote %d bytes\n", k)

		status = "sent"
		dbOutMessage.Status = status
		db.Save(&dbOutMessage)
		return status, dbOutMessage.ID, "file is saved : '" + "temp/email_" + dbOutMessage.ID + ".html'"

	}

}
