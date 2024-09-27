package utils

import (
	"crypto/rand"
	"fmt"
	"golang-template/models"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/joho/godotenv"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GenerateOTP() string {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	otp := n.Int64() + 100000
	return fmt.Sprintf("%06d", otp)
}

func SendOTP(email, otp string) error {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error("Error loading env file:", err)
	}

	companyemail := os.Getenv("COMPANY_EMAIL")

	mj := mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_API_SECRET"))
	log.Println("Mailjet API Key:", os.Getenv("MAILJET_API_KEY"))
	log.Println("Company Email:", companyemail)

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
			Subject:  "Your OTP Code",
			TextPart: fmt.Sprintf("Your OTP code is: %s", otp),
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

func CleanOTP(db *gorm.DB) error {
	var expiredUsers []models.User
	now := time.Now()

	if err := db.Where("user_otp_valid < ? AND user_otp != ''", now).Find(&expiredUsers).Error; err != nil {
		return err
	}

	for _, user := range expiredUsers {
		if err := db.Delete(&user).Error; err != nil {
			return err
		}
	}

	return nil
}