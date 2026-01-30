package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
)

// ELearningSettings handles requests for the /api/elearning_settings/ endpoint
func (as *Server) ELearningSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		settings, err := models.GetELearningSettings()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		// Load all SMTP profiles for the dropdown
		uid := ctx.Get(r, "user_id").(int64)
		smtps, err := models.GetSMTPs(uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		// Return settings along with available SMTP profiles
		response := map[string]interface{}{
			"settings":      settings,
			"smtp_profiles": smtps,
		}
		JSONResponse(w, response, http.StatusOK)

	case "POST":
		settings := models.ELearningSettings{}
		err := json.NewDecoder(r.Body).Decode(&settings)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request data"}, http.StatusBadRequest)
			return
		}

		// Validate that the SMTP profile belongs to the current user
		if settings.SMTPId > 0 {
			uid := ctx.Get(r, "user_id").(int64)
			_, err := models.GetSMTP(settings.SMTPId, uid)
			if err != nil {
				JSONResponse(w, models.Response{Success: false, Message: "Invalid sending profile selected"}, http.StatusBadRequest)
				return
			}
		}

		err = models.PostELearningSettings(&settings)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}

		JSONResponse(w, models.Response{Success: true, Message: "E-Learning settings saved successfully"}, http.StatusOK)

	default:
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusMethodNotAllowed)
	}
}

// ELearningSettingsTestEmail handles requests for the /api/elearning_settings/test endpoint
func (as *Server) ELearningSettingsTestEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	// Parse the test email request
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		JSONResponse(w, models.Response{Success: false, Message: "Please provide a valid email address"}, http.StatusBadRequest)
		return
	}

	// Load the current settings
	settings, err := models.GetELearningSettings()
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	// Check if SMTP is configured
	if settings.SMTPId <= 0 || settings.SMTP.Id == 0 {
		JSONResponse(w, models.Response{Success: false, Message: "Please select a sending profile and save settings first"}, http.StatusBadRequest)
		return
	}

	// Create a sample course for the test email
	sampleCourse := models.Course{
		Name:        "Sample Security Awareness Course",
		Description: "This is a sample course description used for testing the enrollment notification email.",
	}

	// Send the test email using the template
	err = models.SendEnrollmentNotificationWithTemplate(req.Email, sampleCourse, true, "SamplePass123!")
	if err != nil {
		log.Errorf("Failed to send test email: %v", err)
		JSONResponse(w, models.Response{Success: false, Message: fmt.Sprintf("Failed to send test email: %v", err)}, http.StatusInternalServerError)
		return
	}

	JSONResponse(w, models.Response{Success: true, Message: fmt.Sprintf("Test email sent successfully to %s", req.Email)}, http.StatusOK)
}

// ELearningSettingsDefaultTemplate handles requests for the /api/elearning_settings/default_template endpoint
func (as *Server) ELearningSettingsDefaultTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	template := models.GetDefaultEmailTemplate()
	JSONResponse(w, map[string]string{"template": template}, http.StatusOK)
}

// ELearningSettingsPreviewTemplate handles requests for the /api/elearning_settings/preview endpoint
func (as *Server) ELearningSettingsPreviewTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	// Parse the preview request
	var req struct {
		Template string `json:"template"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Template == "" {
		JSONResponse(w, models.Response{Success: false, Message: "Please provide a template to preview"}, http.StatusBadRequest)
		return
	}

	// Load settings for company name and base URL
	settings, _ := models.GetELearningSettings()

	// Create sample data for preview
	data := models.EmailTemplateData{
		UserEmail:         "john.doe@example.com",
		TempPassword:      "TempPass123!",
		CourseName:        "Security Awareness Fundamentals",
		CourseDescription: "Learn how to identify and avoid common security threats like phishing, social engineering, and malware.",
		LoginURL:          fmt.Sprintf("%s/login", settings.BaseURL),
		CompanyName:       settings.CompanyName,
		IsNewUser:         true,
	}

	// Render the template
	tmpl, err := template.New("preview").Parse(req.Template)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: fmt.Sprintf("Template parsing error: %v", err)}, http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: fmt.Sprintf("Template execution error: %v", err)}, http.StatusBadRequest)
		return
	}

	JSONResponse(w, map[string]string{"preview": buf.String()}, http.StatusOK)
}
