package mailer

import (
	"fmt"
	"log"

	"github.com/jerson2000/api-qfirst/config"
	"github.com/wneessen/go-mail"
)

func SendOTPMail(toEmail, OTPCode string) error {
	subject := "Subject: Your OTP Code for Verification\n"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>OTP Verification</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f6f8; padding: 20px;">
			<div style="background-color: #ffffff; padding: 30px; border-radius: 10px; max-width: 600px; margin: auto; box-shadow: 0 2px 5px rgba(0,0,0,0.1);">
				<h2 style="color: #2a7de1;">OTP Verification</h2>
				<p>Hello,</p>
				<p>Your One-Time Password (OTP) is:</p>
				<div style="font-size: 24px; font-weight: bold; color: #2a7de1; margin: 20px 0;">%s</div>
				<p>This OTP is valid for the next <strong>5 minutes</strong>. Do not share it with anyone.</p>
				<p>If you did not request this, please ignore this email.</p>
				<p>Thanks,<br>Your Company Team</p>
			</div>
		</body>
		</html>
	`, OTPCode)

	err := SendMailByGmail(subject, body, toEmail, mail.TypeTextHTML)
	if err != nil {
		return err
	}

	return nil
}

func SendMailByGmail(subject, body, toEmail string, mimeType mail.ContentType) error {

	username := config.Mailer_Email
	password := config.Mailer_Password
	client, err := mail.NewClient("smtp.gmail.com", mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		fmt.Printf("failed to create mail client: %s\n", err)
		return err
	}

	message := mail.NewMsg()
	if err := message.From(username); err != nil {
		log.Fatalf("failed to set From address: %s", err)
		return err
	}
	if err := message.To(toEmail); err != nil {
		log.Fatalf("failed to set To address: %s", err)
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mimeType, body)
	if err = client.DialAndSend(message); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
		return err
	}
	return nil
}
