package notifiers

import (
	"bytes"
	"log"
	"net/smtp"
	"text/template"
)

type EmailNotifier struct {
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<html>
<body>
{{.Body}}
</body>
</html>`

type EmailData struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (e *EmailNotifier) SendMessage(event_name string, target string, data []byte) {
	var err error
	var doc bytes.Buffer

	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		log.Fatal("t.Parse ", err)
	}
	context := &EmailData{
		From:    "Springest Dev <developers@springest.nl>",
		To:      target,
		Subject: "Email subject line",
		Body:    string(data),
	}
	err = t.Execute(&doc, context)
	if err != nil {
		log.Fatal("t.Execute ", err)
	}

	// TODO: setup env variables to support multiple envs
	// TODO: Set up real test mode instead of using mailcatcher
	auth := smtp.PlainAuth("", "", "", "localhost:1025")
	err = smtp.SendMail("localhost:1025", auth, "test@example.com", []string{"recipient@example.com"}, doc.Bytes())
	if err != nil {
		log.Fatal("smtp.SendMail ", err)
	}
}