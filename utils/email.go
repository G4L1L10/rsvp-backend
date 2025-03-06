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

	// Generate RSVP Link with path parameter instead of query parameter
	rsvpLink := fmt.Sprintf("http://localhost:5173/rsvp/%s", rsvpToken) // ✅ Fixed URL format

	// Email subject and HTML body with RSVP button
	subject := "You're Invited to Our Wedding! 🎉"
	body := fmt.Sprintf(
		"Dear %s,<br><br>"+
			"You are invited to our wedding! 🎊<br><br>"+
			"Please confirm your attendance by clicking the button below:<br><br>"+
			"<a href='%s' style='display:inline-block; padding:12px 24px; font-size:16px; "+
			"color:#fff; background-color:#3498db; text-decoration:none; border-radius:5px;'>RSVP Now</a><br><br>"+
			"We look forward to celebrating with you!<br><br>"+
			"Best regards,<br>The Wedding Team 💍",
		guestName, rsvpLink,
	)

	// Format the email message with proper headers for HTML content
	message := fmt.Sprintf("Subject: %s\nMIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\n%s", subject, body)

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
