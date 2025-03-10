package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends a personalized wedding invitation email using Gmail SMTP with an App Password
func SendEmail(guestName, guestEmail, rsvpToken string) error {
	// Load Gmail SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")     // Your Gmail address
	smtpPass := os.Getenv("SMTP_PASSWORD") // App Password

	// Validate SMTP settings
	if smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("❌ SMTP configuration is missing")
	}

	// Generate RSVP Link with path parameter instead of query parameter
	rsvpLink := fmt.Sprintf("https://axeldaphne.com/rsvp/%s", rsvpToken) // ✅ Fixed URL format

	// Email subject and HTML body with RSVP button
	subject := "💍 You're Invited to Axel and Daphne's Wedding Celebration!"
	body := fmt.Sprintf(
		"Dear %s,<br><br>"+
			"With great joy in our hearts, we invite you to celebrate our special day with us! 💍✨<br><br>"+
			"We would love for you to be part of our wedding, creating memories that will last a lifetime.<br><br>"+
			"<strong>Your unique RSVP token: <span style='color:#2c3e50;'>%s</span></strong><br><br>"+ // ✅ Token included
			"<strong style='color:red;'>⚠️ Please do not share your invite token.</strong><br><br>"+ // ✅ Privacy notice
			"To confirm your attendance, please click the button below:<br><br>"+
			"<a href='%s' style='display:inline-block; padding:12px 24px; font-size:16px; "+
			"color:#fff; background-color:#3498db; text-decoration:none; border-radius:5px;'>💌 RSVP Now</a><br><br>"+
			"We truly hope you can join us on this wonderful occasion, and we can't wait to celebrate together! 🎊<br><br>"+
			"With love and excitement,<br>"+
			"<strong>Axel & Daphne 💕</strong>",
		guestName, rsvpToken, rsvpLink, // Pass the token into the email
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
