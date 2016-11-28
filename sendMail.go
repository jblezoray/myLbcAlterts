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

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: Alerte leboncoin.fr, {{.HowMany}} nouveaux résultats
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

<body>
<div style="max-width:750px">
    <ul>
    {{range $ad := .AdData}}
    <li style="list-style:none; margin-bottom:20px; clear:both; background:#EAEBF0; border-top:1px solid #ccc;">
        <div style="float:left; width:180px; padding:20px 0">
            <a href="{{.Url}}"><img src="{{.ThumbSrc}}"></a>
        </div>
        <div style="padding:20px 0 0 220px; background:#ffffff">
            <a href="{{.Url}}" style="font-size:14px; font-weight:bold; color:#369; text-decoration:none;">
                {{.Title}}
            </a>
            <div></div>
            <div>
                {{.LocationTown}} / {{.LocationRegion}}
            </div>
            <div style="float:left; line-height:32px; font-size:14px; font-weight:bold;">
                {{.Price}}&nbsp;&euro;
            </div>
            <div style="float:right; line-height:32px; text-align:right;">
                {{.DateStr}}
            </div>
        </div>
    </li>
    {{end}}
    </ul>
</div>
</body>`

func SendAdsByMail(smtpUser SMTPUser, from string, to string, adData []AdData) error {

	// build document from template
	context := &smtpTemplateData{
		from,
		to,
		len(adData),
		adData,
	}

	t := template.Must(template.New("emailTemplate").Parse(emailTemplate))
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
