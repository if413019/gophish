package models

import (
	"encoding/json"
	"errors"
	"math/rand"
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
	ContentItems []ContentItem  `json:"content_items,omitempty" sql:"-"` // Unified view of modules and quizzes
}

// ContentItem represents a unified view of modules and quizzes for display
type ContentItem struct {
	Type       string        `json:"type"`        // "module" or "quiz"
	OrderIndex int           `json:"order_index"`
	Module     *CourseModule `json:"module,omitempty"`
	Quiz       *CourseQuiz   `json:"quiz,omitempty"`
}

// ContentItemSlice is a slice of ContentItem that implements sort.Interface
type ContentItemSlice []ContentItem

func (c ContentItemSlice) Len() int           { return len(c) }
func (c ContentItemSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ContentItemSlice) Less(i, j int) bool { return c[i].OrderIndex < c[j].OrderIndex }

// GetContentItems returns a unified, sorted list of modules and quizzes
func (c *Course) GetContentItems() []ContentItem {
	items := make([]ContentItem, 0, len(c.Modules)+len(c.Quizzes))

	for i := range c.Modules {
		items = append(items, ContentItem{
			Type:       "module",
			OrderIndex: c.Modules[i].OrderIndex,
			Module:     &c.Modules[i],
		})
	}

	for i := range c.Quizzes {
		items = append(items, ContentItem{
			Type:       "quiz",
			OrderIndex: c.Quizzes[i].OrderIndex,
			Quiz:       &c.Quizzes[i],
		})
	}

	// Sort by OrderIndex
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].OrderIndex < items[i].OrderIndex {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	return items
}

// CourseModule represents a section/module within a course with different types
type CourseModule struct {
	Id           int64     `json:"id"`
	CourseId     int64     `json:"course_id"`
	Name         string    `json:"name" sql:"not null"`
	Description  string    `json:"description"`
	ModuleType      string    `json:"module_type" sql:"not null;default:'html'"` // 'html', 'video', 'presentation', 'quiz'
	Content         string    `json:"content" sql:"type:text"`                   // HTML content for html modules or description/guidelines for others
	VideoURL        string    `json:"video_url"`                                 // Legacy: External video URL
	VideoDuration   int       `json:"video_duration"`                           // Duration in seconds
	VideoFilePath   string    `json:"video_file_path"`                          // Local video file path
	PresentationURL string    `json:"presentation_url"`                          // Legacy: External presentation URL
	PresentationFilePath string `json:"presentation_file_path"`                 // Local presentation file path
	QuizId          int64     `json:"quiz_id"`                                   // Linked quiz ID for quiz modules
	OriginalFilename string   `json:"original_filename"`                        // Original uploaded filename
	FileSize        int64     `json:"file_size"`                               // File size in bytes
	MimeType        string    `json:"mime_type"`                               // MIME type of uploaded file
	UploadedDate    time.Time `json:"uploaded_date"`                           // Upload timestamp
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
	
	// Link quizzes to their corresponding modules
	for i := range c.Modules {
		if c.Modules[i].ModuleType == ModuleTypeQuiz {
			// Find the quiz that belongs to this module
			for j := range c.Quizzes {
				if c.Quizzes[j].ModuleId == c.Modules[i].Id {
					c.Modules[i].QuizId = c.Quizzes[j].Id
					break
				}
			}
		}
	}

	// Populate unified content items view
	c.ContentItems = c.GetContentItems()

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

	// Store modules separately before creating course (GORM may reset nested slices)
	modules := make([]CourseModule, len(c.Modules))
	copy(modules, c.Modules)

	// Store quizzes separately before creating course (GORM may reset nested slices)
	quizzes := make([]CourseQuiz, len(c.Quizzes))
	for i, quiz := range c.Quizzes {
		quizzes[i] = quiz
		// Deep copy questions
		quizzes[i].Questions = make([]CourseQuestion, len(quiz.Questions))
		for j, q := range quiz.Questions {
			quizzes[i].Questions[j] = q
			// Deep copy options
			quizzes[i].Questions[j].Options = make([]QuestionOption, len(q.Options))
			copy(quizzes[i].Questions[j].Options, q.Options)
		}
	}

	log.Info("PostCourse: Received ", len(modules), " modules and ", len(quizzes), " quizzes")
	for i, quiz := range quizzes {
		log.Info("Quiz ", i, ": ", quiz.Name, " has ", len(quiz.Questions), " questions")
		for j, q := range quiz.Questions {
			log.Info("  Question ", j, ": ", q.Question, " has ", len(q.Options), " options")
		}
	}

	tx := db.Begin()

	// Clear nested slices before creating course to avoid GORM issues
	c.Quizzes = nil
	c.Modules = nil

	// Save the course first to get its ID
	err = tx.Create(c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	log.Info("PostCourse: Course created with ID ", c.Id)

	// Save modules
	for i := range modules {
		modules[i].CourseId = c.Id
		modules[i].CreatedDate = time.Now().UTC()
		modules[i].ModifiedDate = modules[i].CreatedDate
		err = tx.Create(&modules[i]).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Restore modules to course
	c.Modules = modules

	// Save quizzes with questions and options from our copy
	for i := range quizzes {
		quizzes[i].CourseId = c.Id
		quizzes[i].CreatedDate = time.Now().UTC()

		// Store questions before creating quiz
		questions := quizzes[i].Questions
		quizzes[i].Questions = nil

		err = tx.Create(&quizzes[i]).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		log.Info("PostCourse: Quiz created with ID ", quizzes[i].Id, ", now saving ", len(questions), " questions")

		// Save questions
		for j := range questions {
			questions[j].QuizId = quizzes[i].Id

			// Store options before creating question
			options := questions[j].Options
			questions[j].Options = nil

			err = tx.Create(&questions[j]).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			log.Info("PostCourse: Question created with ID ", questions[j].Id, ", now saving ", len(options), " options")

			// Save options
			for k := range options {
				options[k].QuestionId = questions[j].Id
				err = tx.Create(&options[k]).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			// Restore options to question
			questions[j].Options = options
		}

		// Restore questions to quiz
		quizzes[i].Questions = questions
	}

	// Restore quizzes to course
	c.Quizzes = quizzes

	tx.Commit()
	return nil
}

// PutCourse updates an existing course
func PutCourse(c *Course, uid int64) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	
	tx := db.Begin()
	
	// Update course basic info
	c.ModifiedDate = time.Now().UTC()
	err = tx.Model(&Course{}).Where("id = ? AND user_id = ?", c.Id, uid).Updates(c).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// Update or create modules
	for i := range c.Modules {
		c.Modules[i].CourseId = c.Id
		c.Modules[i].ModifiedDate = time.Now().UTC()
		
		// Debug logging
		log.Info("Processing module: ", c.Modules[i].Name, " (ID: ", c.Modules[i].Id, ", Type: ", c.Modules[i].ModuleType, ")")
		log.Info("Presentation file path: '", c.Modules[i].PresentationFilePath, "'")
		log.Info("Video file path: '", c.Modules[i].VideoFilePath, "'")
		log.Info("Original filename: '", c.Modules[i].OriginalFilename, "'")
		
		if c.Modules[i].Id == 0 {
			// New module - create it
			c.Modules[i].CreatedDate = time.Now().UTC()
			err = tx.Create(&c.Modules[i]).Error
			log.Info("Created new module with ID: ", c.Modules[i].Id)
		} else {
			// Existing module - update it, preserve created_date if not provided
			if c.Modules[i].CreatedDate.IsZero() {
				// Frontend didn't provide created_date, fetch it from database
				existingModule := CourseModule{}
				err = tx.Where("id = ?", c.Modules[i].Id).First(&existingModule).Error
				if err == nil {
					c.Modules[i].CreatedDate = existingModule.CreatedDate
				}
			}
			err = tx.Save(&c.Modules[i]).Error
			log.Info("Updated existing module ID: ", c.Modules[i].Id)
		}
		
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	
	// Handle module deletion - delete any modules that exist in DB but not in current submission
	var existingModuleIds []int64
	err = tx.Model(&CourseModule{}).Where("course_id = ?", c.Id).Pluck("id", &existingModuleIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// Build list of module IDs from current submission
	var submittedModuleIds []int64
	for _, module := range c.Modules {
		if module.Id != 0 {
			submittedModuleIds = append(submittedModuleIds, module.Id)
		}
	}
	
	// Find modules to delete (exist in DB but not in submission)
	var modulesToDelete []int64
	for _, existingId := range existingModuleIds {
		found := false
		for _, submittedId := range submittedModuleIds {
			if existingId == submittedId {
				found = true
				break
			}
		}
		if !found {
			modulesToDelete = append(modulesToDelete, existingId)
		}
	}
	
	// Delete the modules that were removed
	if len(modulesToDelete) > 0 {
		log.Info("Deleting removed modules: ", modulesToDelete)
		err = tx.Where("id IN (?)", modulesToDelete).Delete(&CourseModule{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update or create quizzes
	for i := range c.Quizzes {
		c.Quizzes[i].CourseId = c.Id

		log.Info("Processing quiz: ", c.Quizzes[i].Name, " (ID: ", c.Quizzes[i].Id, ")")

		if c.Quizzes[i].Id == 0 {
			// New quiz - create it
			c.Quizzes[i].CreatedDate = time.Now().UTC()
			err = tx.Create(&c.Quizzes[i]).Error
			log.Info("Created new quiz with ID: ", c.Quizzes[i].Id)
		} else {
			// Existing quiz - update it
			err = tx.Save(&c.Quizzes[i]).Error
			log.Info("Updated existing quiz ID: ", c.Quizzes[i].Id)
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		// Handle questions for this quiz
		for j := range c.Quizzes[i].Questions {
			c.Quizzes[i].Questions[j].QuizId = c.Quizzes[i].Id

			if c.Quizzes[i].Questions[j].Id == 0 {
				err = tx.Create(&c.Quizzes[i].Questions[j]).Error
			} else {
				err = tx.Save(&c.Quizzes[i].Questions[j]).Error
			}

			if err != nil {
				tx.Rollback()
				return err
			}

			// Handle options for this question
			for k := range c.Quizzes[i].Questions[j].Options {
				c.Quizzes[i].Questions[j].Options[k].QuestionId = c.Quizzes[i].Questions[j].Id

				if c.Quizzes[i].Questions[j].Options[k].Id == 0 {
					err = tx.Create(&c.Quizzes[i].Questions[j].Options[k]).Error
				} else {
					err = tx.Save(&c.Quizzes[i].Questions[j].Options[k]).Error
				}

				if err != nil {
					tx.Rollback()
					return err
				}
			}

			// Delete removed options
			var existingOptionIds []int64
			err = tx.Model(&QuestionOption{}).Where("question_id = ?", c.Quizzes[i].Questions[j].Id).Pluck("id", &existingOptionIds).Error
			if err == nil && len(existingOptionIds) > 0 {
				var submittedOptionIds []int64
				for _, opt := range c.Quizzes[i].Questions[j].Options {
					if opt.Id != 0 {
						submittedOptionIds = append(submittedOptionIds, opt.Id)
					}
				}
				for _, existingId := range existingOptionIds {
					found := false
					for _, submittedId := range submittedOptionIds {
						if existingId == submittedId {
							found = true
							break
						}
					}
					if !found {
						tx.Where("id = ?", existingId).Delete(&QuestionOption{})
					}
				}
			}
		}

		// Delete removed questions
		var existingQuestionIds []int64
		err = tx.Model(&CourseQuestion{}).Where("quiz_id = ?", c.Quizzes[i].Id).Pluck("id", &existingQuestionIds).Error
		if err == nil && len(existingQuestionIds) > 0 {
			var submittedQuestionIds []int64
			for _, q := range c.Quizzes[i].Questions {
				if q.Id != 0 {
					submittedQuestionIds = append(submittedQuestionIds, q.Id)
				}
			}
			for _, existingId := range existingQuestionIds {
				found := false
				for _, submittedId := range submittedQuestionIds {
					if existingId == submittedId {
						found = true
						break
					}
				}
				if !found {
					// Delete options first, then question
					tx.Where("question_id = ?", existingId).Delete(&QuestionOption{})
					tx.Where("id = ?", existingId).Delete(&CourseQuestion{})
				}
			}
		}
	}

	// Handle quiz deletion - delete any quizzes that exist in DB but not in current submission
	var existingQuizIds []int64
	err = tx.Model(&CourseQuiz{}).Where("course_id = ?", c.Id).Pluck("id", &existingQuizIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Build list of quiz IDs from current submission
	var submittedQuizIds []int64
	for _, quiz := range c.Quizzes {
		if quiz.Id != 0 {
			submittedQuizIds = append(submittedQuizIds, quiz.Id)
		}
	}

	// Find quizzes to delete (exist in DB but not in submission)
	var quizzesToDelete []int64
	for _, existingId := range existingQuizIds {
		found := false
		for _, submittedId := range submittedQuizIds {
			if existingId == submittedId {
				found = true
				break
			}
		}
		if !found {
			quizzesToDelete = append(quizzesToDelete, existingId)
		}
	}

	// Delete the quizzes that were removed (with their questions and options)
	if len(quizzesToDelete) > 0 {
		log.Info("Deleting removed quizzes: ", quizzesToDelete)
		// Delete options first
		tx.Exec("DELETE FROM question_options WHERE question_id IN (SELECT id FROM course_questions WHERE quiz_id IN (?))", quizzesToDelete)
		// Delete questions
		tx.Exec("DELETE FROM course_questions WHERE quiz_id IN (?)", quizzesToDelete)
		// Delete quizzes
		err = tx.Where("id IN (?)", quizzesToDelete).Delete(&CourseQuiz{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
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

// GetUserModuleProgress gets module progress for a user's enrollment
func GetUserModuleProgress(enrollmentId int64) ([]ModuleProgress, error) {
	var progress []ModuleProgress
	err := db.Where("enrollment_id = ?", enrollmentId).Find(&progress).Error
	return progress, err
}

// StartModule marks a module as started for a user
func StartModule(enrollmentId, moduleId int64) error {
	// Check if progress already exists
	var existing ModuleProgress
	err := db.Where("enrollment_id = ? AND module_id = ?", enrollmentId, moduleId).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new progress record
		progress := ModuleProgress{
			EnrollmentId: enrollmentId,
			ModuleId:     moduleId,
			StartedDate:  &[]time.Time{time.Now().UTC()}[0],
		}
		return db.Create(&progress).Error
	} else if err != nil {
		return err
	}
	
	// Update existing record if not already started
	if existing.StartedDate == nil {
		now := time.Now().UTC()
		existing.StartedDate = &now
		return db.Save(&existing).Error
	}
	
	return nil
}

// CompleteModule marks a module as completed for a user
func CompleteModule(enrollmentId, moduleId int64) error {
	var progress ModuleProgress
	err := db.Where("enrollment_id = ? AND module_id = ?", enrollmentId, moduleId).First(&progress).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new progress record with both started and completed dates
		now := time.Now().UTC()
		progress = ModuleProgress{
			EnrollmentId:  enrollmentId,
			ModuleId:      moduleId,
			StartedDate:   &now,
			CompletedDate: &now,
		}
		return db.Create(&progress).Error
	} else if err != nil {
		return err
	}
	
	// Update completion date
	now := time.Now().UTC()
	progress.CompletedDate = &now
	if progress.StartedDate == nil {
		progress.StartedDate = &now
	}
	
	return db.Save(&progress).Error
}

// CheckAndUpdateCourseCompletion checks if a course is fully completed and updates enrollment status
// Returns true if course was just completed (to trigger events)
func CheckAndUpdateCourseCompletion(enrollmentId int64) (bool, error) {
	// Get the enrollment
	var enrollment CourseEnrollment
	err := db.Where("id = ?", enrollmentId).First(&enrollment).Error
	if err != nil {
		return false, err
	}
	
	// If already completed, no need to check again
	if enrollment.Status == EnrollmentStatusCompleted {
		return false, nil
	}
	
	// Get the course with modules
	var course Course
	err = db.Where("id = ?", enrollment.CourseId).Preload("Modules").First(&course).Error
	if err != nil {
		return false, err
	}
	
	// Get all module progress for this enrollment
	var moduleProgresses []ModuleProgress
	err = db.Where("enrollment_id = ?", enrollmentId).Find(&moduleProgresses).Error
	if err != nil {
		return false, err
	}
	
	// Create a map of completed modules
	completedModules := make(map[int64]bool)
	for _, progress := range moduleProgresses {
		if progress.CompletedDate != nil {
			completedModules[progress.ModuleId] = true
		}
	}
	
	// Count total modules and required modules
	totalModules := len(course.Modules)
	completedCount := len(completedModules)
	requiredModules := 0
	
	for _, module := range course.Modules {
		if module.MustComplete {
			requiredModules++
		}
	}
	
	// Check if course is completed
	// Course is completed if either:
	// 1. All modules are completed, OR
	// 2. All required modules are completed (if there are any)
	isCompleted := false
	if requiredModules > 0 {
		// Check if all required modules are completed
		requiredCompleted := 0
		for _, module := range course.Modules {
			if module.MustComplete && completedModules[module.Id] {
				requiredCompleted++
			}
		}
		isCompleted = (requiredCompleted == requiredModules)
	} else {
		// All modules must be completed if no specific requirements
		isCompleted = (completedCount == totalModules)
	}
	
	if isCompleted {
		// Calculate completion percentage
		completionPercentage := 100
		if totalModules > 0 {
			completionPercentage = (completedCount * 100) / totalModules
		}
		
		// Update enrollment status
		err = UpdateEnrollmentProgress(enrollmentId, completionPercentage, EnrollmentStatusCompleted)
		if err != nil {
			return false, err
		}
		
		return true, nil // Course was just completed
	}
	
	// Update progress percentage even if not completed
	completionPercentage := 0
	if totalModules > 0 {
		completionPercentage = (completedCount * 100) / totalModules
	}
	
	status := EnrollmentStatusInProgress
	if completedCount == 0 {
		status = EnrollmentStatusEnrolled
	}
	
	err = UpdateEnrollmentProgress(enrollmentId, completionPercentage, status)
	return false, err
}

// CreateCourseCompletionEvent creates a timeline event for course completion
func CreateCourseCompletionEvent(enrollment CourseEnrollment, campaignId int64) error {
	// Get course details
	var course Course
	err := db.Where("id = ?", enrollment.CourseId).First(&course).Error
	if err != nil {
		return err
	}
	
	// Get user details
	var user User
	err = db.Where("id = ?", enrollment.UserId).First(&user).Error
	if err != nil {
		return err
	}
	
	// Create the event
	event := Event{
		CampaignId: campaignId,
		Email:      user.Username, // Username is the email for enrolled users
		Time:       time.Now().UTC(),
		Message:    EventCourseCompleted,
	}
	
	// Add details
	details := map[string]interface{}{
		"course_id":     course.Id,
		"course_name":   course.Name,
		"enrollment_id": enrollment.Id,
		"progress":      enrollment.Progress,
	}
	
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return err
	}
	event.Details = string(detailsJSON)
	
	return AddEvent(&event, campaignId)
}

// QuizResponseSubmission represents a quiz response from the frontend
type QuizResponseSubmission struct {
	QuestionId       int64 `json:"question_id"`
	SelectedOptionId int64 `json:"selected_option_id"`
}

// QuizAttemptResult represents the result of a quiz submission
type QuizAttemptResult struct {
	AttemptId       int64        `json:"attempt_id"`
	Score           int          `json:"score"`
	Passed          bool         `json:"passed"`
	PassingScore    int          `json:"passing_score"`
	TimeTaken       int          `json:"time_taken"`
	CompletedDate   *time.Time   `json:"completed_date"`
	CanRetake       bool         `json:"can_retake"`
	AttemptsLeft    int          `json:"attempts_left"`
	Responses       []QuizResponse `json:"responses,omitempty"`
}

// GetQuizWithQuestions retrieves a quiz with its questions and options
func GetQuizWithQuestions(quizId, courseId int64) (CourseQuiz, error) {
	var quiz CourseQuiz
	err := db.Where("id = ? AND course_id = ?", quizId, courseId).First(&quiz).Error
	if err != nil {
		return quiz, err
	}
	
	// Load questions with options
	err = db.Where("quiz_id = ?", quizId).Order("order_index").Find(&quiz.Questions).Error
	if err != nil {
		return quiz, err
	}
	
	// Load options for each question
	for i := range quiz.Questions {
		err = db.Where("question_id = ?", quiz.Questions[i].Id).Find(&quiz.Questions[i].Options).Error
		if err != nil {
			return quiz, err
		}
	}
	
	return quiz, nil
}

// CanAccessQuiz checks if a user can access a quiz based on module prerequisites
func CanAccessQuiz(enrollmentId, moduleId int64) (bool, error) {
	if moduleId == 0 {
		// No prerequisite module, can access
		return true, nil
	}
	
	// Get all modules for this course to check order
	var enrollment CourseEnrollment
	err := db.Where("id = ?", enrollmentId).First(&enrollment).Error
	if err != nil {
		return false, err
	}
	
	var currentModule CourseModule
	err = db.Where("id = ?", moduleId).First(&currentModule).Error
	if err != nil {
		return false, err
	}
	
	// Get all modules before this one
	var previousModules []CourseModule
	err = db.Where("course_id = ? AND order_index < ?", enrollment.CourseId, currentModule.OrderIndex).Find(&previousModules).Error
	if err != nil {
		return false, err
	}
	
	// Check if all previous modules are completed
	for _, module := range previousModules {
		var progress ModuleProgress
		err = db.Where("enrollment_id = ? AND module_id = ?", enrollmentId, module.Id).First(&progress).Error
		if err == gorm.ErrRecordNotFound || progress.CompletedDate == nil {
			// Previous module not completed
			return false, nil
		}
	}
	
	return true, nil
}

// ShuffleQuestions randomizes the order of questions
func ShuffleQuestions(questions *[]CourseQuestion) {
	if len(*questions) <= 1 {
		return
	}
	
	// Simple Fisher-Yates shuffle
	for i := len(*questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		(*questions)[i], (*questions)[j] = (*questions)[j], (*questions)[i]
	}
}

// StartQuizAttempt creates a new quiz attempt for a user
func StartQuizAttempt(enrollmentId, quizId int64) (QuizAttempt, error) {
	// Check if user has reached max attempts
	var existingAttempts []QuizAttempt
	err := db.Where("enrollment_id = ? AND quiz_id = ?", enrollmentId, quizId).Find(&existingAttempts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return QuizAttempt{}, err
	}
	
	// Get quiz to check max attempts
	var quiz CourseQuiz
	err = db.Where("id = ?", quizId).First(&quiz).Error
	if err != nil {
		return QuizAttempt{}, err
	}
	
	if len(existingAttempts) >= quiz.MaxAttempts {
		return QuizAttempt{}, errors.New("Maximum attempts reached for this quiz")
	}
	
	// Create new attempt
	now := time.Now().UTC()
	attempt := QuizAttempt{
		EnrollmentId:  enrollmentId,
		QuizId:        quizId,
		AttemptNumber: len(existingAttempts) + 1,
		StartedDate:   &now,
		AttemptDate:   now,
	}
	
	err = db.Create(&attempt).Error
	return attempt, err
}

// SubmitQuizAttempt processes quiz responses and calculates score
func SubmitQuizAttempt(attemptId int64, responses []QuizResponseSubmission, enrollmentId, quizId int64) (QuizAttemptResult, error) {
	// Get the attempt
	var attempt QuizAttempt
	err := db.Where("id = ? AND enrollment_id = ?", attemptId, enrollmentId).First(&attempt).Error
	if err != nil {
		return QuizAttemptResult{}, err
	}
	
	// Get quiz details
	quiz, err := GetQuizWithQuestions(quizId, 0) // courseId not needed for validation
	if err != nil {
		return QuizAttemptResult{}, err
	}
	
	// Calculate score
	totalQuestions := len(quiz.Questions)
	correctAnswers := 0
	var quizResponses []QuizResponse
	
	// Create a map of questions for easy lookup
	questionMap := make(map[int64]CourseQuestion)
	for _, question := range quiz.Questions {
		questionMap[question.Id] = question
	}
	
	// Process each response
	for _, response := range responses {
		question, exists := questionMap[response.QuestionId]
		if !exists {
			continue // Skip invalid questions
		}
		
		// Find the correct answer
		isCorrect := false
		for _, option := range question.Options {
			if option.Id == response.SelectedOptionId && option.IsCorrect {
				isCorrect = true
				correctAnswers++
				break
			}
		}
		
		// Create quiz response record
		quizResponse := QuizResponse{
			AttemptId:        attemptId,
			QuestionId:       response.QuestionId,
			SelectedOptionId: response.SelectedOptionId,
			IsCorrect:        isCorrect,
		}
		
		err = db.Create(&quizResponse).Error
		if err != nil {
			return QuizAttemptResult{}, err
		}
		
		quizResponses = append(quizResponses, quizResponse)
	}
	
	// Calculate final score
	score := 0
	if totalQuestions > 0 {
		score = (correctAnswers * 100) / totalQuestions
	}
	
	passed := score >= quiz.PassingScore
	now := time.Now().UTC()
	
	// Calculate time taken
	timeTaken := 0
	if attempt.StartedDate != nil {
		timeTaken = int(now.Sub(*attempt.StartedDate).Seconds())
	}
	
	// Update attempt with results
	attempt.CompletedDate = &now
	attempt.Score = score
	attempt.Passed = passed
	attempt.TimeTaken = timeTaken
	
	err = db.Save(&attempt).Error
	if err != nil {
		return QuizAttemptResult{}, err
	}
	
	// Check attempts left
	var totalAttempts int64
	db.Model(&QuizAttempt{}).Where("enrollment_id = ? AND quiz_id = ?", enrollmentId, quizId).Count(&totalAttempts)
	attemptsLeft := quiz.MaxAttempts - int(totalAttempts)
	if attemptsLeft < 0 {
		attemptsLeft = 0
	}
	
	return QuizAttemptResult{
		AttemptId:     attemptId,
		Score:         score,
		Passed:        passed,
		PassingScore:  quiz.PassingScore,
		TimeTaken:     timeTaken,
		CompletedDate: &now,
		CanRetake:     !passed && attemptsLeft > 0,
		AttemptsLeft:  attemptsLeft,
		Responses:     quizResponses,
	}, nil
}

// GetQuizAttempts retrieves all attempts for a user's quiz
func GetQuizAttempts(enrollmentId, quizId int64) ([]QuizAttempt, error) {
	var attempts []QuizAttempt
	err := db.Where("enrollment_id = ? AND quiz_id = ?", enrollmentId, quizId).
		Order("attempt_number desc").Find(&attempts).Error
	return attempts, err
}