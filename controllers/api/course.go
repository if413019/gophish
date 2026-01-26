package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

// Courses handles the requests for the /api/courses/* endpoints
func (as *Server) Courses(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		as.getCourses(w, r)
	case r.Method == "POST":
		as.postCourse(w, r)
	}
}

// Course handles the requests for the /api/courses/:id endpoints
func (as *Server) Course(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	switch {
	case r.Method == "GET":
		as.getCourse(w, r, id)
	case r.Method == "PUT":
		as.putCourse(w, r, id)
	case r.Method == "DELETE":
		as.deleteCourse(w, r, id)
	}
}

// getCourses returns a list of courses if requested via GET
func (as *Server) getCourses(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)
	courses, err := models.GetCourses(u.Id)
	if err != nil {
		log.Error(err)
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, courses, http.StatusOK)
}

// getCourse returns a course by ID if requested via GET
func (as *Server) getCourse(w http.ResponseWriter, r *http.Request, id int64) {
	u := ctx.Get(r, "user").(models.User)
	c, err := models.GetCourse(id, u.Id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Course not found"}, http.StatusNotFound)
		return
	}
	JSONResponse(w, c, http.StatusOK)
}

// postCourse handles the creation of a new course
func (as *Server) postCourse(w http.ResponseWriter, r *http.Request) {
	c := models.Course{}
	u := ctx.Get(r, "user").(models.User)
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
		return
	}
	err = models.PostCourse(&c, u.Id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	JSONResponse(w, c, http.StatusCreated)
}

// putCourse handles the updating of an existing course
func (as *Server) putCourse(w http.ResponseWriter, r *http.Request, id int64) {
	c := models.Course{}
	u := ctx.Get(r, "user").(models.User)
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
		return
	}
	c.Id = id
	err = models.PutCourse(&c, u.Id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}
	JSONResponse(w, c, http.StatusOK)
}

// deleteCourse handles the deletion of a course
func (as *Server) deleteCourse(w http.ResponseWriter, r *http.Request, id int64) {
	u := ctx.Get(r, "user").(models.User)
	err := models.DeleteCourse(id, u.Id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, models.Response{Success: true, Message: "Course deleted successfully"}, http.StatusOK)
}

// CourseModules handles the requests for course modules
func (as *Server) CourseModules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, _ := strconv.ParseInt(vars["id"], 0, 64)
	u := ctx.Get(r, "user").(models.User)

	switch {
	case r.Method == "POST":
		module := models.CourseModule{}
		err := json.NewDecoder(r.Body).Decode(&module)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		
		// Verify course ownership
		_, err = models.GetCourse(courseId, u.Id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Course not found"}, http.StatusNotFound)
			return
		}
		
		module.CourseId = courseId
		module.CreatedDate = time.Now().UTC()
		
		err = models.DB().Save(&module).Error
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		JSONResponse(w, module, http.StatusCreated)
	}
}

// CourseQuizzes handles the requests for course quizzes
func (as *Server) CourseQuizzes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, _ := strconv.ParseInt(vars["id"], 0, 64)
	u := ctx.Get(r, "user").(models.User)

	switch {
	case r.Method == "POST":
		quiz := models.CourseQuiz{}
		err := json.NewDecoder(r.Body).Decode(&quiz)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		
		// Verify course ownership
		_, err = models.GetCourse(courseId, u.Id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Course not found"}, http.StatusNotFound)
			return
		}
		
		quiz.CourseId = courseId
		quiz.CreatedDate = time.Now().UTC()
		
		tx := models.DB().Begin()
		err = tx.Save(&quiz).Error
		if err != nil {
			tx.Rollback()
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		
		// Save questions and options
		for i := range quiz.Questions {
			quiz.Questions[i].QuizId = quiz.Id
			err = tx.Save(&quiz.Questions[i]).Error
			if err != nil {
				tx.Rollback()
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
				return
			}
			
			for j := range quiz.Questions[i].Options {
				quiz.Questions[i].Options[j].QuestionId = quiz.Questions[i].Id
				err = tx.Save(&quiz.Questions[i].Options[j]).Error
				if err != nil {
					tx.Rollback()
					JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
					return
				}
			}
		}
		
		tx.Commit()
		JSONResponse(w, quiz, http.StatusCreated)
	}
}

// UserEnrollments handles user course enrollments
func (as *Server) UserEnrollments(w http.ResponseWriter, r *http.Request) {
	u := ctx.Get(r, "user").(models.User)
	
	switch {
	case r.Method == "GET":
		enrollments, err := models.GetUserEnrollments(u.Id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, enrollments, http.StatusOK)
		
	case r.Method == "POST":
		var req struct {
			CourseId int64 `json:"course_id"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		
		err = models.EnrollUser(u.Id, req.CourseId, 0)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		
		JSONResponse(w, models.Response{Success: true, Message: "Successfully enrolled in course"}, http.StatusOK)
	}
}

// SubmitQuiz handles quiz submissions
func (as *Server) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	enrollmentId, _ := strconv.ParseInt(vars["enrollment_id"], 0, 64)
	quizId, _ := strconv.ParseInt(vars["quiz_id"], 0, 64)
	u := ctx.Get(r, "user").(models.User)

	var submission struct {
		Responses []struct {
			QuestionId       int64 `json:"question_id"`
			SelectedOptionId int64 `json:"selected_option_id"`
		} `json:"responses"`
	}

	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
		return
	}

	// Verify enrollment belongs to user
	enrollment := models.CourseEnrollment{}
	err = models.DB().Where("id = ? AND user_id = ?", enrollmentId, u.Id).First(&enrollment).Error
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Enrollment not found"}, http.StatusNotFound)
		return
	}

	// Create quiz attempt
	attempt := models.QuizAttempt{
		EnrollmentId: enrollmentId,
		QuizId:       quizId,
		AttemptDate:  time.Now().UTC(),
		Score:        0,
		Passed:       false,
	}

	tx := models.DB().Begin()
	err = tx.Save(&attempt).Error
	if err != nil {
		tx.Rollback()
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	// Process responses and calculate score
	correctAnswers := 0
	totalQuestions := len(submission.Responses)

	for _, resp := range submission.Responses {
		// Get correct answer
		option := models.QuestionOption{}
		err = tx.Where("id = ?", resp.SelectedOptionId).First(&option).Error
		if err != nil {
			tx.Rollback()
			JSONResponse(w, models.Response{Success: false, Message: "Invalid option selected"}, http.StatusBadRequest)
			return
		}

		isCorrect := option.IsCorrect
		if isCorrect {
			correctAnswers++
		}

		// Save response
		response := models.QuizResponse{
			AttemptId:        attempt.Id,
			QuestionId:       resp.QuestionId,
			SelectedOptionId: resp.SelectedOptionId,
			IsCorrect:        isCorrect,
		}
		err = tx.Save(&response).Error
		if err != nil {
			tx.Rollback()
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	// Calculate final score and passing status
	score := (correctAnswers * 100) / totalQuestions
	
	// Get quiz passing score
	quiz := models.CourseQuiz{}
	err = tx.Where("id = ?", quizId).First(&quiz).Error
	if err != nil {
		tx.Rollback()
		JSONResponse(w, models.Response{Success: false, Message: "Quiz not found"}, http.StatusNotFound)
		return
	}

	passed := score >= quiz.PassingScore
	attempt.Score = score
	attempt.Passed = passed

	err = tx.Save(&attempt).Error
	if err != nil {
		tx.Rollback()
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	result := struct {
		Score           int  `json:"score"`
		Passed          bool `json:"passed"`
		RequiredScore   int  `json:"required_score"`
		CorrectAnswers  int  `json:"correct_answers"`
		TotalQuestions  int  `json:"total_questions"`
	}{
		Score:          score,
		Passed:         passed,
		RequiredScore:  quiz.PassingScore,
		CorrectAnswers: correctAnswers,
		TotalQuestions: totalQuestions,
	}

	JSONResponse(w, result, http.StatusOK)
}

// GetQuizQuestions handles getting quiz questions
func (as *Server) GetQuizQuestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quizId, _ := strconv.ParseInt(vars["quiz_id"], 0, 64)
	
	// Get quiz with questions
	quiz, err := models.GetQuizWithQuestions(quizId, 0) // courseId not needed for questions
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Quiz not found"}, http.StatusNotFound)
		return
	}
	
	// Return questions
	JSONResponse(w, map[string]interface{}{"questions": quiz.Questions}, http.StatusOK)
}