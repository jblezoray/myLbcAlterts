package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"
)

type smtpTemplateData struct {
	From        string
	To          string
	HowMany     int
	AdsBySearch map[Search][]AdData
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: Alerte leboncoin.fr, {{.HowMany}} nouveaux résultats
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

<body>
	<div style="max-width:750px">
		<p>Recherches :</p> 
		<ul>
		{{range $search, $ads := .AdsBySearch}}
			<li>
				<a href="http://www.leboncoin.fr/{{$search.Terms}}">{{$search.Name}}</a>
				({{len $ads}} résultat.s)
			</li>
		{{end}}
		</ul>

		{{range $search, $ads := .AdsBySearch}}
		<h2><a href="http://www.leboncoin.fr/{{$search.Terms}}">{{$search.Name}}</a></h2>
		<ul style="padding: 0; overflow: hidden; ">
			{{range $ad := $ads}}
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
		{{end}}
	</div>
</body>`

func SendAdsByMail(config Configuration, adData map[Search][]AdData) error {

	// count ads.
	howMany := 0
	for _, ads := range adData {
		howMany = howMany + len(ads)
	}
	if howMany == 0 {
		fmt.Println("No new Data")
		return nil
	}

	// build document from template
	context := &smtpTemplateData{
		config.MailFrom,
		config.MailTo,
		howMany,
		adData,
	}

	t := template.Must(template.New("emailTemplate").Parse(emailTemplate))
	var err error
	var doc bytes.Buffer
	err = t.Execute(&doc, context)
	if err != nil {
		return err // error trying to execute mail template
	}
	// fmt.Println(doc.String())

	// send mail.
	err = smtp.SendMail(
		config.SMTPUser.Server+":"+strconv.Itoa(config.SMTPUser.Port),
		smtp.PlainAuth("", config.SMTPUser.Username, config.SMTPUser.Password, config.SMTPUser.Server),
		config.MailFrom,
		[]string{config.MailTo},
		doc.Bytes())
	return err
}
