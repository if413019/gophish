$(document).ready(function () {
    $('[data-toggle="tooltip"]').tooltip();
    $("#apiResetForm").submit(function (e) {
        api.reset()
            .success(function (response) {
                user.api_key = response.data
                successFlash(response.message)
                $("#api_key").val(user.api_key)
            })
            .error(function (data) {
                errorFlash(data.message)
            })
        return false
    })
    $("#settingsForm").submit(function (e) {
        $.post("/settings", $(this).serialize())
            .done(function (data) {
                successFlash(data.message)
            })
            .fail(function (data) {
                errorFlash(data.responseJSON.message)
            })
        return false
    })
    //$("#imapForm").submit(function (e) {
    $("#savesettings").click(function() {
        var imapSettings = {}
        imapSettings.host = $("#imaphost").val()
        imapSettings.port = $("#imapport").val()
        imapSettings.username = $("#imapusername").val()
        imapSettings.password = $("#imappassword").val()
        imapSettings.enabled = $('#use_imap').prop('checked')
        imapSettings.tls = $('#use_tls').prop('checked')

        //Advanced settings
        imapSettings.folder = $("#folder").val()
        imapSettings.imap_freq = $("#imapfreq").val()
        imapSettings.restrict_domain = $("#restrictdomain").val()
        imapSettings.ignore_cert_errors = $('#ignorecerterrors').prop('checked')
        imapSettings.delete_reported_campaign_email = $('#deletecampaign').prop('checked')
        
        //To avoid unmarshalling error in controllers/api/imap.go. It would fail gracefully, but with a generic error.
        if (imapSettings.host == ""){
            errorFlash("No IMAP Host specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (imapSettings.port == ""){
            errorFlash("No IMAP Port specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (isNaN(imapSettings.port) || imapSettings.port <1 || imapSettings.port > 65535  ){ 
            errorFlash("Invalid IMAP Port")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (imapSettings.imap_freq == ""){
            imapSettings.imap_freq = "60"
        }

        api.IMAP.post(imapSettings).done(function (data) {
                if (data.success == true) {
                    successFlashFade("Successfully updated IMAP settings.", 2)
                } else {
                    errorFlash("Unable to update IMAP settings.")
                }
            })
            .success(function (data){
                loadIMAPSettings()
            })
            .fail(function (data) {
                errorFlash(data.responseJSON.message)
            })
            .always(function (data){
                document.body.scrollTop = 0;
                document.documentElement.scrollTop = 0;
            })
        
        return false
    })

    $("#validateimap").click(function() {

        // Query validate imap server endpoint
        var server = {}
        server.host = $("#imaphost").val()
        server.port = $("#imapport").val()
        server.username = $("#imapusername").val()
        server.password = $("#imappassword").val()
        server.tls = $('#use_tls').prop('checked')
        server.ignore_cert_errors = $('#ignorecerterrors').prop('checked')

        //To avoid unmarshalling error in controllers/api/imap.go. It would fail gracefully, but with a generic error. 
        if (server.host == ""){
            errorFlash("No IMAP Host specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (server.port == ""){
            errorFlash("No IMAP Port specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (isNaN(server.port) || server.port <1 || server.port > 65535  ){
            errorFlash("Invalid IMAP Port")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }

        var oldHTML = $("#validateimap").html();
        // Disable inputs and change button text
        $("#imaphost").attr("disabled", true);
        $("#imapport").attr("disabled", true);
        $("#imapusername").attr("disabled", true);
        $("#imappassword").attr("disabled", true);
        $("#use_imap").attr("disabled", true);
        $("#use_tls").attr("disabled", true);
        $('#ignorecerterrors').attr("disabled", true);
        $("#folder").attr("disabled", true);
        $("#restrictdomain").attr("disabled", true);
        $('#deletecampaign').attr("disabled", true);
        $('#lastlogin').attr("disabled", true);
        $('#imapfreq').attr("disabled", true);
        $("#validateimap").attr("disabled", true);  
        $("#validateimap").html("<i class='fa fa-circle-o-notch fa-spin'></i> Testing...");
        
        api.IMAP.validate(server).done(function(data) {
            if (data.success == true) {
                Swal.fire({
                    title: "Success",
                    html: "Logged into <b>" + escapeHtml($("#imaphost").val()) + "</b>",
                    type: "success",
                })
            } else {
                Swal.fire({
                    title: "Failed!",
                    html: "Unable to login to <b>" + escapeHtml($("#imaphost").val()) + "</b>.",
                    type: "error",
                    showCancelButton: true,
                    cancelButtonText: "Close",
                    confirmButtonText: "More Info",
                    confirmButtonColor: "#428bca",
                    allowOutsideClick: false,
                }).then(function(result) {
                    if (result.value) {
                        Swal.fire({
                            title: "Error:",
                            text: data.message,
                        })
                    }
                  })
            }
            
          })
          .fail(function() {
            Swal.fire({
                title: "Failed!",
                text: "An unecpected error occured.",
                type: "error",
            })
          })
          .always(function() {
            //Re-enable inputs and change button text
            $("#imaphost").attr("disabled", false);
            $("#imapport").attr("disabled", false);
            $("#imapusername").attr("disabled", false);
            $("#imappassword").attr("disabled", false);
            $("#use_imap").attr("disabled", false);
            $("#use_tls").attr("disabled", false);
            $('#ignorecerterrors').attr("disabled", false);
            $("#folder").attr("disabled", false);
            $("#restrictdomain").attr("disabled", false);
            $('#deletecampaign').attr("disabled", false);
            $('#lastlogin').attr("disabled", false);
            $('#imapfreq').attr("disabled", false);
            $("#validateimap").attr("disabled", false);
            $("#validateimap").html(oldHTML);

          });

      }); //end testclick

    $("#reporttab").click(function() {
        loadIMAPSettings()
    })

    $("#advanced").click(function() {
        $("#advancedarea").toggle();
    })

    function loadIMAPSettings(){
        api.IMAP.get()
        .success(function (imap) {
            if (imap.length == 0){
                $('#lastlogindiv').hide()
            } else {
                imap = imap[0]
                if (imap.enabled == false){
                    $('#lastlogindiv').hide()
                } else {
                    $('#lastlogindiv').show()
                }
                $("#imapusername").val(imap.username)
                $("#imaphost").val(imap.host)
                $("#imapport").val(imap.port)
                $("#imappassword").val(imap.password)
                $('#use_tls').prop('checked', imap.tls)
                $('#ignorecerterrors').prop('checked', imap.ignore_cert_errors)
                $('#use_imap').prop('checked', imap.enabled)
                $("#folder").val(imap.folder)
                $("#restrictdomain").val(imap.restrict_domain)
                $('#deletecampaign').prop('checked', imap.delete_reported_campaign_email)
                $('#lastloginraw').val(imap.last_login)
                $('#lastlogin').val(moment.utc(imap.last_login).fromNow())
                $('#imapfreq').val(imap.imap_freq)
            }  

        })
        .error(function () {
            errorFlash("Error fetching IMAP settings")
        })
    }

    var use_map = localStorage.getItem('gophish.use_map')
    $("#use_map").prop('checked', JSON.parse(use_map))
    $("#use_map").on('change', function () {
        localStorage.setItem('gophish.use_map', JSON.stringify(this.checked))
    })

    loadIMAPSettings()

    // E-Learning Settings Tab
    var elearningSettingsLoaded = false;

    $("#elearningtab").click(function() {
        if (!elearningSettingsLoaded) {
            loadELearningSettings()
            elearningSettingsLoaded = true
        }
    })

    // Load E-Learning Settings
    function loadELearningSettings() {
        api.elearningSettings.get()
            .success(function(data) {
                var settings = data.settings
                var smtpProfiles = data.smtp_profiles

                // Populate SMTP dropdown
                var $smtpSelect = $("#elearning_smtp")
                $smtpSelect.empty()
                $smtpSelect.append('<option value="">-- Select Sending Profile --</option>')

                if (smtpProfiles && smtpProfiles.length > 0) {
                    smtpProfiles.forEach(function(profile) {
                        var selected = settings.smtp_id === profile.id ? ' selected' : ''
                        $smtpSelect.append('<option value="' + profile.id + '"' + selected + '>' + escapeHtml(profile.name) + '</option>')
                    })
                }

                // Populate form fields
                $("#elearning_base_url").val(settings.base_url || 'https://localhost:3333')
                $("#elearning_company_name").val(settings.company_name || 'Your Organization')
                $("#elearning_email_subject").val(settings.email_subject || 'Security Awareness Training Required')

                // Load email template
                if (settings.email_html) {
                    $("#elearning_email_html").val(settings.email_html)
                } else {
                    // Load default template if none exists
                    loadDefaultTemplate()
                }
            })
            .error(function(data) {
                errorFlash("Error loading E-Learning settings")
            })
    }

    // Load default template
    function loadDefaultTemplate() {
        api.elearningSettings.getDefaultTemplate()
            .success(function(data) {
                $("#elearning_email_html").val(data.template)
            })
            .error(function() {
                errorFlash("Error loading default template")
            })
    }

    // Reset to default template button
    $("#reset-template").click(function() {
        Swal.fire({
            title: "Reset Template?",
            text: "This will replace your current template with the default template.",
            type: "warning",
            showCancelButton: true,
            confirmButtonText: "Reset",
            confirmButtonColor: "#dc3545"
        }).then(function(result) {
            if (result.value) {
                loadDefaultTemplate()
                successFlashFade("Template reset to default", 2)
            }
        })
    })

    // Template editor/preview toggle buttons
    $("#btn-template-editor").click(function() {
        $(this).addClass('active')
        $("#btn-template-preview").removeClass('active')
        $("#template-editor-pane").show()
        $("#template-preview-pane").hide()
    })

    $("#btn-template-preview").click(function() {
        $(this).addClass('active')
        $("#btn-template-editor").removeClass('active')
        $("#template-editor-pane").hide()
        $("#template-preview-pane").show()

        // Generate preview
        var template = $("#elearning_email_html").val()
        var iframe = document.getElementById('template-preview-frame')
        var iframeDoc = iframe.contentDocument || iframe.contentWindow.document

        if (!template) {
            iframeDoc.open()
            iframeDoc.write('<html><body style="font-family: sans-serif; padding: 20px; color: #666;"><p>No template to preview. Please enter an HTML template.</p></body></html>')
            iframeDoc.close()
            return
        }

        api.elearningSettings.previewTemplate(template)
            .success(function(data) {
                iframeDoc.open()
                iframeDoc.write(data.preview)
                iframeDoc.close()
            })
            .error(function(data) {
                var message = data.responseJSON ? data.responseJSON.message : "Error generating preview"
                iframeDoc.open()
                iframeDoc.write('<html><body style="font-family: sans-serif; padding: 20px;"><div style="color: #721c24; background-color: #f8d7da; border: 1px solid #f5c6cb; padding: 15px; border-radius: 4px;">' + escapeHtml(message) + '</div></body></html>')
                iframeDoc.close()
            })
    })

    // Save E-Learning Settings
    $("#save-elearning-settings").click(function() {
        var settings = {
            smtp_id: parseInt($("#elearning_smtp").val()) || 0,
            base_url: $("#elearning_base_url").val().trim(),
            company_name: $("#elearning_company_name").val().trim(),
            email_subject: $("#elearning_email_subject").val().trim(),
            email_html: $("#elearning_email_html").val()
        }

        // Validation
        if (!settings.smtp_id) {
            errorFlash("Please select a Sending Profile")
            return
        }
        if (!settings.base_url) {
            errorFlash("Please enter a Base URL")
            return
        }
        if (!settings.company_name) {
            errorFlash("Please enter a Company Name")
            return
        }

        var $btn = $(this)
        var oldHTML = $btn.html()
        $btn.attr("disabled", true)
        $btn.html('<i class="fa fa-circle-o-notch fa-spin"></i> Saving...')

        api.elearningSettings.post(settings)
            .success(function(data) {
                if (data.success) {
                    successFlashFade("E-Learning settings saved successfully", 3)
                } else {
                    errorFlash(data.message || "Error saving settings")
                }
            })
            .error(function(data) {
                var message = data.responseJSON ? data.responseJSON.message : "Error saving settings"
                errorFlash(message)
            })
            .always(function() {
                $btn.attr("disabled", false)
                $btn.html(oldHTML)
            })
    })

    // Send Test Email
    $("#send-test-email").click(function() {
        var email = $("#test_email_address").val().trim()
        if (!email) {
            errorFlash("Please enter a test email address")
            return
        }

        // Basic email validation
        var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(email)) {
            errorFlash("Please enter a valid email address")
            return
        }

        var $btn = $(this)
        var oldHTML = $btn.html()
        $btn.attr("disabled", true)
        $btn.html('<i class="fa fa-circle-o-notch fa-spin"></i> Sending...')

        api.elearningSettings.testEmail(email)
            .success(function(data) {
                if (data.success) {
                    Swal.fire({
                        title: "Success!",
                        text: data.message,
                        type: "success"
                    })
                } else {
                    Swal.fire({
                        title: "Failed",
                        text: data.message || "Failed to send test email",
                        type: "error"
                    })
                }
            })
            .error(function(data) {
                var message = data.responseJSON ? data.responseJSON.message : "Error sending test email"
                Swal.fire({
                    title: "Error",
                    text: message,
                    type: "error"
                })
            })
            .always(function() {
                $btn.attr("disabled", false)
                $btn.html(oldHTML)
            })
    })
})
