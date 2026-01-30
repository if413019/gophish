package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
)

// ELearningSettings stores the configuration for e-learning enrollment notifications
// This is a singleton table - only one row with id=1
type ELearningSettings struct {
	Id           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	SMTPId       int64     `json:"smtp_id" gorm:"column:smtp_id"`
	SMTP         SMTP      `json:"smtp" gorm:"-"`
	EmailSubject string    `json:"email_subject" gorm:"column:email_subject"`
	EmailHTML    string    `json:"email_html" gorm:"column:email_html"`
	BaseURL      string    `json:"base_url" gorm:"column:base_url"`
	CompanyName  string    `json:"company_name" gorm:"column:company_name"`
	ModifiedDate time.Time `json:"modified_date" gorm:"column:modified_date"`
}

// TableName specifies the database tablename for Gorm to use
func (s ELearningSettings) TableName() string {
	return "elearning_settings"
}

// ErrBaseURLRequired is thrown when base_url is not specified
var ErrBaseURLRequired = errors.New("Base URL is required")

// ErrCompanyNameRequired is thrown when company_name is not specified
var ErrCompanyNameRequired = errors.New("Company Name is required")

// ErrSMTPIdRequired is thrown when smtp_id is not specified
var ErrSMTPIdRequired = errors.New("Sending Profile is required")

// Validate ensures that ELearningSettings are valid
func (s *ELearningSettings) Validate() error {
	if s.BaseURL == "" {
		return ErrBaseURLRequired
	}
	if s.CompanyName == "" {
		return ErrCompanyNameRequired
	}
	if s.SMTPId <= 0 {
		return ErrSMTPIdRequired
	}
	return nil
}

// GetELearningSettings retrieves the e-learning settings from the database
// If no settings exist, returns empty settings with default values
func GetELearningSettings() (ELearningSettings, error) {
	settings := ELearningSettings{}

	// Try to get existing settings (singleton with id=1)
	err := db.Where("id = ?", 1).First(&settings).Error
	if err != nil {
		// Return default settings if none exist
		log.Info("No e-learning settings found, returning defaults")
		return ELearningSettings{
			Id:           1,
			EmailSubject: "Security Awareness Training Required",
			EmailHTML:    GetDefaultEmailTemplate(),
			BaseURL:      "https://localhost:3333",
			CompanyName:  "Your Organization",
		}, nil
	}

	// Load the associated SMTP profile if smtp_id is set
	if settings.SMTPId > 0 {
		smtp := SMTP{}
		err = db.Where("id = ?", settings.SMTPId).First(&smtp).Error
		if err == nil {
			settings.SMTP = smtp
		} else {
			log.Warnf("Could not load SMTP profile %d for e-learning settings: %v", settings.SMTPId, err)
		}
	}

	return settings, nil
}

// PostELearningSettings creates or updates the e-learning settings
func PostELearningSettings(settings *ELearningSettings) error {
	// Validate the settings
	err := settings.Validate()
	if err != nil {
		log.Error(err)
		return err
	}

	// Always use id=1 for singleton
	settings.Id = 1
	settings.ModifiedDate = time.Now().UTC()

	// Check if settings exist
	existing := ELearningSettings{}
	err = db.Where("id = ?", 1).First(&existing).Error

	if err != nil {
		// No existing settings, create new
		log.Info("Creating new e-learning settings")
		err = db.Create(settings).Error
	} else {
		// Update existing settings
		log.Info("Updating e-learning settings")
		err = db.Model(&existing).Updates(map[string]interface{}{
			"smtp_id":       settings.SMTPId,
			"email_subject": settings.EmailSubject,
			"email_html":    settings.EmailHTML,
			"base_url":      settings.BaseURL,
			"company_name":  settings.CompanyName,
			"modified_date": settings.ModifiedDate,
		}).Error
	}

	if err != nil {
		log.Error(err)
		return err
	}

	// Reload to get the full settings with SMTP profile
	*settings, err = GetELearningSettings()
	return err
}

// GetDefaultEmailTemplate returns the default HTML email template for enrollment notifications
func GetDefaultEmailTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Security Awareness Training Required</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
            border-radius: 10px 10px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px;
            border: 1px solid #e0e0e0;
            border-top: none;
        }
        .alert-box {
            background: #fff3cd;
            border: 1px solid #ffc107;
            border-radius: 5px;
            padding: 15px;
            margin-bottom: 20px;
        }
        .course-info {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 5px;
            padding: 20px;
            margin: 20px 0;
        }
        .course-title {
            color: #667eea;
            font-size: 18px;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .login-info {
            background: #e8f4fd;
            border: 1px solid #bee5eb;
            border-radius: 5px;
            padding: 15px;
            margin: 20px 0;
        }
        .btn {
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-decoration: none;
            padding: 12px 30px;
            border-radius: 5px;
            font-weight: bold;
            margin: 10px 0;
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 12px;
            padding: 20px;
            border: 1px solid #e0e0e0;
            border-top: none;
            border-radius: 0 0 10px 10px;
            background: #f9f9f9;
        }
        .password-box {
            background: #fff;
            border: 2px dashed #667eea;
            border-radius: 5px;
            padding: 15px;
            margin: 10px 0;
            text-align: center;
        }
        .password {
            font-family: monospace;
            font-size: 18px;
            color: #667eea;
            letter-spacing: 2px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Security Awareness Training</h1>
        <p>{{.CompanyName}}</p>
    </div>

    <div class="content">
        <div class="alert-box">
            <strong>Important Notice:</strong> You have been identified in a recent phishing simulation conducted by our security team. This training is required to help improve your security awareness.
        </div>

        <p>Dear User,</p>

        <p>You have fallen for a phishing simulation conducted by {{.CompanyName}}. Don't worry - this was a test, and no harm was done. However, this incident shows the importance of cybersecurity awareness.</p>

        <p>As part of your security awareness training, you have been enrolled in the following course:</p>

        <div class="course-info">
            <div class="course-title">{{.CourseName}}</div>
            <p>{{.CourseDescription}}</p>
        </div>

        <h3>How to Access Your Training</h3>

        <div class="login-info">
            <p><strong>Login Page:</strong> <a href="{{.LoginURL}}">{{.LoginURL}}</a></p>
            <p><strong>Username:</strong> {{.UserEmail}}</p>
            {{if .IsNewUser}}
            <div class="password-box">
                <p><strong>Your Temporary Password:</strong></p>
                <p class="password">{{.TempPassword}}</p>
                <p style="color: #dc3545; font-size: 12px;"><strong>Important:</strong> You must change this password on your first login.</p>
            </div>
            {{else}}
            <p><strong>Password:</strong> Use your existing password</p>
            {{end}}
        </div>

        <p style="text-align: center;">
            <a href="{{.LoginURL}}" class="btn">Start Your Training</a>
        </p>

        <h3>What You'll Learn</h3>
        <ul>
            <li>How to recognize phishing attempts</li>
            <li>Understanding security best practices</li>
            <li>Protecting yourself and our organization from cyber threats</li>
        </ul>

        <p>This training is mandatory and must be completed within the specified timeframe. If you have any questions, please contact your IT Security team.</p>

        <p>Best regards,<br>Security Awareness Team<br>{{.CompanyName}}</p>
    </div>

    <div class="footer">
        <p>This email was sent because you interacted with a security awareness campaign.</p>
        <p>This is an automated message from the Gophish security awareness platform.</p>
    </div>
</body>
</html>`
}
