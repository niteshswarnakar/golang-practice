package testmail

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendMail() {
	fmt.Println("Sending mail ...")
	auth := smtp.PlainAuth(
		"",
		"tuntun.broda@gmail.com",
		"vdee fkke zral jwqh",
		"smtp.gmail.com",
	)

	msg := "Subject: This is the subject\nTunTun - 654234"

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"tuntun.broda@gmail.com",
		[]string{"nepbytechannel@gmail.com"},
		[]byte(msg),
	)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Email sent successfully...")
}
