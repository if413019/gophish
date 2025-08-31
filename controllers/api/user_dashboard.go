package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
)

// UserStats provides statistics for the user dashboard
func (as *Server) UserStats(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		user := ctx.Get(r, "user").(models.User)
		
		// Get user enrollment statistics
		stats, err := getUserDashboardStats(user.Id)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, stats, http.StatusOK)
	}
}

// UserActiveCourses returns courses the user is currently enrolled in
func (as *Server) UserActiveCourses(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		user := ctx.Get(r, "user").(models.User)
		
		courses, err := getUserActiveCourses(user.Id)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, courses, http.StatusOK)
	}
}

// UserActivity returns recent activity for the user
func (as *Server) UserActivity(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		user := ctx.Get(r, "user").(models.User)
		
		activities, err := getUserActivity(user.Id)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, activities, http.StatusOK)
	}
}

// UserAchievements returns user achievements
func (as *Server) UserAchievements(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		user := ctx.Get(r, "user").(models.User)
		
		achievements, err := getUserAchievements(user.Id)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, achievements, http.StatusOK)
	}
}

// UserDashboardStats represents the statistics shown on the user dashboard
type UserDashboardStats struct {
	PhishedCount      int    `json:"phished_count"`
	ReportedCount     int    `json:"reported_count"`
	AssignedCourses   int    `json:"assigned_courses"`
	CompletedCourses  int    `json:"completed_courses"`
	PhishedStatus     string `json:"phished_status"`
	ReportedStatus    string `json:"reported_status"`
	CourseStatus      string `json:"course_status"`
	CompletionStatus  string `json:"completion_status"`
}

// UserCourse represents a course from the user's perspective
type UserCourse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Progress    int    `json:"progress"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline,omitempty"`
}

// UserActivity represents a user activity event
type UserActivity struct {
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"created_date"`
}

// UserAchievement represents a user achievement
type UserAchievement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Unlocked    bool   `json:"unlocked"`
}

// getUserDashboardStats calculates statistics for the user dashboard
func getUserDashboardStats(userID int64) (UserDashboardStats, error) {
	stats := UserDashboardStats{}
	
	// Get user enrollments to count assigned courses
	enrollments, err := models.GetUserEnrollments(userID)
	if err != nil {
		return stats, err
	}
	
	stats.AssignedCourses = len(enrollments)
	
	// Calculate completed courses and course progress
	completedCourses := 0
	if stats.AssignedCourses > 0 {
		for _, enrollment := range enrollments {
			courseProgress := calculateCourseProgress(enrollment.Id, enrollment.CourseId)
			if courseProgress == 100 {
				completedCourses++
			}
		}
	}
	stats.CompletedCourses = completedCourses
	
	// Calculate phishing statistics from campaign results
	phishedCount, reportedCount := calculatePhishingStats(userID)
	stats.PhishedCount = phishedCount
	stats.ReportedCount = reportedCount
	
	// Status indicators
	if stats.PhishedCount > 0 {
		stats.PhishedStatus = "needs_attention"
	} else {
		stats.PhishedStatus = "good"
	}
	
	if stats.ReportedCount > 0 {
		stats.ReportedStatus = "good"
	} else {
		stats.ReportedStatus = "pending"
	}
	
	if stats.AssignedCourses > 0 {
		completionPercentage := (completedCourses * 100) / stats.AssignedCourses
		if completionPercentage >= 80 {
			stats.CourseStatus = "good"
		} else if completionPercentage >= 50 {
			stats.CourseStatus = "in_progress"
		} else {
			stats.CourseStatus = "needs_attention"
		}
		
		if completedCourses == stats.AssignedCourses {
			stats.CompletionStatus = "completed"
		} else if completedCourses > 0 {
			stats.CompletionStatus = "in_progress"
		} else {
			stats.CompletionStatus = "needs_attention"
		}
	}
	
	return stats, nil
}

// getUserActiveCourses gets courses the user is enrolled in
func getUserActiveCourses(userID int64) ([]UserCourse, error) {
	var courses []UserCourse
	
	log.Infof("Getting active courses for user ID: %d", userID)
	
	// Get user enrollments
	enrollments, err := models.GetUserEnrollments(userID)
	if err != nil {
		log.Errorf("Error getting user enrollments: %v", err)
		return courses, err
	}
	
	log.Infof("Found %d enrollments for user %d", len(enrollments), userID)
	
	// Get course details for each enrollment
	for i, enrollment := range enrollments {
		log.Infof("Processing enrollment %d: CourseID=%d, Status=%s", i, enrollment.CourseId, enrollment.Status)
		
		if enrollment.Status == "completed" {
			log.Infof("Skipping completed course %d", enrollment.CourseId)
			continue // Skip completed courses
		}
		
		// Get course details - try with admin user first (ID 1)
		course, err := models.GetCourse(enrollment.CourseId, 1)
		if err != nil {
			log.Errorf("Error getting course %d: %v", enrollment.CourseId, err)
			continue // Skip if can't get course details
		}
		
		log.Infof("Successfully loaded course: %s (ID: %d)", course.Name, course.Id)
		
		// Calculate actual progress based on completed modules
		actualProgress := calculateCourseProgress(enrollment.Id, course.Id)
		
		// Determine status based on progress
		status := enrollment.Status
		if actualProgress > 0 && status == "enrolled" {
			status = "in_progress"
		}
		if actualProgress == 100 {
			status = "completed"
		}
		
		userCourse := UserCourse{
			ID:          course.Id,
			Name:        course.Name,
			Description: course.Description,
			Progress:    actualProgress,
			Status:      status,
		}
		
		courses = append(courses, userCourse)
	}
	
	return courses, nil
}

// getUserActivity gets recent user activity
func getUserActivity(userID int64) ([]UserActivity, error) {
	var activities []UserActivity
	
	// Get enrollment activities
	enrollments, err := models.GetUserEnrollments(userID)
	if err != nil {
		return activities, err
	}
	
	// Add enrollment activities
	for _, enrollment := range enrollments {
		course, err := models.GetCourse(enrollment.CourseId, 1) // Use admin user ID
		if err != nil {
			log.Errorf("Error getting course %d for activity: %v", enrollment.CourseId, err)
			continue
		}
		
		activity := UserActivity{
			Type:        "enrollment",
			Title:       "Enrolled in Course",
			Description: "You were enrolled in " + course.Name,
			CreatedDate: enrollment.EnrolledDate,
		}
		activities = append(activities, activity)
	}
	
	// Add module completion activities
	for _, enrollment := range enrollments {
		moduleProgress, err := models.GetUserModuleProgress(enrollment.Id)
		if err != nil {
			log.Errorf("Error getting module progress for enrollment %d: %v", enrollment.Id, err)
			continue
		}
		
		course, err := models.GetCourse(enrollment.CourseId, 1)
		if err != nil {
			continue
		}
		
		// Add activity for each completed module
		for _, progress := range moduleProgress {
			if progress.CompletedDate != nil {
				// Find the module name
				var moduleName string
				for _, module := range course.Modules {
					if module.Id == progress.ModuleId {
						moduleName = module.Name
						break
					}
				}
				
				activity := UserActivity{
					Type:        "module_completion",
					Title:       "Completed Module",
					Description: fmt.Sprintf("You completed \"%s\" in course \"%s\"", moduleName, course.Name),
					CreatedDate: *progress.CompletedDate,
				}
				activities = append(activities, activity)
			} else if progress.StartedDate != nil {
				// Add activity for started modules
				var moduleName string
				for _, module := range course.Modules {
					if module.Id == progress.ModuleId {
						moduleName = module.Name
						break
					}
				}
				
				activity := UserActivity{
					Type:        "module_start",
					Title:       "Started Module",
					Description: fmt.Sprintf("You started \"%s\" in course \"%s\"", moduleName, course.Name),
					CreatedDate: *progress.StartedDate,
				}
				activities = append(activities, activity)
			}
		}
	}
	
	// Sort activities by date (most recent first)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].CreatedDate.After(activities[j].CreatedDate)
	})
	
	return activities, nil
}

// getUserAchievements gets user achievements
func getUserAchievements(userID int64) ([]UserAchievement, error) {
	var achievements []UserAchievement
	
	// Get user stats for achievement calculation
	stats, err := getUserDashboardStats(userID)
	if err != nil {
		return achievements, err
	}
	
	// Define phishing-focused achievements
	achievementDefs := []UserAchievement{
		{
			Title:       "Lesson Learned",
			Description: "Complete your first security awareness course after being phished",
			Icon:        "fa-graduation-cap",
			Unlocked:    stats.PhishedCount > 0 && stats.CompletedCourses > 0,
		},
		{
			Title:       "Security Defender",
			Description: "Report your first suspicious email",
			Icon:        "fa-shield",
			Unlocked:    stats.ReportedCount > 0,
		},
		{
			Title:       "Phishing Hunter",
			Description: "Report 5 or more phishing attempts",
			Icon:        "fa-flag",
			Unlocked:    stats.ReportedCount >= 5,
		},
		{
			Title:       "Victim to Victor",
			Description: "Complete training and avoid falling for phishing twice",
			Icon:        "fa-trophy",
			Unlocked:    stats.PhishedCount <= 1 && stats.CompletedCourses > 0,
		},
		{
			Title:       "Security Champion",
			Description: "Complete all assigned training courses",
			Icon:        "fa-star",
			Unlocked:    stats.CompletedCourses == stats.AssignedCourses && stats.AssignedCourses > 0,
		},
		{
			Title:       "Eagle Eye",
			Description: "Report suspicious emails without falling victim",
			Icon:        "fa-eye",
			Unlocked:    stats.ReportedCount > 0 && stats.PhishedCount == 0,
		},
	}
	
	return achievementDefs, nil
}

// calculateCourseProgress calculates the completion percentage for a course based on completed modules
func calculateCourseProgress(enrollmentId, courseId int64) int {
	// Get total number of modules in the course
	course, err := models.GetCourse(courseId, 1) // Use admin user ID
	if err != nil {
		log.Errorf("Error getting course for progress calculation: %v", err)
		return 0
	}
	
	totalModules := len(course.Modules)
	if totalModules == 0 {
		return 0
	}
	
	// Get completed modules
	moduleProgress, err := models.GetUserModuleProgress(enrollmentId)
	if err != nil {
		log.Errorf("Error getting module progress: %v", err)
		return 0
	}
	
	// Count completed modules
	completedModules := 0
	for _, progress := range moduleProgress {
		if progress.CompletedDate != nil {
			completedModules++
		}
	}
	
	// Calculate percentage
	percentage := (completedModules * 100) / totalModules
	
	log.Infof("Course progress: %d/%d modules completed = %d%%", completedModules, totalModules, percentage)
	
	return percentage
}

// calculatePhishingStats calculates phishing-related statistics from campaign results
func calculatePhishingStats(userID int64) (int, int) {
	// Get the user to find their email/username
	user, err := models.GetUser(userID)
	if err != nil {
		log.Errorf("Error getting user for phishing stats: %v", err)
		return 0, 0
	}
	
	// Query campaign results to count phishing interactions
	// In Gophish, results are linked by email address in the BaseRecipient
	phishedCount := 0
	reportedCount := 0
	
	// Count how many times user clicked links or submitted data
	err = models.DB().Model(&models.Result{}).
		Where("email = ? AND (status = ? OR status = ?)", user.Username, models.EventClicked, models.EventDataSubmit).
		Count(&phishedCount).Error
	
	if err != nil {
		log.Errorf("Error counting phished results: %v", err)
	}
	
	// Count how many times user reported phishing emails
	err = models.DB().Model(&models.Result{}).
		Where("email = ? AND reported = ?", user.Username, true).
		Count(&reportedCount).Error
	
	if err != nil {
		log.Errorf("Error counting reported results: %v", err)
	}
	
	log.Infof("Phishing stats for user %s (ID: %d): phished=%d, reported=%d", user.Username, userID, phishedCount, reportedCount)
	
	return phishedCount, reportedCount
}

// UserModuleProgress returns module progress for a user's course enrollment
func (as *Server) UserModuleProgress(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["id"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Get user enrollment for this course
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			log.Error("Error getting user enrollment: ", err)
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Get module progress
		progress, err := models.GetUserModuleProgress(enrollment.Id)
		if err != nil {
			log.Error("Error getting module progress: ", err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, progress, http.StatusOK)
	}
}

// UserStartModule marks a module as started for a user
func (as *Server) UserStartModule(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		moduleId, err := strconv.ParseInt(vars["moduleId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid module ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Get user enrollment for this course
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			log.Error("Error getting user enrollment: ", err)
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Start the module
		err = models.StartModule(enrollment.Id, moduleId)
		if err != nil {
			log.Error("Error starting module: ", err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, models.Response{Success: true, Message: "Module started"}, http.StatusOK)
	}
}

// UserCompleteModule marks a module as completed for a user
func (as *Server) UserCompleteModule(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		moduleId, err := strconv.ParseInt(vars["moduleId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid module ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Get user enrollment for this course
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			log.Error("Error getting user enrollment: ", err)
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Complete the module
		err = models.CompleteModule(enrollment.Id, moduleId)
		if err != nil {
			log.Error("Error completing module: ", err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Check if the course is now completed and create timeline event if needed
		courseCompleted, err := models.CheckAndUpdateCourseCompletion(enrollment.Id)
		if err != nil {
			log.Error("Error checking course completion: ", err)
			// Don't fail the request, just log the error
		} else if courseCompleted && enrollment.CampaignId > 0 {
			// Create a course completion timeline event
			// Reload enrollment to get updated data
			updatedEnrollment, err := models.GetCourseEnrollment(user.Id, courseId)
			if err == nil {
				err = models.CreateCourseCompletionEvent(updatedEnrollment, enrollment.CampaignId)
				if err != nil {
					log.Error("Error creating course completion event: ", err)
				} else {
					log.Infof("Created course completion event for user %s in campaign %d", user.Username, enrollment.CampaignId)
				}
			}
		}
		
		JSONResponse(w, models.Response{Success: true, Message: "Module completed"}, http.StatusOK)
	}
}

// UserGetQuizQuestions retrieves quiz questions for a user (randomized if configured)
func (as *Server) UserGetQuizQuestions(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		quizId, err := strconv.ParseInt(vars["quizId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid quiz ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Verify user enrollment
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Get quiz with questions
		quiz, err := models.GetQuizWithQuestions(quizId, courseId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Quiz not found"}, http.StatusNotFound)
			return
		}
		
		// Check if user can access this quiz (previous modules completed)
		canAccess, err := models.CanAccessQuiz(enrollment.Id, quiz.ModuleId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		if !canAccess {
			JSONResponse(w, models.Response{Success: false, Message: "Complete previous modules before taking this quiz"}, http.StatusForbidden)
			return
		}
		
		// Randomize questions if configured
		questions := quiz.Questions
		if quiz.ShuffleQuestions {
			models.ShuffleQuestions(&questions)
		}
		
		// Remove correct answers from response (for security)
		for i := range questions {
			for j := range questions[i].Options {
				questions[i].Options[j].IsCorrect = false
			}
		}
		
		JSONResponse(w, map[string]interface{}{
			"quiz":      quiz,
			"questions": questions,
		}, http.StatusOK)
	}
}

// UserStartQuizAPI starts a new quiz attempt
func (as *Server) UserStartQuizAPI(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		quizId, err := strconv.ParseInt(vars["quizId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid quiz ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Verify user enrollment
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Create quiz attempt
		attempt, err := models.StartQuizAttempt(enrollment.Id, quizId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, map[string]interface{}{
			"attempt_id": attempt.Id,
			"started_at": attempt.StartedDate,
		}, http.StatusOK)
	}
}

// UserSubmitQuizAPI submits quiz responses and calculates score
func (as *Server) UserSubmitQuizAPI(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		quizId, err := strconv.ParseInt(vars["quizId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid quiz ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Parse submission data
		type QuizSubmission struct {
			AttemptId int64                    `json:"attempt_id"`
			Responses []models.QuizResponseSubmission `json:"responses"`
		}
		
		var submission QuizSubmission
		err = json.NewDecoder(r.Body).Decode(&submission)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid submission format"}, http.StatusBadRequest)
			return
		}
		
		// Verify user enrollment
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Process quiz submission
		result, err := models.SubmitQuizAttempt(submission.AttemptId, submission.Responses, enrollment.Id, quizId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// If quiz passed and this is a module quiz, mark module complete
		if result.Passed {
			quiz, err := models.GetQuizWithQuestions(quizId, courseId)
			if err == nil && quiz.ModuleId > 0 {
				err = models.CompleteModule(enrollment.Id, quiz.ModuleId)
				if err != nil {
					log.Error("Error completing module after quiz success: ", err)
				} else {
					// Check for course completion
					courseCompleted, err := models.CheckAndUpdateCourseCompletion(enrollment.Id)
					if err != nil {
						log.Error("Error checking course completion: ", err)
					} else if courseCompleted && enrollment.CampaignId > 0 {
						updatedEnrollment, err := models.GetCourseEnrollment(user.Id, courseId)
						if err == nil {
							err = models.CreateCourseCompletionEvent(updatedEnrollment, enrollment.CampaignId)
							if err != nil {
								log.Error("Error creating course completion event: ", err)
							}
						}
					}
				}
			}
		}
		
		JSONResponse(w, result, http.StatusOK)
	}
}

// UserGetQuizAttempts retrieves user's quiz attempts and scores
func (as *Server) UserGetQuizAttempts(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		vars := mux.Vars(r)
		courseId, err := strconv.ParseInt(vars["courseId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid course ID"}, http.StatusBadRequest)
			return
		}
		
		quizId, err := strconv.ParseInt(vars["quizId"], 0, 64)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid quiz ID"}, http.StatusBadRequest)
			return
		}
		
		user := ctx.Get(r, "user").(models.User)
		
		// Verify user enrollment
		enrollment, err := models.GetCourseEnrollment(user.Id, courseId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
			return
		}
		
		// Get quiz attempts
		attempts, err := models.GetQuizAttempts(enrollment.Id, quizId)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, attempts, http.StatusOK)
	}
}