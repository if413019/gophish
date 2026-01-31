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
    $("#content-container").empty()
}

function newCourse(){
    course = {}
    $("#modalLabel").text("New Course")
    $("#name").val("")
    $("#description").val("")
    $("#content-container").empty()
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
    var moduleCount = $("#content-container .module-panel").length
    var moduleElement = $(moduleHtml)
    var uniqueId = Date.now() + "_" + moduleCount

    // Fix radio button names to be unique per module
    moduleElement.find('input[name="video-source"]').attr('name', 'video-source-' + uniqueId)
    moduleElement.find('input[name="presentation-source"]').attr('name', 'presentation-source-' + uniqueId)

    // Fix checkbox id and label for attribute to be unique
    var checkboxId = 'module-must-complete-' + uniqueId
    moduleElement.find('.module-must-complete').attr('id', checkboxId)
    moduleElement.find('label[for="module-must-complete-placeholder"]').attr('for', checkboxId)

    // Handle remove module
    moduleElement.find('.remove-module').click(function() {
        $(this).closest('.module-panel').remove()
        updateOrderIndices()
    })

    // Handle module type change
    moduleElement.find('.module-type').change(function() {
        toggleModuleFields($(this))
    })

    // Handle move up/down
    moduleElement.find('.move-up').click(function() {
        moveItemUp($(this).closest('.content-item'))
    })
    moduleElement.find('.move-down').click(function() {
        moveItemDown($(this).closest('.content-item'))
    })

    $("#content-container").append(moduleElement)

    // Trigger initial field visibility based on default selection
    toggleModuleFields(moduleElement.find('.module-type'))
    updateOrderIndices()
}

function toggleModuleFields(moduleTypeSelect) {
    var modulePanel = moduleTypeSelect.closest('.module-panel')
    var moduleType = moduleTypeSelect.val()

    console.log('Toggling module fields for type:', moduleType)

    // Hide all type-specific fields first
    modulePanel.find('.module-video-fields').hide()
    modulePanel.find('.module-presentation-fields').hide()

    // Show appropriate fields based on type
    switch(moduleType) {
        case 'video':
            console.log('Showing video fields')
            modulePanel.find('.module-video-fields').show()
            setupVideoSourceToggle(modulePanel)
            break
        case 'presentation':
            console.log('Showing presentation fields')
            modulePanel.find('.module-presentation-fields').show()
            setupPresentationSourceToggle(modulePanel)
            break
        case 'html':
        default:
            console.log('HTML type selected - no extra fields needed')
            // HTML is default, no extra fields needed
            break
    }
}

function setupVideoSourceToggle(modulePanel) {
    var videoRadios = modulePanel.find('input[type="radio"]').filter(function() {
        return $(this).attr('name') && $(this).attr('name').startsWith('video-source-')
    })
    var uploadSection = modulePanel.find('.video-upload-section')
    var urlSection = modulePanel.find('.video-url-section')
    
    videoRadios.change(function() {
        if ($(this).val() === 'upload') {
            uploadSection.show()
            urlSection.hide()
        } else {
            uploadSection.hide()
            urlSection.show()
        }
    })
    
    // Trigger initial state
    videoRadios.filter(':checked').trigger('change')
    
    // Setup file upload handler
    var fileInput = modulePanel.find('.module-video-file')
    fileInput.change(function() {
        handleFileUpload($(this), 'video', modulePanel)
    })
    
    // Setup remove uploaded file handler
    modulePanel.find('.remove-uploaded-file').click(function() {
        removeUploadedFile($(this), 'video', modulePanel)
    })
}

function setupPresentationSourceToggle(modulePanel) {
    var presentationRadios = modulePanel.find('input[type="radio"]').filter(function() {
        return $(this).attr('name') && $(this).attr('name').startsWith('presentation-source-')
    })
    var uploadSection = modulePanel.find('.presentation-upload-section')
    var urlSection = modulePanel.find('.presentation-url-section')
    
    presentationRadios.change(function() {
        if ($(this).val() === 'upload') {
            uploadSection.show()
            urlSection.hide()
        } else {
            uploadSection.hide()
            urlSection.show()
        }
    })
    
    // Trigger initial state
    presentationRadios.filter(':checked').trigger('change')
    
    // Setup file upload handler
    var fileInput = modulePanel.find('.module-presentation-file')
    fileInput.change(function() {
        handleFileUpload($(this), 'presentation', modulePanel)
    })
    
    // Setup remove uploaded file handler
    modulePanel.find('.remove-uploaded-file').click(function() {
        removeUploadedFile($(this), 'presentation', modulePanel)
    })
}

function handleFileUpload(fileInput, fileType, modulePanel) {
    var file = fileInput[0].files[0]
    if (!file) return
    
    var progressDiv = modulePanel.find('.upload-progress')
    var progressBar = progressDiv.find('.progress-bar')
    var statusText = progressDiv.find('.upload-status')
    var uploadedInfo = modulePanel.find('.uploaded-file-info')
    
    // Show progress
    progressDiv.show()
    progressBar.css('width', '0%')
    statusText.text('Uploading...')
    uploadedInfo.hide()
    
    // Create form data
    var formData = new FormData()
    formData.append('file', file)
    
    // Upload file
    $.ajax({
        url: '/api/courses/upload/' + fileType,
        type: 'POST',
        data: formData,
        processData: false,
        contentType: false,
        xhr: function() {
            var xhr = new window.XMLHttpRequest()
            xhr.upload.addEventListener('progress', function(e) {
                if (e.lengthComputable) {
                    var percent = Math.round((e.loaded / e.total) * 100)
                    progressBar.css('width', percent + '%')
                    statusText.text('Uploading... ' + percent + '%')
                }
            }, false)
            return xhr
        },
        success: function(response) {
            if (response.success) {
                progressDiv.hide()
                uploadedInfo.show()
                uploadedInfo.find('.filename').text(response.original_filename)
                
                // Store file info for form submission
                modulePanel.data('uploaded-file', {
                    filename: response.filename,
                    original_filename: response.original_filename,
                    file_path: response.file_path,
                    file_size: response.file_size,
                    mime_type: response.mime_type
                })
                
                successFlash('File uploaded successfully: ' + response.original_filename)
            } else {
                progressDiv.hide()
                modalError(response.message || 'Upload failed')
            }
        },
        error: function(xhr) {
            progressDiv.hide()
            var errorMsg = 'Upload failed'
            if (xhr.responseJSON && xhr.responseJSON.message) {
                errorMsg = xhr.responseJSON.message
            }
            modalError(errorMsg)
        }
    })
}

function removeUploadedFile(button, fileType, modulePanel) {
    var uploadedFile = modulePanel.data('uploaded-file')
    if (!uploadedFile) return
    
    // Delete file from server
    $.ajax({
        url: '/api/courses/files/' + fileType + '/' + uploadedFile.filename,
        type: 'DELETE',
        success: function(response) {
            if (response.success) {
                // Clear UI
                modulePanel.find('.uploaded-file-info').hide()
                modulePanel.find('.module-' + fileType + '-file').val('')
                modulePanel.removeData('uploaded-file')
                successFlash('File removed successfully')
            } else {
                modalError(response.message || 'Failed to remove file')
            }
        },
        error: function(xhr) {
            var errorMsg = 'Failed to remove file'
            if (xhr.responseJSON && xhr.responseJSON.message) {
                errorMsg = xhr.responseJSON.message
            }
            modalError(errorMsg)
        }
    })
}

// Move item up in the content container
function moveItemUp(item) {
    var prev = item.prev('.content-item')
    if (prev.length) {
        item.insertBefore(prev)
        updateOrderIndices()
    }
}

// Move item down in the content container
function moveItemDown(item) {
    var next = item.next('.content-item')
    if (next.length) {
        item.insertAfter(next)
        updateOrderIndices()
    }
}

// Update visual order indices for all items
function updateOrderIndices() {
    $("#content-container .content-item").each(function(index) {
        $(this).data('order-index', index)
    })
}

function addQuiz() {
    var quizHtml = $("#quiz-template").html()
    var quizElement = $(quizHtml)

    quizElement.find('.remove-quiz').click(function() {
        $(this).closest('.quiz-panel').remove()
        updateOrderIndices()
    })

    quizElement.find('.add-question').click(function() {
        addQuestion($(this).closest('.quiz-panel'))
    })

    // Handle move up/down
    quizElement.find('.move-up').click(function() {
        moveItemUp($(this).closest('.content-item'))
    })
    quizElement.find('.move-down').click(function() {
        moveItemDown($(this).closest('.content-item'))
    })

    $("#content-container").append(quizElement)
    updateOrderIndices()
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

    // Serialize all content items in order from the unified container
    $("#content-container .content-item").each(function(index) {
        var itemType = $(this).data('item-type')

        if (itemType === 'module') {
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

            // Include module ID and timestamps if this is an existing module
            var moduleId = modulePanel.data('module-id')
            if (moduleId) {
                module.id = moduleId
                module.created_date = modulePanel.data('created-date')
                // Don't include modified_date - let the backend set it
            }

            // Add type-specific fields
            if (moduleType === 'video') {
                var videoSource = modulePanel.find('input[type="radio"]:checked').filter(function() {
                    return $(this).attr('name') && $(this).attr('name').startsWith('video-source-')
                }).val()
                if (videoSource === 'upload') {
                    var uploadedFile = modulePanel.data('uploaded-file')
                    if (uploadedFile) {
                        module.video_file_path = uploadedFile.file_path
                        module.original_filename = uploadedFile.original_filename
                        module.file_size = uploadedFile.file_size
                        module.mime_type = uploadedFile.mime_type
                    }
                } else {
                    module.video_url = modulePanel.find('.module-video-url').val()
                }
            } else if (moduleType === 'presentation') {
                var presentationSource = modulePanel.find('input[type="radio"]:checked').filter(function() {
                    return $(this).attr('name') && $(this).attr('name').startsWith('presentation-source-')
                }).val()
                var uploadedFile = modulePanel.data('uploaded-file')
                if (presentationSource === 'upload') {
                    if (uploadedFile) {
                        module.presentation_file_path = uploadedFile.file_path
                        module.original_filename = uploadedFile.original_filename
                        module.file_size = uploadedFile.file_size
                        module.mime_type = uploadedFile.mime_type
                    }
                } else {
                    module.presentation_url = modulePanel.find('.module-presentation-url').val()
                }
            }

            courseData.modules.push(module)

        } else if (itemType === 'quiz') {
            var quizPanel = $(this)
            var quiz = {
                name: quizPanel.find('.quiz-name').val(),
                description: quizPanel.find('.quiz-description').val(),
                passing_score: parseInt(quizPanel.find('.quiz-passing-score').val()) || 70,
                time_limit: parseInt(quizPanel.find('.quiz-time-limit').val()) || 0,
                max_attempts: parseInt(quizPanel.find('.quiz-max-attempts').val()) || 3,
                order_index: index,
                questions: []
            }

            // Include quiz ID if this is an existing quiz
            var quizId = quizPanel.data('quiz-id')
            if (quizId) {
                quiz.id = quizId
            }

            quizPanel.find('.question-panel').each(function(questionIndex) {
                var questionPanel = $(this)
                var question = {
                    question: questionPanel.find('.question-text').val(),
                    order_index: questionIndex,
                    options: []
                }

                // Include question ID if this is an existing question
                var questionId = questionPanel.data('question-id')
                if (questionId) {
                    question.id = questionId
                }

                questionPanel.find('.option-group').each(function(optionIndex) {
                    var optionGroup = $(this)
                    var isCorrect = optionGroup.find('.correct-option').is(':checked')
                    var option = {
                        option: optionGroup.find('.option-text').val(),
                        is_correct: isCorrect,
                        order_index: optionIndex
                    }

                    // Include option ID if this is an existing option
                    var optionId = optionGroup.data('option-id')
                    if (optionId) {
                        option.id = optionId
                    }

                    question.options.push(option)
                })

                quiz.questions.push(question)
            })

            courseData.quizzes.push(quiz)
        }
    })

    return courseData
}

$(document).ready(function(){
    // Setup the courses table with modern configuration
    $("#courseTable").DataTable({
        columnDefs: [
            {
                orderable: false,
                targets: "no-sort"
            }
        ],
        language: {
            search: "",
            searchPlaceholder: "Search courses...",
            lengthMenu: "Show _MENU_ courses per page",
            info: "Showing _START_ to _END_ of _TOTAL_ courses",
            infoEmpty: "No courses available",
            infoFiltered: "(filtered from _MAX_ total courses)",
            paginate: {
                first: "First",
                last: "Last",
                next: "Next",
                previous: "Previous"
            },
            emptyTable: "No courses created yet"
        },
        pageLength: 10,
        lengthMenu: [[5, 10, 25, 50, -1], [5, 10, 25, 50, "All"]],
        dom: '<"dataTables_top_controls"<"dataTables_length"l><"dataTables_filter"f>>t<"dataTables_bottom_controls"<"dataTables_info"i><"dataTables_paginate"p>>',
        drawCallback: function() {
            // Add modern styling to pagination buttons after each draw
            $('.dataTables_wrapper .paginate_button').each(function() {
                var $btn = $(this);
                if ($btn.hasClass('previous')) {
                    $btn.html('<i class="fa fa-chevron-left"></i>');
                } else if ($btn.hasClass('next')) {
                    $btn.html('<i class="fa fa-chevron-right"></i>');
                }
            });
        },
        initComplete: function() {
            // Add search icon and improve styling after initialization
            $('.dataTables_filter label').prepend('<i class="fa fa-search" style="margin-right: 0.5rem; color: #9ca3af;"></i>');
            $('.dataTables_filter input').attr('placeholder', 'Search courses...');
            
            // Add length menu icon
            $('.dataTables_length label').prepend('<i class="fa fa-list" style="margin-right: 0.5rem; color: #9ca3af;"></i>');
        }
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
                "<div class='course-name'>" + escapeHtml(course.name) + "</div>" +
                "<div class='course-description' title='" + escapeHtml(course.description) + "'>" + escapeHtml(course.description) + "</div>",
                "<div class='course-meta'>" +
                    "<span class='meta-badge modules'><i class='fa fa-list-ol'></i> " + moduleCount + " Modules</span>" +
                    "<span class='meta-badge quizzes'><i class='fa fa-question-circle'></i> " + quizCount + " Quizzes</span>" +
                "</div>",
                "<div class='course-date'>" + moment(course.created_date).format('MMM DD, YYYY') + "</div>",
                "<div class='course-date'>" + moment(course.modified_date).format('MMM DD, YYYY') + "</div>",
                "<div class='action-buttons'>" +
                    "<button class='btn-sm-modern btn-secondary-modern' data-toggle='tooltip' data-placement='left' title='Preview Course' onclick='window.location.href=\"/courses/" + course.id + "/preview\"'>" +
                        "<i class='fa fa-eye'></i>" +
                    "</button>" +
                    "<button class='btn-sm-modern btn-primary-modern' data-toggle='tooltip' data-placement='left' title='Edit Course' onclick='save(" + i + ")'>" +
                        "<i class='fa fa-pencil'></i>" +
                    "</button>" +
                    "<button class='btn-sm-modern btn-danger-modern' data-toggle='tooltip' data-placement='left' title='Delete Course' onclick='deleteCourse(" + i + ")'>" +
                        "<i class='fa fa-trash-o'></i>" +
                    "</button>" +
                "</div>"
            ])
        })
        if (courseRows.length === 0) {
            // Show empty state
            $("#courseTable").parent().html(
                "<div class='empty-state'>" +
                    "<i class='fa fa-graduation-cap'></i>" +
                    "<h3>No courses yet</h3>" +
                    "<p>Get started by creating your first e-learning course</p>" +
                    "<button class='btn-modern btn-primary-modern' onclick='newCourse()'>" +
                        "<i class='fa fa-plus'></i> Create Your First Course" +
                    "</button>" +
                "</div>"
            )
        } else {
            $("#courseTable").DataTable().rows.add(courseRows).draw()
            $('[data-toggle="tooltip"]').tooltip()
        }
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
    $("#content-container").empty()

    // Combine modules and quizzes into a single sorted list
    var contentItems = []

    if (course.modules) {
        $.each(course.modules, function(i, module) {
            contentItems.push({
                type: 'module',
                data: module,
                order_index: module.order_index || i
            })
        })
    }

    if (course.quizzes) {
        $.each(course.quizzes, function(i, quiz) {
            contentItems.push({
                type: 'quiz',
                data: quiz,
                order_index: quiz.order_index || i
            })
        })
    }

    // Sort by order_index
    contentItems.sort(function(a, b) {
        return a.order_index - b.order_index
    })

    // Add items in sorted order
    $.each(contentItems, function(i, item) {
        if (item.type === 'module') {
            loadModuleItem(item.data)
        } else if (item.type === 'quiz') {
            loadQuizItem(item.data)
        }
    })

    updateOrderIndices()
}

function loadModuleItem(module) {
    addModule()
    var modulePanel = $("#content-container .module-panel").last()

    // Store the module ID and timestamps for existing modules
    if (module.id) {
        modulePanel.data('module-id', module.id)
        modulePanel.data('created-date', module.created_date)
        modulePanel.data('modified-date', module.modified_date)
    }

    modulePanel.find('.module-name').val(module.name)
    modulePanel.find('.module-description').val(module.description)
    modulePanel.find('.module-content').val(module.content)
    modulePanel.find('.module-type').val(module.module_type).trigger('change')
    modulePanel.find('.module-must-complete').prop('checked', module.must_complete)
    modulePanel.find('.module-min-time').val(module.min_time_spent)

    // Load type-specific data
    if (module.module_type === 'video') {
        if (module.video_file_path) {
            // Uploaded file
            modulePanel.find('input[type="radio"]').filter(function() {
                return $(this).attr('name') && $(this).attr('name').startsWith('video-source-') && $(this).val() === 'upload'
            }).prop('checked', true).trigger('change')
            if (module.original_filename) {
                modulePanel.find('.uploaded-file-info').show()
                modulePanel.find('.filename').text(module.original_filename)
                modulePanel.data('uploaded-file', {
                    filename: module.video_file_path.split('/').pop(),
                    original_filename: module.original_filename,
                    file_path: module.video_file_path,
                    file_size: module.file_size,
                    mime_type: module.mime_type
                })
            }
        } else if (module.video_url) {
            // External URL
            modulePanel.find('input[type="radio"]').filter(function() {
                return $(this).attr('name') && $(this).attr('name').startsWith('video-source-') && $(this).val() === 'url'
            }).prop('checked', true).trigger('change')
            modulePanel.find('.module-video-url').val(module.video_url)
        }
    } else if (module.module_type === 'presentation') {
        if (module.presentation_file_path) {
            // Uploaded file
            modulePanel.find('input[type="radio"]').filter(function() {
                return $(this).attr('name') && $(this).attr('name').startsWith('presentation-source-') && $(this).val() === 'upload'
            }).prop('checked', true).trigger('change')
            if (module.original_filename) {
                modulePanel.find('.uploaded-file-info').show()
                modulePanel.find('.filename').text(module.original_filename)
                modulePanel.data('uploaded-file', {
                    filename: module.presentation_file_path.split('/').pop(),
                    original_filename: module.original_filename,
                    file_path: module.presentation_file_path,
                    file_size: module.file_size,
                    mime_type: module.mime_type
                })
            }
        } else if (module.presentation_url) {
            // External URL
            modulePanel.find('input[type="radio"]').filter(function() {
                return $(this).attr('name') && $(this).attr('name').startsWith('presentation-source-') && $(this).val() === 'url'
            }).prop('checked', true).trigger('change')
            modulePanel.find('.module-presentation-url').val(module.presentation_url)
        }
    }
}

function loadQuizItem(quiz) {
    addQuiz()
    var quizPanel = $("#content-container .quiz-panel").last()

    // Store the quiz ID for existing quizzes
    if (quiz.id) {
        quizPanel.data('quiz-id', quiz.id)
    }

    quizPanel.find('.quiz-name').val(quiz.name)
    quizPanel.find('.quiz-description').val(quiz.description)
    quizPanel.find('.quiz-passing-score').val(quiz.passing_score)
    quizPanel.find('.quiz-time-limit').val(quiz.time_limit || 0)
    quizPanel.find('.quiz-max-attempts').val(quiz.max_attempts || 3)

    // Load questions
    if (quiz.questions) {
        $.each(quiz.questions, function(j, question) {
            addQuestion(quizPanel)
            var questionPanel = quizPanel.find('.question-panel').last()

            // Store the question ID for existing questions
            if (question.id) {
                questionPanel.data('question-id', question.id)
            }

            questionPanel.find('.question-text').val(question.question)

            // Load options
            if (question.options) {
                $.each(question.options, function(k, option) {
                    addOption(questionPanel)
                    var optionGroup = questionPanel.find('.option-group').last()

                    // Store the option ID for existing options
                    if (option.id) {
                        optionGroup.data('option-id', option.id)
                    }

                    optionGroup.find('.option-text').val(option.option)
                    if (option.is_correct) {
                        optionGroup.find('.correct-option').prop('checked', true)
                    }
                })
            }
        })
    }
}

function save(idx) {
    // Fetch full course data including questions and options
    api.courseId.get(courses[idx].id)
    .success(function(fullCourse) {
        course = fullCourse
        edit(course)
        $("#modal").modal('show')
    })
    .error(function(data) {
        errorFlash("Error loading course details")
    })
}