package api

import (
	"net/http"
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
	PhishingTests     int    `json:"phishing_tests"`
	AssignedCourses   int    `json:"assigned_courses"`
	CompletionRate    int    `json:"completion_rate"`
	LearningStreak    int    `json:"learning_streak"`
	PhishingStatus    string `json:"phishing_status"`
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
	
	// Calculate completion rate based on actual module progress
	if stats.AssignedCourses > 0 {
		totalProgress := 0
		for _, enrollment := range enrollments {
			courseProgress := calculateCourseProgress(enrollment.Id, enrollment.CourseId)
			totalProgress += courseProgress
		}
		stats.CompletionRate = totalProgress / stats.AssignedCourses
	}
	
	// Placeholder values for other stats
	stats.PhishingTests = 0
	stats.LearningStreak = 1
	
	// Status indicators
	if stats.AssignedCourses > 0 {
		if stats.CompletionRate >= 80 {
			stats.CourseStatus = "good"
		} else if stats.CompletionRate >= 50 {
			stats.CourseStatus = "in_progress"
		} else {
			stats.CourseStatus = "needs_attention"
		}
		
		if stats.CompletionRate == 100 {
			stats.CompletionStatus = "completed"
		} else if stats.CompletionRate >= 50 {
			stats.CompletionStatus = "in_progress"
		} else {
			stats.CompletionStatus = "needs_attention"
		}
	}
	
	stats.PhishingStatus = "good" // Placeholder
	
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
		course, err := models.GetCourse(enrollment.CourseId, userID)
		if err != nil {
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
	
	// TODO: Add phishing test activities when we have proper models for them
	
	// Sort activities by date (most recent first)
	// For now, we'll keep them as is since we only have enrollments
	
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
	
	// Define achievements
	achievementDefs := []UserAchievement{
		{
			Title:       "First Steps",
			Description: "Complete your first security awareness course",
			Icon:        "fa-baby",
			Unlocked:    stats.CompletionRate > 0,
		},
		{
			Title:       "Security Aware",
			Description: "Complete 3 security awareness courses",
			Icon:        "fa-shield",
			Unlocked:    stats.AssignedCourses >= 3 && stats.CompletionRate >= 50,
		},
		{
			Title:       "Phishing Expert",
			Description: "Recognize and report 5 phishing attempts",
			Icon:        "fa-eye",
			Unlocked:    false, // Placeholder - need to track reported phishing
		},
		{
			Title:       "Perfect Score",
			Description: "Complete all assigned courses with 100% score",
			Icon:        "fa-star",
			Unlocked:    stats.CompletionRate == 100 && stats.AssignedCourses > 0,
		},
		{
			Title:       "Streak Master",
			Description: "Maintain a 30-day learning streak",
			Icon:        "fa-fire",
			Unlocked:    stats.LearningStreak >= 30,
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
		
		JSONResponse(w, models.Response{Success: true, Message: "Module completed"}, http.StatusOK)
	}
}