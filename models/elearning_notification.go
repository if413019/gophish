package models

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	log "github.com/gophish/gophish/logger"
	"golang.org/x/crypto/bcrypt"
)

// ELearningConfig holds the configuration for e-learning notifications
type ELearningConfig struct {
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	SMTPFrom        string
	SMTPUseTLS      bool
	BaseURL         string
	DefaultPassword string
	EmailSubject    string
	CompanyName     string
}

var elearningConfig *ELearningConfig

// LoadELearningConfig loads configuration from .env file
func LoadELearningConfig() (*ELearningConfig, error) {
	if elearningConfig != nil {
		return elearningConfig, nil
	}

	config := &ELearningConfig{
		// Default values
		SMTPHost:        "localhost",
		SMTPPort:        587,
		SMTPUseTLS:      true,
		DefaultPassword: "TempPass123!",
		EmailSubject:    "Security Awareness Training Required",
		CompanyName:     "Your Organization",
		BaseURL:         "https://localhost:3333",
	}

	// Try to read .env file
	file, err := os.Open(".env")
	if err != nil {
		log.Warnf("Could not open .env file: %v. Using defaults.", err)
		elearningConfig = config
		return config, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ELEARNING_SMTP_HOST":
			config.SMTPHost = value
		case "ELEARNING_SMTP_PORT":
			if port, err := strconv.Atoi(value); err == nil {
				config.SMTPPort = port
			}
		case "ELEARNING_SMTP_USERNAME":
			config.SMTPUsername = value
		case "ELEARNING_SMTP_PASSWORD":
			config.SMTPPassword = value
		case "ELEARNING_SMTP_FROM":
			config.SMTPFrom = value
		case "ELEARNING_SMTP_USE_TLS":
			config.SMTPUseTLS = strings.ToLower(value) == "true"
		case "ELEARNING_BASE_URL":
			config.BaseURL = value
		case "ELEARNING_DEFAULT_PASSWORD":
			config.DefaultPassword = value
		case "ELEARNING_EMAIL_SUBJECT":
			config.EmailSubject = value
		case "ELEARNING_COMPANY_NAME":
			config.CompanyName = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("Error reading .env file: %v", err)
	}

	elearningConfig = config
	return config, nil
}

// SendEnrollmentNotification sends an email notification about course enrollment
func SendEnrollmentNotification(userEmail string, course Course, isNewUser bool, tempPassword string) error {
	log.Infof("SendEnrollmentNotification called for %s, course: %s, isNewUser: %t", userEmail, course.Name, isNewUser)
	
	config, err := LoadELearningConfig()
	if err != nil {
		log.Errorf("Failed to load e-learning config: %v", err)
		return fmt.Errorf("failed to load e-learning config: %v", err)
	}
	log.Infof("E-learning config loaded - SMTP Host: %s, Username: %s", config.SMTPHost, config.SMTPUsername)

	// Skip if SMTP is not configured
	if config.SMTPUsername == "" || config.SMTPPassword == "" {
		log.Warn("E-learning SMTP not configured. Skipping enrollment notification.")
		return nil
	}

	subject := config.EmailSubject
	body := createEnrollmentEmailBody(userEmail, course, isNewUser, tempPassword, config)
	log.Infof("Email content created for %s, subject: %s", userEmail, subject)

	log.Infof("Attempting to send email to %s", userEmail)
	err = sendEmail(config, userEmail, subject, body)
	if err != nil {
		log.Errorf("Failed to send email to %s: %v", userEmail, err)
		return err
	}
	log.Infof("Email sent successfully to %s", userEmail)
	return nil
}

// createEnrollmentEmailBody creates the email body for enrollment notification
func createEnrollmentEmailBody(userEmail string, course Course, isNewUser bool, tempPassword string, config *ELearningConfig) string {
	courseURL := fmt.Sprintf("%s/courses/%d/preview", config.BaseURL, course.Id)
	loginURL := fmt.Sprintf("%s/login", config.BaseURL)

	body := fmt.Sprintf(`
Dear User,

You have fallen for a phishing simulation conducted by %s. This is part of our ongoing security awareness program.

Don't worry - this was a test, and no harm was done. However, this incident shows the importance of cybersecurity awareness.

As part of your security awareness training, you have been enrolled in the following course:

Course: %s
Description: %s

To access your training:
1. Visit: %s
2. Login with:
   - Username: %s`, config.CompanyName, course.Name, course.Description, loginURL, userEmail)

	if isNewUser {
		body += fmt.Sprintf(`
   - Temporary Password: %s

IMPORTANT: You must change this password on your first login for security reasons.`, tempPassword)
	} else {
		body += `
   - Use your existing password`
	}

	body += fmt.Sprintf(`

3. Complete the required training modules

Course Preview: %s

This training is mandatory and must be completed within the specified timeframe. The course will help you:
- Recognize phishing attempts
- Understand security best practices  
- Protect yourself and our organization from cyber threats

If you have any questions about this training, please contact your IT Security team.

Best regards,
Security Awareness Team
%s

---
This email was sent because you interacted with a security awareness campaign. 
This is an automated message from the Gophish security awareness platform.
`, courseURL, config.CompanyName)

	return body
}

// sendEmail sends an email using the configured SMTP settings
func sendEmail(config *ELearningConfig, to, subject, body string) error {
	// For SMTP MAIL FROM command, use only the email address, not display name
	fromAddr := config.SMTPUsername
	
	// For the message headers, we can use the full display name
	fromDisplay := config.SMTPFrom
	if fromDisplay == "" {
		fromDisplay = config.SMTPUsername
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", fromDisplay, to, subject, body)

	// Connect to SMTP server
	serverAddr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)
	log.Infof("Connecting to SMTP server: %s (TLS: %t)", serverAddr, config.SMTPUseTLS)
	
	var c *smtp.Client
	var err error

	if config.SMTPUseTLS {
		// For Gmail and most modern SMTP servers, use STARTTLS instead of direct TLS
		c, err = smtp.Dial(serverAddr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
		
		// Start TLS if supported
		if ok, _ := c.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName: config.SMTPHost,
			}
			if err = c.StartTLS(tlsConfig); err != nil {
				c.Close()
				return fmt.Errorf("failed to start TLS: %v", err)
			}
		}
	} else {
		c, err = smtp.Dial(serverAddr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
	}
	defer c.Close()

	// Authenticate
	if config.SMTPUsername != "" && config.SMTPPassword != "" {
		log.Infof("Authenticating with SMTP server for user: %s", config.SMTPUsername)
		auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %v", err)
		}
		log.Infof("SMTP authentication successful")
	}

	// Send email
	log.Infof("Sending email from %s to %s", fromAddr, to)
	if err = c.Mail(fromAddr); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	log.Infof("Enrollment notification sent successfully to %s", to)
	return nil
}

// GenerateSecurePassword generates a secure random password
func GenerateSecurePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}

// UpdateUserPasswordForFirstLogin sets up a user for first-time login
func UpdateUserPasswordForFirstLogin(user *User, tempPassword string) error {
	// Hash the temporary password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update user with hashed password and mark for password change
	user.Hash = string(hashedPassword)
	user.PasswordChangeRequired = true

	return db.Save(user).Error
}