package models

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
	"github.com/oschwald/maxminddb-golang"
)

type mmCity struct {
	GeoPoint mmGeoPoint `maxminddb:"location"`
}

type mmGeoPoint struct {
	Latitude  float64 `maxminddb:"latitude"`
	Longitude float64 `maxminddb:"longitude"`
}

// Result contains the fields for a result object,
// which is a representation of a target in a campaign.
type Result struct {
	Id           int64     `json:"-"`
	CampaignId   int64     `json:"-"`
	UserId       int64     `json:"-"`
	RId          string    `json:"id"`
	Status       string    `json:"status" sql:"not null"`
	IP           string    `json:"ip"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	SendDate     time.Time `json:"send_date"`
	Reported     bool      `json:"reported" sql:"not null"`
	ModifiedDate time.Time `json:"modified_date"`
	BaseRecipient
}

func (r *Result) createEvent(status string, details interface{}) (*Event, error) {
	e := &Event{Email: r.Email, Message: status}
	if details != nil {
		dj, err := json.Marshal(details)
		if err != nil {
			return nil, err
		}
		e.Details = string(dj)
	}
	AddEvent(e, r.CampaignId)
	return e, nil
}

// HandleEmailSent updates a Result to indicate that the email has been
// successfully sent to the remote SMTP server
func (r *Result) HandleEmailSent() error {
	event, err := r.createEvent(EventSent, nil)
	if err != nil {
		return err
	}
	r.SendDate = event.Time
	r.Status = EventSent
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleEmailError updates a Result to indicate that there was an error when
// attempting to send the email to the remote SMTP server.
func (r *Result) HandleEmailError(err error) error {
	event, err := r.createEvent(EventSendingError, EventError{Error: err.Error()})
	if err != nil {
		return err
	}
	r.Status = Error
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleEmailBackoff updates a Result to indicate that the email received a
// temporary error and needs to be retried
func (r *Result) HandleEmailBackoff(err error, sendDate time.Time) error {
	event, err := r.createEvent(EventSendingError, EventError{Error: err.Error()})
	if err != nil {
		return err
	}
	r.Status = StatusRetry
	r.SendDate = sendDate
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleEmailOpened updates a Result in the case where the recipient opened the
// email.
func (r *Result) HandleEmailOpened(details EventDetails) error {
	event, err := r.createEvent(EventOpened, details)
	if err != nil {
		return err
	}
	// Don't update the status if the user already clicked the link
	// or submitted data to the campaign
	if r.Status == EventClicked || r.Status == EventDataSubmit {
		return nil
	}
	r.Status = EventOpened
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleClickedLink updates a Result in the case where the recipient clicked
// the link in an email.
func (r *Result) HandleClickedLink(details EventDetails) error {
	event, err := r.createEvent(EventClicked, details)
	if err != nil {
		return err
	}
	// Don't update the status if the user has already submitted data via the
	// landing page form.
	if r.Status == EventDataSubmit {
		return nil
	}
	r.Status = EventClicked
	r.ModifiedDate = event.Time
	
	// Check if auto-enrollment should occur
	campaign := Campaign{}
	err = db.Where("id = ?", r.CampaignId).First(&campaign).Error
	if err == nil && campaign.CourseId > 0 {
		// Get or create user for enrollment
		// For now, we'll use the email as a basic user identifier
		// In a full implementation, this would integrate with the existing user system
		user, err := r.getOrCreateUserForEnrollment()
		if err != nil {
			log.Errorf("Failed to get/create user for enrollment: %v", err)
		} else {
			// Enroll the user in the course
			err = EnrollUser(user.Id, campaign.CourseId, r.CampaignId)
			if err != nil {
				log.Errorf("Failed to enroll user in course: %v", err)
			} else {
				log.Infof("User %s automatically enrolled in course %d from campaign %d", r.Email, campaign.CourseId, r.CampaignId)
				// Send enrollment notification email
				log.Infof("Sending enrollment notification to user: %s", r.Email)
				err = r.sendEnrollmentNotification(user, campaign)
				if err != nil {
					log.Errorf("Failed to send enrollment notification to %s: %v", r.Email, err)
				} else {
					log.Infof("Enrollment notification sent successfully to %s", r.Email)
				}
			}
		}
	}
	
	return db.Save(r).Error
}

// HandleFormSubmit updates a Result in the case where the recipient submitted
// credentials to the form on a Landing Page.
func (r *Result) HandleFormSubmit(details EventDetails) error {
	event, err := r.createEvent(EventDataSubmit, details)
	if err != nil {
		return err
	}
	r.Status = EventDataSubmit
	r.ModifiedDate = event.Time
	
	// Check if auto-enrollment should occur (form submission also triggers enrollment)
	campaign := Campaign{}
	err = db.Where("id = ?", r.CampaignId).First(&campaign).Error
	if err == nil && campaign.CourseId > 0 {
		// Get or create user for enrollment
		user, err := r.getOrCreateUserForEnrollment()
		if err != nil {
			log.Errorf("Failed to get/create user for enrollment: %v", err)
		} else {
			// Enroll the user in the course
			err = EnrollUser(user.Id, campaign.CourseId, r.CampaignId)
			if err != nil {
				log.Errorf("Failed to enroll user in course: %v", err)
			} else {
				log.Infof("User %s automatically enrolled in course %d from campaign %d (form submit)", r.Email, campaign.CourseId, r.CampaignId)
				// Send enrollment notification email
				log.Infof("Sending enrollment notification to user: %s (form submit)", r.Email)
				err = r.sendEnrollmentNotification(user, campaign)
				if err != nil {
					log.Errorf("Failed to send enrollment notification to %s (form submit): %v", r.Email, err)
				} else {
					log.Infof("Enrollment notification sent successfully to %s (form submit)", r.Email)
				}
			}
		}
	}
	
	return db.Save(r).Error
}

// getOrCreateUserForEnrollment gets or creates a user for course enrollment
// This creates a basic user record for course enrollment purposes
func (r *Result) getOrCreateUserForEnrollment() (User, error) {
	// Try to find existing user by email
	var user User
	err := db.Where("username = ?", r.Email).First(&user).Error
	if err == nil {
		// User already exists
		return user, nil
	}
	
	if err != gorm.ErrRecordNotFound {
		// Some other error occurred
		return user, err
	}
	
	// User doesn't exist, create a basic user record for enrollment
	// This would be a limited user account specifically for course access
	user = User{
		Username: r.Email,
		// Set a placeholder API key - this user won't use API functionality
		ApiKey:   generateSecureKey(),
		// These users don't have admin access to the system
		RoleID: 2, // Assuming role ID 2 is a basic user role
	}
	
	err = db.Save(&user).Error
	if err != nil {
		return user, err
	}
	
	return user, nil
}

// sendEnrollmentNotification sends an email to notify the user about their enrollment
func (r *Result) sendEnrollmentNotification(user User, campaign Campaign) error {
	log.Infof("Starting enrollment notification process for %s", r.Email)
	log.Infof("Campaign CourseId: %d, Campaign UserId: %d", campaign.CourseId, campaign.UserId)
	
	// Load the course details
	course, err := GetCourse(campaign.CourseId, campaign.UserId)
	if err != nil {
		log.Errorf("Failed to load course details for CourseId %d, UserId %d: %v", campaign.CourseId, campaign.UserId, err)
		return fmt.Errorf("failed to load course details: %v", err)
	}
	log.Infof("Successfully loaded course: %s (ID: %d)", course.Name, course.Id)

	// Check if this is a new user or needs password reset (needs password setup)
	isNewUser := user.Hash == "" || user.PasswordChangeRequired
	log.Infof("User %s isNewUser: %t (hash empty: %t, password_change_required: %t)", r.Email, isNewUser, user.Hash == "", user.PasswordChangeRequired)
	var tempPassword string
	
	if isNewUser {
		log.Infof("Generating temporary password for new user: %s", r.Email)
		// Generate a temporary password
		tempPassword, err = GenerateSecurePassword(12)
		if err != nil {
			return fmt.Errorf("failed to generate temporary password: %v", err)
		}
		
		// Set up the user for first login
		err = UpdateUserPasswordForFirstLogin(&user, tempPassword)
		if err != nil {
			return fmt.Errorf("failed to set up user password: %v", err)
		}
		log.Infof("Temporary password set up for user: %s", r.Email)
	} else {
		log.Infof("User %s already has password, skipping temporary password setup", r.Email)
	}

	// Send the notification email
	log.Infof("Calling SendEnrollmentNotification for %s", r.Email)
	err = SendEnrollmentNotification(r.Email, course, isNewUser, tempPassword)
	if err != nil {
		log.Errorf("SendEnrollmentNotification failed for %s: %v", r.Email, err)
		return err
	}
	log.Infof("SendEnrollmentNotification completed successfully for %s", r.Email)
	return nil
}

// HandleEmailReport updates a Result in the case where they report a simulated
// phishing email using the HTTP handler.
func (r *Result) HandleEmailReport(details EventDetails) error {
	event, err := r.createEvent(EventReported, details)
	if err != nil {
		return err
	}
	r.Reported = true
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// UpdateGeo updates the latitude and longitude of the result in
// the database given an IP address
func (r *Result) UpdateGeo(addr string) error {
	// Open a connection to the maxmind db
	mmdb, err := maxminddb.Open("static/db/geolite2-city.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer mmdb.Close()
	ip := net.ParseIP(addr)
	var city mmCity
	// Get the record
	err = mmdb.Lookup(ip, &city)
	if err != nil {
		return err
	}
	// Update the database with the record information
	r.IP = addr
	r.Latitude = city.GeoPoint.Latitude
	r.Longitude = city.GeoPoint.Longitude
	return db.Save(r).Error
}

func generateResultId() (string, error) {
	const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	k := make([]byte, 7)
	for i := range k {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNum))))
		if err != nil {
			return "", err
		}
		k[i] = alphaNum[idx.Int64()]
	}
	return string(k), nil
}

// GenerateId generates a unique key to represent the result
// in the database
func (r *Result) GenerateId(tx *gorm.DB) error {
	// Keep trying until we generate a unique key (shouldn't take more than one or two iterations)
	for {
		rid, err := generateResultId()
		if err != nil {
			return err
		}
		r.RId = rid
		err = tx.Table("results").Where("r_id=?", r.RId).First(&Result{}).Error
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	return nil
}

// GetResult returns the Result object from the database
// given the ResultId
func GetResult(rid string) (Result, error) {
	r := Result{}
	err := db.Where("r_id=?", rid).First(&r).Error
	return r, err
}
