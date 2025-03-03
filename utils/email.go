package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends a personalized email invitation using Gmail SMTP with App Password
func SendEmail(guestName, guestEmail, rsvpToken string) error {
	// Load Gmail SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")     // Your Gmail address
	smtpPass := os.Getenv("SMTP_PASSWORD") // App Password

	// Validate Gmail settings
	if smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("❌ SMTP configuration is missing")
	}

	// Generate RSVP Link
	rsvpLink := fmt.Sprintf("https://rsvp.example.com?token=%s", rsvpToken)

	// Email subject and body with guest's name
	subject := "You're Invited to Our Wedding! 🎉"
	body := fmt.Sprintf(
		"Dear %s,\n\n"+
			"You are invited to our wedding! 🎊\n\n"+
			"Please click the link below to RSVP:\n%s\n\n"+
			"We look forward to celebrating with you!\n\n"+
			"Best regards,\nThe Wedding Team 💍",
		guestName, rsvpLink,
	)

	// Format the email message
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Build Gmail SMTP server address
	serverAddress := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Set up authentication using Gmail App Password
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send the email
	err := smtp.SendMail(serverAddress, auth, smtpUser, []string{guestEmail}, []byte(message))
	if err != nil {
		log.Printf("❌ Error sending email to %s: %v", guestEmail, err)
		return err
	}

	log.Printf("✅ RSVP invitation email sent to %s", guestEmail)
	return nil
}
