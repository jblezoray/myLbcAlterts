package main

import (
	"bytes"
	"fmt"
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

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: Alerte leboncoin.fr, {{.HowMany}} nouveaux résultats
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

{{range $ad := .AdData}}
<li>{{.Title}}</li>
{{end}}
</ul>

Sincerely,

{{.From}}
`

func SendAdsByMail(emailUser EmailUser, from string, to string, adData []AdData) error {

	var err error
	var doc bytes.Buffer

	// build document from template
	context := &smtpTemplateData{
		from,
		to,
		len(adData),
		adData,
	}
	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		return err // error trying to parse mail template
	}
	err = t.Execute(&doc, context)
	if err != nil {
		return err // error trying to execute mail template
	}

	fmt.Println(doc.String())

	// send mail.
	err = smtp.SendMail(
		emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer),
		from,
		[]string{to},
		doc.Bytes())
	return err
}
