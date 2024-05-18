package emailVerification

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

func isValidEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func hasValidMX(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	return err == nil && len(mxRecords) > 0
}

func checkSMTP(email string) (bool, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, errors.New("invalid emailVerification format")
	}
	domain := parts[1]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, err
	}

	client, err := smtp.Dial(mxRecords[0].Host + ":25")
	if err != nil {
		return false, err
	}
	defer client.Close()

	err = client.Hello("example.com")
	err = client.Mail("you@example.com")
	if err := client.Rcpt(email); err != nil {
		return false, err
	}

	return true, nil
}

func IsVerifiedEmail(email string) (bool, error) {
	if !isValidEmailFormat(email) {
		return false, errors.New("invalid emailVerification format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, errors.New("invalid emailVerification format")
	}
	domain := parts[1]

	if !hasValidMX(domain) {
		return false, errors.New("invalid emailVerification domain")
	}

	smtpCheck, err := checkSMTP(email)
	if err != nil {
		return false, err
	}
	return smtpCheck, nil
}

func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(email string, token string) error {
	from := "betterpc@mail.ru"
	password := "YkrvYMZ8KgEnqqtyUtGG"
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	to := []string{email}
	subject := "Verify your email address" //change domain in body!! !! ! !! !
	body := fmt.Sprintf("Please click the following link to verify your email address: http://localhost:8080/verify?token=%s", token)

	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err
}
