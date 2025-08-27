package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
)

// DB returns the database connection
func DB() *gorm.DB {
	return db
}

// Course represents an e-learning course
type Course struct {
	Id           int64          `json:"id"`
	UserId       int64          `json:"-"`
	Name         string         `json:"name" sql:"not null"`
	Description  string         `json:"description"`
	CreatedDate  time.Time      `json:"created_date"`
	ModifiedDate time.Time      `json:"modified_date"`
	Modules      []CourseModule `json:"modules,omitempty"`
	Quizzes      []CourseQuiz   `json:"quizzes,omitempty"`
}

// CourseModule represents a section/module within a course with different types
type CourseModule struct {
	Id           int64     `json:"id"`
	CourseId     int64     `json:"course_id"`
	Name         string    `json:"name" sql:"not null"`
	Description  string    `json:"description"`
	ModuleType      string    `json:"module_type" sql:"not null;default:'html'"` // 'html', 'video', 'presentation', 'quiz'
	Content         string    `json:"content" sql:"type:text"`                   // HTML content for html modules or description/guidelines for others
	VideoURL        string    `json:"video_url"`                                 // Video URL for video modules
	VideoDuration   int       `json:"video_duration"`                           // Duration in seconds
	PresentationURL string    `json:"presentation_url"`                          // Google Presentation URL/ID for presentation modules
	QuizId          int64     `json:"quiz_id"`                                   // Linked quiz ID for quiz modules
	MustComplete bool      `json:"must_complete" sql:"default:true"`         // Whether completion is required
	MinTimeSpent int       `json:"min_time_spent" sql:"default:0"`           // Minimum time to spend (seconds)
	OrderIndex   int       `json:"order_index"`
	CreatedDate  time.Time `json:"created_date"`
	ModifiedDate time.Time `json:"modified_date"`
}

// CourseQuiz represents a quiz within a course with enhanced features
type CourseQuiz struct {
	Id              int64              `json:"id"`
	CourseId        int64              `json:"course_id"`
	ModuleId        int64              `json:"module_id,omitempty"`        // Optional: linked to a quiz module
	Name            string             `json:"name" sql:"not null"`
	Description     string             `json:"description"`
	OrderIndex      int                `json:"order_index"`
	PassingScore    int                `json:"passing_score"`               // Percentage needed to pass
	TimeLimit       int                `json:"time_limit"`                  // Time limit in minutes, 0 = no limit
	MaxAttempts     int                `json:"max_attempts" sql:"default:3"` // Maximum attempts allowed
	ShuffleQuestions bool              `json:"shuffle_questions" sql:"default:false"` // Randomize question order
	ShowResults     bool               `json:"show_results" sql:"default:true"`       // Show results after completion
	Questions       []CourseQuestion   `json:"questions,omitempty"`
	CreatedDate     time.Time          `json:"created_date"`
}

// CourseQuestion represents a question within a quiz with enhanced features
type CourseQuestion struct {
	Id           int64             `json:"id"`
	QuizId       int64             `json:"quiz_id"`
	Question     string            `json:"question" sql:"type:text;not null"`
	QuestionType string            `json:"question_type" sql:"default:'multiple_choice'"` // 'multiple_choice', 'true_false', 'text'
	Points       int               `json:"points" sql:"default:1"`                        // Points for this question
	Explanation  string            `json:"explanation" sql:"type:text"`                   // Explanation shown after answering
	OrderIndex   int               `json:"order_index"`
	Options      []QuestionOption  `json:"options,omitempty"`
}

// QuestionOption represents an answer option for a quiz question
type QuestionOption struct {
	Id         int64  `json:"id"`
	QuestionId int64  `json:"question_id"`
	Option     string `json:"option" sql:"not null"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index"`
}

// CourseEnrollment tracks user enrollment and progress in courses
type CourseEnrollment struct {
	Id               int64                  `json:"id"`
	UserId           int64                  `json:"user_id"`
	CourseId         int64                  `json:"course_id"`
	CampaignId       int64                  `json:"campaign_id,omitempty"`
	EnrolledDate     time.Time              `json:"enrolled_date"`
	StartedDate      *time.Time             `json:"started_date,omitempty"`
	CompletedDate    *time.Time             `json:"completed_date,omitempty"`
	Status           string                 `json:"status"`
	Progress         int                    `json:"progress"`
	Course           Course                 `json:"course,omitempty"`
	ModuleProgress   []ModuleProgress       `json:"module_progress,omitempty"`
	QuizAttempts     []QuizAttempt          `json:"quiz_attempts,omitempty"`
}

// ModuleProgress tracks user progress through individual modules with detailed tracking
type ModuleProgress struct {
	Id                 int64      `json:"id"`
	EnrollmentId       int64      `json:"enrollment_id"`
	ModuleId           int64      `json:"module_id"`
	StartedDate        *time.Time `json:"started_date,omitempty"`
	CompletedDate      *time.Time `json:"completed_date,omitempty"`
	TimeSpent          int        `json:"time_spent" sql:"default:0"`        // Time spent in seconds
	ProgressPercentage int        `json:"progress_percentage" sql:"default:0"` // 0-100
	Status             string     `json:"status" sql:"default:'not_started'"` // 'not_started', 'in_progress', 'completed'
	CompletionData     string     `json:"completion_data" sql:"type:text"`    // JSON data for module-specific info
}

// QuizAttempt represents a user's attempt at a quiz with enhanced tracking
type QuizAttempt struct {
	Id            int64          `json:"id"`
	EnrollmentId  int64          `json:"enrollment_id"`
	QuizId        int64          `json:"quiz_id"`
	AttemptNumber int            `json:"attempt_number" sql:"default:1"`    // Which attempt this is
	StartedDate   *time.Time     `json:"started_date,omitempty"`            // When attempt was started
	CompletedDate *time.Time     `json:"completed_date,omitempty"`          // When attempt was completed
	AttemptDate   time.Time      `json:"attempt_date"`                      // For backward compatibility
	TimeTaken     int            `json:"time_taken"`                        // Time taken in seconds
	Score         int            `json:"score"`                             // Score as percentage
	Passed        bool           `json:"passed"`                            // Whether the attempt passed
	Responses     []QuizResponse `json:"responses,omitempty"`
}

// QuizResponse represents a user's answer to a quiz question
type QuizResponse struct {
	Id           int64 `json:"id"`
	AttemptId    int64 `json:"attempt_id"`
	QuestionId   int64 `json:"question_id"`
	SelectedOptionId int64 `json:"selected_option_id"`
	IsCorrect    bool  `json:"is_correct"`
}

// Enrollment status constants
const (
	EnrollmentStatusEnrolled   = "enrolled"
	EnrollmentStatusInProgress = "in_progress"
	EnrollmentStatusCompleted  = "completed"
	EnrollmentStatusFailed     = "failed"
)

// Module type constants
const (
	ModuleTypeHTML         = "html"
	ModuleTypeVideo        = "video"
	ModuleTypePresentation = "presentation"
	ModuleTypeQuiz         = "quiz"
)

// Module progress status constants
const (
	ModuleStatusNotStarted = "not_started"
	ModuleStatusInProgress = "in_progress"
	ModuleStatusCompleted  = "completed"
)

// Question type constants
const (
	QuestionTypeMultipleChoice = "multiple_choice"
	QuestionTypeTrueFalse      = "true_false"
	QuestionTypeText           = "text"
)

// Course validation errors
var (
	ErrCourseNameNotSpecified = errors.New("Course name not specified")
	ErrCourseNotFound         = errors.New("Course not found")
	ErrEnrollmentNotFound     = errors.New("Enrollment not found")
	ErrQuizNotFound           = errors.New("Quiz not found")
	ErrQuestionNotFound       = errors.New("Question not found")
)

// Validate checks for required fields in course
func (c *Course) Validate() error {
	switch {
	case c.Name == "":
		return ErrCourseNameNotSpecified
	}
	return nil
}

// GetCourses returns all courses for a user
func GetCourses(uid int64) ([]Course, error) {
	courses := []Course{}
	err := db.Where("user_id = ?", uid).Find(&courses).Error
	if err != nil {
		log.Error(err)
		return courses, err
	}
	
	for i := range courses {
		err = db.Where("course_id = ?", courses[i].Id).Find(&courses[i].Modules).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
		err = db.Where("course_id = ?", courses[i].Id).Find(&courses[i].Quizzes).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
	}
	
	return courses, nil
}

// GetCourse returns a course by ID
func GetCourse(id int64, uid int64) (Course, error) {
	c := Course{}
	err := db.Where("id = ? AND user_id = ?", id, uid).Find(&c).Error
	if err != nil {
		return c, err
	}
	
	err = db.Where("course_id = ?", c.Id).Order("order_index").Find(&c.Modules).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
	}
	
	err = db.Where("course_id = ?", c.Id).Order("order_index").Find(&c.Quizzes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
	}
	
	for i := range c.Quizzes {
		err = db.Where("quiz_id = ?", c.Quizzes[i].Id).Order("order_index").Find(&c.Quizzes[i].Questions).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
		
		for j := range c.Quizzes[i].Questions {
			err = db.Where("question_id = ?", c.Quizzes[i].Questions[j].Id).Order("order_index").Find(&c.Quizzes[i].Questions[j].Options).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error(err)
			}
		}
	}
	
	return c, nil
}

// GetCourseByName returns a course by name
func GetCourseByName(name string, uid int64) (Course, error) {
	c := Course{}
	err := db.Where("name = ? AND user_id = ?", name, uid).Find(&c).Error
	return c, err
}

// PostCourse creates a new course
func PostCourse(c *Course, uid int64) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	
	c.UserId = uid
	c.CreatedDate = time.Now().UTC()
	c.ModifiedDate = c.CreatedDate
	
	return db.Save(c).Error
}

// PutCourse updates an existing course
func PutCourse(c *Course, uid int64) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	
	c.ModifiedDate = time.Now().UTC()
	err = db.Model(&Course{}).Where("id = ? AND user_id = ?", c.Id, uid).Updates(c).Error
	return err
}

// DeleteCourse deletes a course and all related data
func DeleteCourse(id int64, uid int64) error {
	tx := db.Begin()
	
	// Delete quiz responses
	tx.Exec("DELETE FROM quiz_responses WHERE attempt_id IN (SELECT id FROM quiz_attempts WHERE enrollment_id IN (SELECT id FROM course_enrollments WHERE course_id = ?))", id)
	
	// Delete quiz attempts
	tx.Exec("DELETE FROM quiz_attempts WHERE enrollment_id IN (SELECT id FROM course_enrollments WHERE course_id = ?)", id)
	
	// Delete module progress
	tx.Exec("DELETE FROM module_progress WHERE enrollment_id IN (SELECT id FROM course_enrollments WHERE course_id = ?)", id)
	
	// Delete enrollments
	tx.Where("course_id = ?", id).Delete(&CourseEnrollment{})
	
	// Delete question options
	tx.Exec("DELETE FROM question_options WHERE question_id IN (SELECT id FROM course_questions WHERE quiz_id IN (SELECT id FROM course_quizzes WHERE course_id = ?))", id)
	
	// Delete questions
	tx.Exec("DELETE FROM course_questions WHERE quiz_id IN (SELECT id FROM course_quizzes WHERE course_id = ?)", id)
	
	// Delete quizzes
	tx.Where("course_id = ?", id).Delete(&CourseQuiz{})
	
	// Delete modules
	tx.Where("course_id = ?", id).Delete(&CourseModule{})
	
	// Delete course
	err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&Course{}).Error
	
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// EnrollUser enrolls a user in a course
func EnrollUser(userId, courseId, campaignId int64) error {
	enrollment := CourseEnrollment{
		UserId:       userId,
		CourseId:     courseId,
		CampaignId:   campaignId,
		EnrolledDate: time.Now().UTC(),
		Status:       EnrollmentStatusEnrolled,
		Progress:     0,
	}
	
	// Check if already enrolled
	var count int64
	db.Model(&CourseEnrollment{}).Where("user_id = ? AND course_id = ?", userId, courseId).Count(&count)
	if count > 0 {
		return nil // Already enrolled
	}
	
	return db.Save(&enrollment).Error
}

// GetUserEnrollments returns all enrollments for a user
func GetUserEnrollments(userId int64) ([]CourseEnrollment, error) {
	enrollments := []CourseEnrollment{}
	err := db.Where("user_id = ?", userId).Find(&enrollments).Error
	if err != nil {
		return enrollments, err
	}
	
	for i := range enrollments {
		err = db.Where("id = ?", enrollments[i].CourseId).Find(&enrollments[i].Course).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		}
	}
	
	return enrollments, nil
}

// GetCourseEnrollment returns a specific enrollment
func GetCourseEnrollment(userId, courseId int64) (CourseEnrollment, error) {
	enrollment := CourseEnrollment{}
	err := db.Where("user_id = ? AND course_id = ?", userId, courseId).Find(&enrollment).Error
	if err != nil {
		return enrollment, err
	}
	
	err = db.Where("id = ?", enrollment.CourseId).Find(&enrollment.Course).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
	}
	
	return enrollment, nil
}

// UpdateEnrollmentProgress updates user progress in a course
func UpdateEnrollmentProgress(enrollmentId int64, progress int, status string) error {
	updates := map[string]interface{}{
		"progress": progress,
		"status":   status,
	}
	
	if status == EnrollmentStatusInProgress {
		updates["started_date"] = time.Now().UTC()
	} else if status == EnrollmentStatusCompleted {
		updates["completed_date"] = time.Now().UTC()
	}
	
	return db.Model(&CourseEnrollment{}).Where("id = ?", enrollmentId).Updates(updates).Error
}