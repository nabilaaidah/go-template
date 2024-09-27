package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/sirupsen/logrus"
)

func SendResetEmail(email, token string) error {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error("Error loading env file:", err)
	}

	companyemail := os.Getenv("COMPANY_EMAIL")
	ip := os.Getenv("IP_PROD")

	mj := mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_API_SECRET"))

	resetLink := fmt.Sprintf("http://%s:8080/reset-password?t=%s", ip, token)
	plainTextContent := fmt.Sprintf("Hello,\n\nPlease click the link to reset your password: %s\n\nIf you did not request this, please ignore this email.", resetLink)
	htmlContent := fmt.Sprintf(`
		<p>Hello,</p>
		<p>You requested a password reset.
		<p>Please click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>If you did not request this, please ignore this email.</p>
	`, resetLink)

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: companyemail,
				Name:  "Arena",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
				},
			},
			Subject:  "Password Reset Request",
			TextPart: plainTextContent,
			HTMLPart: htmlContent,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		logrus.Error("Failed to send OTP email:", err)
		return err
	}
	logrus.Infof("OTP email sent successfully: %v", res)

	return nil
}