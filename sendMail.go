package main

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

type smtpTemplateData struct {
	From    string
	To      string
	HowMany int
	AdData  []AdData
}

const emailTemplateFileName string = "sendMail.tmpl"

func SendAdsByMail(emailUser EmailUser, from string, to string, adData []AdData) error {

	// build document from template
	context := &smtpTemplateData{
		from,
		to,
		len(adData),
		adData,
	}

	t := template.Must(template.New(emailTemplateFileName).ParseFiles(emailTemplateFileName))
	var err error
	var doc bytes.Buffer
	err = t.Execute(&doc, context)
	if err != nil {
		return err // error trying to execute mail template
	}
	//fmt.Println(doc.String())

	// send mail.
	err = smtp.SendMail(
		emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer),
		from,
		[]string{to},
		doc.Bytes())
	return err
}
