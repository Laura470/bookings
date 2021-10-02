package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/Laura470/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {
	//faccio partire una routine con una funzione anonima che funziona in back ground
	//è un for loop infinito
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	//altre cose per un live server

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	//se non ho scelto template
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else { //se ho scelto template
		//vado a leggere dal disco
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))

		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}
