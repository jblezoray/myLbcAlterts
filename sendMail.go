package main

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"
)

type SMTPUser struct {
	Username string
	Password string
	Server   string
	Port     int
}

type smtpTemplateData struct {
	From    string
	To      string
	HowMany int
	AdData  []AdData
}

const emailTemplateFileName string = "sendMail.tmpl"

func SendAdsByMail(smtpUser SMTPUser, from string, to string, adData []AdData) error {

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
		smtpUser.Server+":"+strconv.Itoa(smtpUser.Port),
		smtp.PlainAuth("", smtpUser.Username, smtpUser.Password, smtpUser.Server),
		from,
		[]string{to},
		doc.Bytes())
	return err
}
