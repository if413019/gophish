var courses = []
var course = {}

function save(idx){
    course = courses[idx]
    $("#modal").modal('show')
}

function dismiss(){
    $("#modal").modal('hide')
    $("#name").val("")
    $("#description").val("")
    $("#modules-container").empty()
    $("#quizzes-container").empty()
}

function newCourse(){
    course = {}
    $("#modalLabel").text("New Course")
    $("#name").val("")
    $("#description").val("")  
    $("#modules-container").empty()
    $("#quizzes-container").empty()
    $("#modal").modal('show')
}

function deleteCourse(idx){
    if (confirm("Delete " + courses[idx].name + "?")){
        api.courseId.delete(courses[idx].id)
        .success(function(data){
            successFlash(data.message)
            location.reload()
        })
        .error(function(data){
            modalError(data.responseJSON.message)
        })
    }
}

function addModule() {
    var moduleHtml = $("#module-template").html()
    var moduleCount = $("#modules-container .module-panel").length
    var moduleElement = $(moduleHtml)
    
    // Handle remove module
    moduleElement.find('.remove-module').click(function() {
        $(this).closest('.module-panel').remove()
    })
    
    // Handle module type change
    moduleElement.find('.module-type').change(function() {
        toggleModuleFields($(this))
    })
    
    // Load available quizzes for quiz selection
    loadQuizzesForModule(moduleElement)
    
    $("#modules-container").append(moduleElement)
    
    // Trigger initial field visibility based on default selection
    toggleModuleFields(moduleElement.find('.module-type'))
}

function toggleModuleFields(moduleTypeSelect) {
    var modulePanel = moduleTypeSelect.closest('.module-panel')
    var moduleType = moduleTypeSelect.val()
    
    console.log('Toggling module fields for type:', moduleType)
    
    // Hide all type-specific fields first
    modulePanel.find('.module-video-fields').hide()
    modulePanel.find('.module-presentation-fields').hide() 
    modulePanel.find('.module-quiz-fields').hide()
    
    // Show appropriate fields based on type
    switch(moduleType) {
        case 'video':
            console.log('Showing video fields')
            modulePanel.find('.module-video-fields').show()
            break
        case 'presentation':
            console.log('Showing presentation fields')
            modulePanel.find('.module-presentation-fields').show()
            break
        case 'quiz':
            console.log('Showing quiz fields')
            modulePanel.find('.module-quiz-fields').show()
            break
        case 'html':
        default:
            console.log('HTML type selected - no extra fields needed')
            // HTML is default, no extra fields needed
            break
    }
}

function loadQuizzesForModule(moduleElement) {
    // Load available quizzes from current form
    var quizSelect = moduleElement.find('.module-quiz-id')
    quizSelect.empty().append('<option value="">Select an existing quiz...</option>')
    
    // Get quizzes from current form (from quizzes tab)
    $("#quizzes-container .quiz-panel").each(function(index) {
        var quizName = $(this).find('.quiz-name').val() || ('Quiz ' + (index + 1))
        quizSelect.append('<option value="' + index + '">' + quizName + '</option>')
    })
}

function addQuiz() {
    var quizHtml = $("#quiz-template").html()
    var quizElement = $(quizHtml)
    
    quizElement.find('.remove-quiz').click(function() {
        $(this).closest('.quiz-panel').remove()
        // Refresh all module quiz dropdowns after removing a quiz
        refreshAllModuleQuizDropdowns()
    })
    
    quizElement.find('.add-question').click(function() {
        addQuestion($(this).closest('.quiz-panel'))
    })
    
    // Listen for quiz name changes to update module dropdowns
    quizElement.find('.quiz-name').on('input', function() {
        refreshAllModuleQuizDropdowns()
    })
    
    $("#quizzes-container").append(quizElement)
    
    // Refresh all module quiz dropdowns after adding a new quiz
    refreshAllModuleQuizDropdowns()
}

function refreshAllModuleQuizDropdowns() {
    $("#modules-container .module-panel").each(function() {
        loadQuizzesForModule($(this))
    })
}

function addQuestion(quizPanel) {
    var questionHtml = $("#question-template").html()
    var questionElement = $(questionHtml)
    
    questionElement.find('.remove-question').click(function() {
        $(this).closest('.question-panel').remove()
    })
    
    questionElement.find('.add-option').click(function() {
        addOption($(this).closest('.question-panel'))
    })
    
    quizPanel.find('.questions-container').append(questionElement)
}

function addOption(questionPanel) {
    var optionHtml = $("#option-template").html()
    var optionElement = $(optionHtml)
    var questionId = questionPanel.index()
    
    // Set unique name for radio buttons within this question
    optionElement.find('.correct-option').attr('name', 'correct-option-' + questionId + '-' + Date.now())
    
    optionElement.find('.remove-option').click(function() {
        $(this).closest('.option-group').remove()
    })
    
    questionPanel.find('.options-container').append(optionElement)
}

function serializeCourse() {
    var courseData = {
        name: $("#name").val(),
        description: $("#description").val(),
        modules: [],
        quizzes: []
    }
    
    // Serialize modules
    $("#modules-container .module-panel").each(function(index) {
        var modulePanel = $(this)
        var moduleType = modulePanel.find('.module-type').val()
        
        var module = {
            name: modulePanel.find('.module-name').val(),
            description: modulePanel.find('.module-description').val(),
            content: modulePanel.find('.module-content').val(),
            module_type: moduleType,
            must_complete: modulePanel.find('.module-must-complete').is(':checked'),
            min_time_spent: parseInt(modulePanel.find('.module-min-time').val()) || 0,
            order_index: index
        }
        
        // Add type-specific fields
        if (moduleType === 'video') {
            module.video_url = modulePanel.find('.module-video-url').val()
        } else if (moduleType === 'presentation') {
            module.presentation_url = modulePanel.find('.module-presentation-url').val()
        } else if (moduleType === 'quiz') {
            var selectedQuizIndex = modulePanel.find('.module-quiz-id').val()
            if (selectedQuizIndex !== '') {
                module.quiz_id = parseInt(selectedQuizIndex)
            }
        }
        
        courseData.modules.push(module)
    })
    
    // Serialize quizzes
    $("#quizzes-container .quiz-panel").each(function(quizIndex) {
        var quiz = {
            name: $(this).find('.quiz-name').val(),
            description: $(this).find('.quiz-description').val(),
            passing_score: parseInt($(this).find('.quiz-passing-score').val()) || 70,
            order_index: quizIndex,
            questions: []
        }
        
        $(this).find('.question-panel').each(function(questionIndex) {
            var question = {
                question: $(this).find('.question-text').val(),
                order_index: questionIndex,
                options: []
            }
            
            $(this).find('.option-group').each(function(optionIndex) {
                var isCorrect = $(this).find('.correct-option').is(':checked')
                var option = {
                    option: $(this).find('.option-text').val(),
                    is_correct: isCorrect,
                    order_index: optionIndex
                }
                question.options.push(option)
            })
            
            quiz.questions.push(question)
        })
        
        courseData.quizzes.push(quiz)
    })
    
    return courseData
}

$(document).ready(function(){
    // Setup the courses table
    $("#courseTable").DataTable({
        columnDefs: [
            {
                orderable: false,
                targets: "no-sort"
            }
        ]
    });
    load()
    
    // Setup modal form submission
    $("#modalSubmit").click(function(){
        var courseData = serializeCourse()
        
        if (course.id){
            // Update existing course
            courseData.id = course.id
            api.courseId.put(course.id, courseData)
            .success(function(data){
                successFlash("Course updated successfully!")
                location.reload()
            })
            .error(function(data){
                modalError(data.responseJSON.message)
            })
        } else {
            // Create new course
            api.courses.post(courseData)
            .success(function(data){
                successFlash("Course created successfully!")
                location.reload()
            })
            .error(function(data){
                modalError(data.responseJSON.message)
            })
        }
    })
})

function load(){
    api.courses.get()
    .success(function(cs){
        courses = cs
        $("#courseTable").DataTable().clear()
        courseRows = []
        $.each(courses, function(i, course){
            var moduleCount = course.modules ? course.modules.length : 0
            var quizCount = course.quizzes ? course.quizzes.length : 0
            
            courseRows.push([
                escapeHtml(course.name),
                escapeHtml(course.description),
                moduleCount,
                quizCount,
                moment(course.created_date).format('MMMM Do YYYY, h:mm:ss a'),
                moment(course.modified_date).format('MMMM Do YYYY, h:mm:ss a'),
                "<div class='pull-right'><span data-toggle='tooltip' data-placement='left' title='Edit Course'><button class='btn btn-primary' data-toggle='modal' data-backdrop='static' data-target='#modal' onclick='save(" + i + ")'>\
                    <i class='fa fa-pencil'></i>\
                    </button></span>\
                    <span data-toggle='tooltip' data-placement='left' title='Delete Course'><button class='btn btn-danger' onclick='deleteCourse(" + i + ")'>\
                    <i class='fa fa-trash-o'></i>\
                    </button></span></div>"
            ])
        })
        $("#courseTable").DataTable().rows.add(courseRows).draw()
        $('[data-toggle="tooltip"]').tooltip()
    })
    .error(function(){
        errorFlash("Error fetching courses")
    })
}

function edit(course) {
    $("#modalLabel").text("Edit Course")
    $("#name").val(course.name)
    $("#description").val(course.description)
    
    // Clear containers
    $("#modules-container").empty()
    $("#quizzes-container").empty()
    
    // Load modules
    if (course.modules) {
        $.each(course.modules, function(i, module) {
            addModule()
            var modulePanel = $("#modules-container .module-panel").last()
            modulePanel.find('.module-name').val(module.name)
            modulePanel.find('.module-description').val(module.description)
            modulePanel.find('.module-content').val(module.content)
        })
    }
    
    // Load quizzes
    if (course.quizzes) {
        $.each(course.quizzes, function(i, quiz) {
            addQuiz()
            var quizPanel = $("#quizzes-container .quiz-panel").last()
            quizPanel.find('.quiz-name').val(quiz.name)
            quizPanel.find('.quiz-description').val(quiz.description)
            quizPanel.find('.quiz-passing-score').val(quiz.passing_score)
            
            // Load questions
            if (quiz.questions) {
                $.each(quiz.questions, function(j, question) {
                    addQuestion(quizPanel)
                    var questionPanel = quizPanel.find('.question-panel').last()
                    questionPanel.find('.question-text').val(question.question)
                    
                    // Load options
                    if (question.options) {
                        $.each(question.options, function(k, option) {
                            addOption(questionPanel)
                            var optionGroup = questionPanel.find('.option-group').last()
                            optionGroup.find('.option-text').val(option.option)
                            if (option.is_correct) {
                                optionGroup.find('.correct-option').prop('checked', true)
                            }
                        })
                    }
                })
            }
        })
    }
}

function save(idx) {
    course = courses[idx]
    edit(course)
    $("#modal").modal('show')
}