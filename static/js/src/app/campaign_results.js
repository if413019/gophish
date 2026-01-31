var map = null
var doPoll = true;

// statuses is a helper map to point result statuses to ui classes
var statuses = {
    "Email Sent": {
        color: "#1abc9c",
        label: "label-modern label-success-modern",
        icon: "fa-envelope",
        point: "ct-point-sent"
    },
    "Emails Sent": {
        color: "#1abc9c",
        label: "label-modern label-success-modern",
        icon: "fa-envelope",
        point: "ct-point-sent"
    },
    "In progress": {
        label: "label-modern label-primary-modern"
    },
    "Queued": {
        label: "label-modern label-info-modern"
    },
    "Completed": {
        label: "label-modern label-success-modern"
    },
    "Email Opened": {
        color: "#f9bf3b",
        label: "label-modern label-warning-modern",
        icon: "fa-envelope-open",
        point: "ct-point-opened"
    },
    "Clicked Link": {
        color: "#F39C12",
        label: "label-modern label-warning-modern",
        icon: "fa-mouse-pointer",
        point: "ct-point-clicked"
    },
    "Success": {
        color: "#f05b4f",
        label: "label-modern label-danger-modern",
        icon: "fa-exclamation",
        point: "ct-point-clicked"
    },
    //not a status, but is used for the campaign timeline and user timeline
    "Email Reported": {
        color: "#45d6ef",
        label: "label-modern label-info-modern",
        icon: "fa-bullhorn",
        point: "ct-point-reported"
    },
    "Course Enrolled": {
        color: "#9b59b6",
        label: "label-modern label-primary-modern",
        icon: "fa-graduation-cap",
        point: "ct-point-enrolled"
    },
    "Course Completed": {
        color: "#27ae60",
        label: "label-modern label-success-modern",
        icon: "fa-trophy",
        point: "ct-point-completed"
    },
    "Error": {
        color: "#6c7a89",
        label: "label-modern label-secondary-modern",
        icon: "fa-times",
        point: "ct-point-error"
    },
    "Error Sending Email": {
        color: "#6c7a89",
        label: "label-modern label-secondary-modern",
        icon: "fa-times",
        point: "ct-point-error"
    },
    "Submitted Data": {
        color: "#f05b4f",
        label: "label-modern label-danger-modern",
        icon: "fa-exclamation",
        point: "ct-point-clicked"
    },
    "Unknown": {
        color: "#6c7a89",
        label: "label-modern label-secondary-modern",
        icon: "fa-question",
        point: "ct-point-error"
    },
    "Sending": {
        color: "#428bca",
        label: "label-modern label-primary-modern",
        icon: "fa-spinner",
        point: "ct-point-sending"
    },
    "Retrying": {
        color: "#6c7a89",
        label: "label-modern label-secondary-modern",
        icon: "fa-clock-o",
        point: "ct-point-error"
    },
    "Scheduled": {
        color: "#428bca",
        label: "label-modern label-primary-modern",
        icon: "fa-clock-o",
        point: "ct-point-sending"
    },
    "Campaign Created": {
        label: "label-modern label-success-modern",
        icon: "fa-rocket"
    }
}

var statusMapping = {
    "Email Sent": "sent",
    "Email Opened": "opened",
    "Clicked Link": "clicked",
    "Submitted Data": "submitted_data",
    "Email Reported": "reported",
}

// Learning status styles
var learningStatuses = {
    "not-needed": {
        label: "label-modern label-secondary-modern",
        text: "Not Needed",
        icon: "fa-minus"
    },
    "not-started": {
        label: "label-modern label-info-modern",
        text: "Not Started",
        icon: "fa-clock-o"
    },
    "in-progress": {
        label: "label-modern label-warning-modern",
        text: "In Progress",
        icon: "fa-spinner"
    },
    "completed": {
        label: "label-modern label-success-modern",
        text: "Completed",
        icon: "fa-check"
    }
}

// This is an underwhelming attempt at an enum
// until I have time to refactor this appropriately.
var progressListing = [
    "Email Sent",
    "Email Opened",
    "Clicked Link",
    "Submitted Data"
]

var campaign = {}
var bubbles = []

function dismiss() {
    $("#modal\\.flashes").empty()
    $("#modal").modal('hide')
    $("#resultsTable").dataTable().DataTable().clear().draw()
}

// Deletes a campaign after prompting the user
function deleteCampaign() {
    Swal.fire({
        title: "Are you sure?",
        text: "This will delete the campaign. This can't be undone!",
        type: "warning",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Delete Campaign",
        confirmButtonColor: "#ff6b6b",
        reverseButtons: true,
        allowOutsideClick: false,
        showLoaderOnConfirm: true,
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                api.campaignId.delete(campaign.id)
                    .success(function (msg) {
                        resolve()
                    })
                    .error(function (data) {
                        reject(data.responseJSON.message)
                    })
            })
        }
    }).then(function (result) {
        if(result.value){
            Swal.fire(
                'Campaign Deleted!',
                'This campaign has been deleted!',
                'success'
            );
        }
        $('button:contains("OK")').on('click', function () {
            location.href = '/campaigns'
        })
    })
}

// Completes a campaign after prompting the user
function completeCampaign() {
    Swal.fire({
        title: "Are you sure?",
        text: "Gophish will stop processing events for this campaign",
        type: "warning",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Complete Campaign",
        confirmButtonColor: "#667eea",
        reverseButtons: true,
        allowOutsideClick: false,
        showLoaderOnConfirm: true,
        preConfirm: function () {
            return new Promise(function (resolve, reject) {
                api.campaignId.complete(campaign.id)
                    .success(function (msg) {
                        resolve()
                    })
                    .error(function (data) {
                        reject(data.responseJSON.message)
                    })
            })
        }
    }).then(function (result) {
        if (result.value){
            Swal.fire(
                'Campaign Completed!',
                'This campaign has been completed!',
                'success'
            );
            $('#complete_button')[0].disabled = true;
            $('#complete_button').text('Completed!')
            doPoll = false;
        }
    })
}

// Exports campaign results as a CSV file
function exportAsCSV(scope) {
    exportHTML = $("#exportButton").html()
    var csvScope = null
    var filename = campaign.name + ' - ' + capitalize(scope) + '.csv'
    switch (scope) {
        case "results":
            csvScope = campaign.results
            break;
        case "events":
            csvScope = campaign.timeline
            break;
    }
    if (!csvScope) {
        return
    }
    $("#exportButton").html('<i class="fa fa-spinner fa-spin"></i>')
    var csvString = Papa.unparse(csvScope, {
        'escapeFormulae': true
    })
    var csvData = new Blob([csvString], {
        type: 'text/csv;charset=utf-8;'
    });
    if (navigator.msSaveBlob) {
        navigator.msSaveBlob(csvData, filename);
    } else {
        var csvURL = window.URL.createObjectURL(csvData);
        var dlLink = document.createElement('a');
        dlLink.href = csvURL;
        dlLink.setAttribute('download', filename)
        document.body.appendChild(dlLink)
        dlLink.click();
        document.body.removeChild(dlLink)
    }
    $("#exportButton").html(exportHTML)
}

function replay(event_idx) {
    request = campaign.timeline[event_idx]
    details = JSON.parse(request.details)
    url = null
    form = $('<form>').attr({
        method: 'POST',
        target: '_blank',
    })
    /* Create a form object and submit it */
    $.each(Object.keys(details.payload), function (i, param) {
        if (param == "rid") {
            return true;
        }
        if (param == "__original_url") {
            url = details.payload[param];
            return true;
        }
        $('<input>').attr({
            name: param,
        }).val(details.payload[param]).appendTo(form);
    })
    /* Ensure we know where to send the user */
    // Prompt for the URL
    Swal.fire({
        title: 'Where do you want the credentials submitted to?',
        input: 'text',
        showCancelButton: true,
        inputPlaceholder: "http://example.com/login",
        inputValue: url || "",
        inputValidator: function (value) {
            return new Promise(function (resolve, reject) {
                if (value) {
                    resolve();
                } else {
                    reject('Invalid URL.');
                }
            });
        }
    }).then(function (result) {
        if (result.value){
            url = result.value
            submitForm()
        }
    })
    return
    submitForm()

    function submitForm() {
        form.attr({
            action: url
        })
        form.appendTo('body').submit().remove()
    }
}

/**
 * Returns an HTML string that displays the OS and browser that clicked the link
 * or submitted credentials.
 * 
 * @param {object} event_details - The "details" parameter for a campaign
 *  timeline event
 * 
 */
var renderDevice = function (event_details) {
    var ua = UAParser(details.browser['user-agent'])
    var detailsString = '<div class="timeline-device-details">'

    var deviceIcon = 'laptop'
    if (ua.device.type) {
        if (ua.device.type == 'tablet' || ua.device.type == 'mobile') {
            deviceIcon = ua.device.type
        }
    }

    var deviceVendor = ''
    if (ua.device.vendor) {
        deviceVendor = ua.device.vendor.toLowerCase()
        if (deviceVendor == 'microsoft') deviceVendor = 'windows'
    }

    var deviceName = 'Unknown'
    if (ua.os.name) {
        deviceName = ua.os.name
        if (deviceName == "Mac OS") {
            deviceVendor = 'apple'
        } else if (deviceName == "Windows") {
            deviceVendor = 'windows'
        }
        if (ua.device.vendor && ua.device.model) {
            deviceName = ua.device.vendor + ' ' + ua.device.model
        }
    }

    if (ua.os.version) {
        deviceName = deviceName + ' (OS Version: ' + ua.os.version + ')'
    }

    deviceString = '<div class="timeline-device-os"><span class="fa fa-stack">' +
        '<i class="fa fa-' + escapeHtml(deviceIcon) + ' fa-stack-2x"></i>' +
        '<i class="fa fa-vendor-icon fa-' + escapeHtml(deviceVendor) + ' fa-stack-1x"></i>' +
        '</span> ' + escapeHtml(deviceName) + '</div>'

    detailsString += deviceString

    var deviceBrowser = 'Unknown'
    var browserIcon = 'info-circle'
    var browserVersion = ''

    if (ua.browser && ua.browser.name) {
        deviceBrowser = ua.browser.name
        // Handle the "mobile safari" case
        deviceBrowser = deviceBrowser.replace('Mobile ', '')
        if (deviceBrowser) {
            browserIcon = deviceBrowser.toLowerCase()
            if (browserIcon == 'ie') browserIcon = 'internet-explorer'
        }
        browserVersion = '(Version: ' + ua.browser.version + ')'
    }

    var browserString = '<div class="timeline-device-browser"><span class="fa fa-stack">' +
        '<i class="fa fa-' + escapeHtml(browserIcon) + ' fa-stack-1x"></i></span> ' +
        deviceBrowser + ' ' + browserVersion + '</div>'

    detailsString += browserString
    detailsString += '</div>'
    return detailsString
}

function renderTimeline(data) {
    record = {
        "id": data[0],
        "first_name": data[2],
        "last_name": data[3],
        "email": data[4],
        "position": data[5],
        "status": data[6],
        "reported": data[7],
        "send_date": data[8]
    }
    results = '<div class="timeline col-sm-12 well well-lg">' +
        '<h6>Timeline for ' + escapeHtml(record.first_name) + ' ' + escapeHtml(record.last_name) +
        '</h6><span class="subtitle">Email: ' + escapeHtml(record.email) +
        '<br>Result ID: ' + escapeHtml(record.id) + '</span>' +
        '<div class="timeline-graph col-sm-6">'
    $.each(campaign.timeline, function (i, event) {
        if (!event.email || event.email == record.email) {
            // Add the event
            results += '<div class="timeline-entry">' +
                '    <div class="timeline-bar"></div>'
            results +=
                '    <div class="timeline-icon ' + statuses[event.message].label + '">' +
                '    <i class="fa ' + statuses[event.message].icon + '"></i></div>' +
                '    <div class="timeline-message">' + escapeHtml(event.message) +
                '    <span class="timeline-date">' + moment.utc(event.time).local().format('MMMM Do YYYY h:mm:ss a') + '</span>'
            if (event.details) {
                details = JSON.parse(event.details)
                if (event.message == "Clicked Link" || event.message == "Submitted Data") {
                    deviceView = renderDevice(details)
                    if (deviceView) {
                        results += deviceView
                    }
                }
                if (event.message == "Course Enrolled") {
                    results += '<div class="timeline-enrollment-details">'
                    results += '<div class="enrollment-course-info">'
                    results += '<strong>Course:</strong> ' + escapeHtml(details.course_name || 'Unknown Course') + '<br>'
                    results += '<strong>Triggered by:</strong> ' + escapeHtml(details.enrollment_trigger || 'Unknown') + '<br>'
                    results += '<strong>Course ID:</strong> ' + escapeHtml(details.course_id || 'N/A')
                    results += '</div></div>'
                }
                if (event.message == "Course Completed") {
                    results += '<div class="timeline-completion-details">'
                    results += '<div class="completion-course-info">'
                    results += '<strong>Course:</strong> ' + escapeHtml(details.course_name || 'Unknown Course') + '<br>'
                    results += '<strong>Progress:</strong> ' + escapeHtml(details.progress || '0') + '%<br>'
                    results += '<strong>Course ID:</strong> ' + escapeHtml(details.course_id || 'N/A')
                    results += '</div></div>'
                }
                if (event.message == "Submitted Data") {
                    results += '<div class="timeline-replay-button"><button onclick="replay(' + i + ')" class="btn btn-success">'
                    results += '<i class="fa fa-refresh"></i> Replay Credentials</button></div>'
                    results += '<div class="timeline-event-details"><i class="fa fa-caret-right"></i> View Details</div>'
                }
                if (details.payload) {
                    results += '<div class="timeline-event-results">'
                    results += '    <table class="table table-condensed table-bordered table-striped">'
                    results += '        <thead><tr><th>Parameter</th><th>Value(s)</tr></thead><tbody>'
                    $.each(Object.keys(details.payload), function (i, param) {
                        if (param == "rid") {
                            return true;
                        }
                        results += '    <tr>'
                        results += '        <td>' + escapeHtml(param) + '</td>'
                        results += '        <td>' + escapeHtml(details.payload[param]) + '</td>'
                        results += '    </tr>'
                    })
                    results += '       </tbody></table>'
                    results += '</div>'
                }
                if (details.error) {
                    results += '<div class="timeline-event-details"><i class="fa fa-caret-right"></i> View Details</div>'
                    results += '<div class="timeline-event-results">'
                    results += '<span class="label label-default">Error</span> ' + details.error
                    results += '</div>'
                }
            }
            results += '</div></div>'
        }
    })
    // Add the scheduled send event at the bottom
    if (record.status == "Scheduled" || record.status == "Retrying") {
        results += '<div class="timeline-entry">' +
            '    <div class="timeline-bar"></div>'
        results +=
            '    <div class="timeline-icon ' + statuses[record.status].label + '">' +
            '    <i class="fa ' + statuses[record.status].icon + '"></i></div>' +
            '    <div class="timeline-message">' + "Scheduled to send at " + record.send_date + '</span>'
    }
    results += '</div></div>'
    return results
}

var renderTimelineChart = function (chartopts) {
    return Highcharts.chart('timeline_chart', {
        chart: {
            zoomType: 'x',
            type: 'line',
            height: "200px",
            backgroundColor: 'transparent',
            style: {
                fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, sans-serif"
            }
        },
        title: {
            text: null
        },
        xAxis: {
            type: 'datetime',
            dateTimeLabelFormats: {
                second: '%l:%M:%S',
                minute: '%l:%M',
                hour: '%l:%M',
                day: '%b %d, %Y',
                week: '%b %d, %Y',
                month: '%b %Y'
            },
            lineColor: '#e2e8f0',
            tickColor: '#e2e8f0',
            labels: {
                style: {
                    color: '#718096'
                }
            }
        },
        yAxis: {
            min: 0,
            max: 2,
            visible: false,
            tickInterval: 1,
            labels: {
                enabled: false
            },
            title: {
                text: ""
            }
        },
        tooltip: {
            backgroundColor: 'rgba(45, 55, 72, 0.95)',
            borderColor: '#667eea',
            borderRadius: 8,
            style: {
                color: '#ffffff'
            },
            formatter: function () {
                return Highcharts.dateFormat('%A, %b %d %l:%M:%S %P', new Date(this.x)) +
                    '<br>Event: ' + this.point.message + '<br>Email: <b>' + this.point.email + '</b>'
            }
        },
        legend: {
            enabled: false
        },
        plotOptions: {
            series: {
                marker: {
                    enabled: true,
                    symbol: 'circle',
                    radius: 5,
                    lineWidth: 2,
                    lineColor: '#ffffff'
                },
                cursor: 'pointer',
            },
            line: {
                states: {
                    hover: {
                        lineWidth: 2
                    }
                }
            }
        },
        credits: {
            enabled: false
        },
        series: [{
            data: chartopts['data'],
            dashStyle: "shortdash",
            color: "#667eea",
            lineWidth: 2,
            turboThreshold: 0
        }]
    })
}

/* Renders a mini pie chart in stat cards */
var renderPieChart = function (chartopts) {
    // Update the count display
    var countId = chartopts['elemId'].replace('_chart', '-count');
    $('#' + countId).text(chartopts['data'][0].count);

    return Highcharts.chart(chartopts['elemId'], {
        chart: {
            type: 'pie',
            backgroundColor: 'transparent',
            height: 80,
            width: 80,
            margin: [0, 0, 0, 0],
            spacing: [0, 0, 0, 0]
        },
        title: {
            text: null
        },
        credits: {
            enabled: false
        },
        tooltip: {
            enabled: false
        },
        plotOptions: {
            pie: {
                innerSize: '70%',
                dataLabels: {
                    enabled: false
                },
                states: {
                    hover: {
                        enabled: false
                    }
                }
            }
        },
        series: [{
            data: chartopts['data'],
            colors: chartopts['colors'],
        }]
    })
}

/* Renders the big combined stats chart */
var renderCombinedStatsChart = function (statsData, total, hasLearning, learningStats) {
    var categories = ['Sent', 'Opened', 'Clicked', 'Submitted', 'Reported'];
    var colors = ['#1abc9c', '#f9bf3b', '#F39C12', '#f05b4f', '#45d6ef'];
    var data = [
        { name: 'Sent', y: statsData['Email Sent'] || 0, color: '#1abc9c' },
        { name: 'Opened', y: statsData['Email Opened'] || 0, color: '#f9bf3b' },
        { name: 'Clicked', y: statsData['Clicked Link'] || 0, color: '#F39C12' },
        { name: 'Submitted', y: statsData['Submitted Data'] || 0, color: '#f05b4f' },
        { name: 'Reported', y: statsData['Email Reported'] || 0, color: '#45d6ef' }
    ];

    // Add learning stats if available
    if (hasLearning && learningStats) {
        data.push({
            name: 'Completed',
            y: learningStats.completed_count || 0,
            color: '#27ae60'
        });
    }

    return Highcharts.chart('combined_stats_chart', {
        chart: {
            type: 'bar',
            backgroundColor: 'transparent',
            height: 220,
            style: {
                fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, sans-serif"
            }
        },
        title: { text: null },
        xAxis: {
            categories: data.map(function(d) { return d.name; }),
            labels: {
                style: { color: '#718096', fontSize: '11px' }
            },
            lineColor: '#e2e8f0'
        },
        yAxis: {
            min: 0,
            max: total > 0 ? total : 1,
            title: { text: null },
            labels: {
                style: { color: '#718096' }
            },
            gridLineColor: '#e2e8f0'
        },
        tooltip: {
            backgroundColor: 'rgba(45, 55, 72, 0.95)',
            borderColor: '#667eea',
            borderRadius: 8,
            style: { color: '#ffffff' },
            formatter: function () {
                var pct = total > 0 ? Math.round((this.y / total) * 100) : 0;
                return '<b>' + this.point.name + '</b><br/>' +
                       this.y + ' of ' + total + ' (' + pct + '%)';
            }
        },
        legend: { enabled: false },
        credits: { enabled: false },
        plotOptions: {
            bar: {
                borderRadius: 4,
                dataLabels: {
                    enabled: true,
                    format: '{y}',
                    style: {
                        color: '#718096',
                        fontWeight: '600',
                        textOutline: 'none'
                    }
                }
            }
        },
        series: [{
            name: 'Count',
            data: data,
            colorByPoint: true
        }]
    });
}

/* Updates the bubbles on the map

@param {campaign.result[]} results - The campaign results to process
*/
var updateMap = function (results) {
    if (!map) {
        return
    }
    bubbles = []
    $.each(campaign.results, function (i, result) {
        // Check that it wasn't an internal IP
        if (result.latitude == 0 && result.longitude == 0) {
            return true;
        }
        newIP = true
        $.each(bubbles, function (i, bubble) {
            if (bubble.ip == result.ip) {
                bubbles[i].radius += 1
                newIP = false
                return false
            }
        })
        if (newIP) {
            bubbles.push({
                latitude: result.latitude,
                longitude: result.longitude,
                name: result.ip,
                fillKey: "point",
                radius: 2
            })
        }
    })
    map.bubbles(bubbles)
}

/**
 * Creates a status label for use in the results datatable
 * @param {string} status 
 * @param {moment(datetime)} send_date 
 */
function createStatusLabel(status, send_date) {
    var label = statuses[status].label || "label-default";
    var statusColumn = "<span class=\"label " + label + "\">" + status + "</span>"
    // Add the tooltip if the email is scheduled to be sent
    if (status == "Scheduled" || status == "Retrying") {
        var sendDateMessage = "Scheduled to send at " + send_date
        statusColumn = "<span class=\"label " + label + "\" data-toggle=\"tooltip\" data-placement=\"top\" data-html=\"true\" title=\"" + sendDateMessage + "\">" + status + "</span>"
    }
    return statusColumn
}

/**
 * Creates a learning status label for use in the results datatable
 * @param {string} status - Learning status (not-needed, not-started, in-progress, completed)
 * @param {number} progress - Progress percentage (0-100)
 */
function createLearningLabel(status, progress) {
    var statusInfo = learningStatuses[status] || learningStatuses["not-needed"];
    var label = statusInfo.label;
    var text = statusInfo.text;

    // For in-progress, show progress percentage
    if (status === "in-progress" && progress > 0) {
        text = progress + "% Done";
    }

    return "<span class=\"label " + label + "\"><i class=\"fa " + statusInfo.icon + "\"></i> " + text + "</span>";
}

/**
 * Renders the course completion pie chart
 */
function renderCourseCompletionChart(learningStats) {
    if (!learningStats) {
        return;
    }

    // Show the card and learning column (show even if no one is phished yet)
    $("#course-completion-card").show();
    $(".learning-column").show();

    // Update the count display
    $("#course-completion-count").text(learningStats.completed_count + "/" + learningStats.phished_count);

    // Calculate percentage (handle division by zero)
    var completedPercent = 0;
    if (learningStats.phished_count > 0) {
        completedPercent = Math.floor((learningStats.completed_count / learningStats.phished_count) * 100);
    }
    var chartData = [
        { name: "Completed", y: completedPercent, count: learningStats.completed_count },
        { name: "", y: 100 - completedPercent }
    ];

    return Highcharts.chart('course_completion_chart', {
        chart: {
            type: 'pie',
            backgroundColor: 'transparent',
            height: 80,
            width: 80,
            margin: [0, 0, 0, 0],
            spacing: [0, 0, 0, 0]
        },
        title: { text: null },
        credits: { enabled: false },
        tooltip: { enabled: false },
        plotOptions: {
            pie: {
                innerSize: '70%',
                dataLabels: { enabled: false },
                states: { hover: { enabled: false } }
            }
        },
        series: [{
            data: chartData,
            colors: ['#27ae60', '#dddddd']
        }]
    });
}

/* poll - Queries the API and updates the UI with the results
 *
 * Updates:
 * * Timeline Chart
 * * Email (Donut) Chart
 * * Map Bubbles
 * * Datatables
 */
function poll() {
    api.campaignId.results(campaign.id)
        .success(function (c) {
            campaign = c
            /* Update the timeline */
            var timeline_series_data = []
            $.each(campaign.timeline, function (i, event) {
                var event_date = moment.utc(event.time).local()
                timeline_series_data.push({
                    email: event.email,
                    message: event.message,
                    x: event_date.valueOf(),
                    y: 1,
                    marker: {
                        fillColor: statuses[event.message].color
                    }
                })
            })
            var timeline_chart = $("#timeline_chart").highcharts()
            timeline_chart.series[0].update({
                data: timeline_series_data
            })
            /* Update the results donut chart */
            var email_series_data = {}
            // Load the initial data
            Object.keys(statusMapping).forEach(function (k) {
                email_series_data[k] = 0
            });
            $.each(campaign.results, function (i, result) {
                email_series_data[result.status]++;
                if (result.reported) {
                    email_series_data['Email Reported']++
                }
                // Backfill status values
                var step = progressListing.indexOf(result.status)
                for (var i = 0; i < step; i++) {
                    email_series_data[progressListing[i]]++
                }
            })
            $.each(email_series_data, function (status, count) {
                var email_data = []
                if (!(status in statusMapping)) {
                    return true
                }
                email_data.push({
                    name: status,
                    y: Math.floor((count / campaign.results.length) * 100),
                    count: count
                })
                email_data.push({
                    name: '',
                    y: 100 - Math.floor((count / campaign.results.length) * 100)
                })
                var chart = $("#" + statusMapping[status] + "_chart").highcharts()
                chart.series[0].update({
                    data: email_data
                })
            })

            /* Update the course completion chart if learning stats exist */
            if (campaign.learning_stats) {
                var completionChart = $("#course_completion_chart").highcharts();
                if (completionChart) {
                    var completedPercent = 0;
                    if (campaign.learning_stats.phished_count > 0) {
                        completedPercent = Math.floor((campaign.learning_stats.completed_count / campaign.learning_stats.phished_count) * 100);
                    }
                    var chartData = [
                        { name: "Completed", y: completedPercent, count: campaign.learning_stats.completed_count },
                        { name: "", y: 100 - completedPercent }
                    ];
                    completionChart.series[0].update({ data: chartData });
                    $("#course-completion-count").text(campaign.learning_stats.completed_count + "/" + campaign.learning_stats.phished_count);
                }
            }

            /* Update the combined stats chart */
            var combinedChart = $("#combined_stats_chart").highcharts();
            if (combinedChart) {
                var hasLearning = campaign.learning_stats !== null && campaign.learning_stats !== undefined;
                var newData = [
                    { name: 'Sent', y: email_series_data['Email Sent'] || 0, color: '#1abc9c' },
                    { name: 'Opened', y: email_series_data['Email Opened'] || 0, color: '#f9bf3b' },
                    { name: 'Clicked', y: email_series_data['Clicked Link'] || 0, color: '#F39C12' },
                    { name: 'Submitted', y: email_series_data['Submitted Data'] || 0, color: '#f05b4f' },
                    { name: 'Reported', y: email_series_data['Email Reported'] || 0, color: '#45d6ef' }
                ];
                if (hasLearning && campaign.learning_stats) {
                    newData.push({ name: 'Completed', y: campaign.learning_stats.completed_count || 0, color: '#27ae60' });
                }
                combinedChart.series[0].setData(newData, true);
                combinedChart.yAxis[0].setExtremes(0, campaign.results.length > 0 ? campaign.results.length : 1);
            }

            /* Update the datatable */
            resultsTable = $("#resultsTable").DataTable()
            resultsTable.rows().every(function (i, tableLoop, rowLoop) {
                var row = this.row(i)
                var rowData = row.data()
                var rid = rowData[0]
                var email = rowData[4] // email is at index 4
                $.each(campaign.results, function (j, result) {
                    if (result.id == rid) {
                        rowData[8] = moment(result.send_date).format('MMMM Do YYYY, h:mm:ss a')
                        rowData[7] = result.reported
                        rowData[6] = result.status
                        // Update learning progress
                        if (campaign.learning_progress && campaign.learning_progress[result.email]) {
                            rowData[9] = campaign.learning_progress[result.email];
                        }
                        resultsTable.row(i).data(rowData)
                        if (row.child.isShown()) {
                            $(row.node()).find("#caret").removeClass("fa-caret-right")
                            $(row.node()).find("#caret").addClass("fa-caret-down")
                            row.child(renderTimeline(row.data()))
                        }
                        return false
                    }
                })
            })
            resultsTable.draw(false)
            /* Update the map information */
            updateMap(campaign.results)
            $('[data-toggle="tooltip"]').tooltip()
            $("#refresh_message").hide()
            $("#refresh_btn").show()
        })
}

function load() {
    campaign.id = window.location.pathname.split('/').slice(-1)[0]
    var use_map = JSON.parse(localStorage.getItem('gophish.use_map'))
    api.campaignId.results(campaign.id)
        .success(function (c) {
            campaign = c
            if (campaign) {
                $("title").text(c.name + " - Gophish")
                $("#loading").hide()
                $("#campaignResults").show()
                // Set the title
                $("#page-title").text("Results for " + c.name)
                if (c.status == "Completed") {
                    $('#complete_button')[0].disabled = true;
                    $('#complete_button').text('Completed!');
                    doPoll = false;
                }
                // Setup viewing the details of a result
                $("#resultsTable").on("click", ".timeline-event-details", function () {
                    // Show the parameters
                    payloadResults = $(this).parent().find(".timeline-event-results")
                    if (payloadResults.is(":visible")) {
                        $(this).find("i").removeClass("fa-caret-down")
                        $(this).find("i").addClass("fa-caret-right")
                        payloadResults.hide()
                    } else {
                        $(this).find("i").removeClass("fa-caret-right")
                        $(this).find("i").addClass("fa-caret-down")
                        payloadResults.show()
                    }
                })
                // Check if campaign has learning stats (course associated)
                var hasLearningStats = campaign.learning_stats !== null && campaign.learning_stats !== undefined;

                // Setup the results table
                var columnDefsList = [{
                        orderable: false,
                        targets: "no-sort"
                    }, {
                        className: "details-control",
                        "targets": [1]
                    }, {
                        "visible": false,
                        "targets": [0, 8]
                    },
                    {
                        "render": function (data, type, row) {
                            return createStatusLabel(data, row[8])
                        },
                        "targets": [6]
                    },
                    {
                        className: "text-center",
                        "render": function (reported, type, row) {
                            if (type == "display") {
                                if (reported) {
                                    return "<i class='fa fa-check-circle text-center text-success'></i>"
                                }
                                return "<i role='button' class='fa fa-times-circle text-center text-muted' onclick='report_mail(\"" + row[0] + "\", \"" + campaign.id + "\");'></i>"
                            }
                            return reported
                        },
                        "targets": [7]
                    }
                ];

                // Add learning column renderer if learning stats exist
                if (hasLearningStats) {
                    columnDefsList.push({
                        className: "text-center",
                        "render": function (data, type, row) {
                            if (type == "display" && data) {
                                return createLearningLabel(data.status, data.progress);
                            }
                            return data ? data.status : "";
                        },
                        "targets": [9]
                    });
                } else {
                    // Hide learning column if no learning stats
                    columnDefsList.push({
                        "visible": false,
                        "targets": [9]
                    });
                }

                resultsTable = $("#resultsTable").DataTable({
                    destroy: true,
                    "order": [
                        [2, "asc"]
                    ],
                    columnDefs: columnDefsList
                });
                resultsTable.clear();
                var email_series_data = {}
                var timeline_series_data = []
                Object.keys(statusMapping).forEach(function (k) {
                    email_series_data[k] = 0
                });
                $.each(campaign.results, function (i, result) {
                    // Get learning progress for this result
                    var learningProgress = null;
                    if (campaign.learning_progress && campaign.learning_progress[result.email]) {
                        learningProgress = campaign.learning_progress[result.email];
                    }

                    resultsTable.row.add([
                        result.id,
                        "<i id=\"caret\" class=\"fa fa-caret-right\"></i>",
                        escapeHtml(result.first_name) || "",
                        escapeHtml(result.last_name) || "",
                        escapeHtml(result.email) || "",
                        escapeHtml(result.position) || "",
                        result.status,
                        result.reported,
                        moment(result.send_date).format('MMMM Do YYYY, h:mm:ss a'),
                        learningProgress
                    ])
                    email_series_data[result.status]++;
                    if (result.reported) {
                        email_series_data['Email Reported']++
                    }
                    // Backfill status values
                    var step = progressListing.indexOf(result.status)
                    for (var i = 0; i < step; i++) {
                        email_series_data[progressListing[i]]++
                    }
                })
                resultsTable.draw();
                // Setup tooltips
                $('[data-toggle="tooltip"]').tooltip()
                // Setup the individual timelines
                $('#resultsTable tbody').on('click', 'td.details-control', function () {
                    var tr = $(this).closest('tr');
                    var row = resultsTable.row(tr);
                    if (row.child.isShown()) {
                        // This row is already open - close it
                        row.child.hide();
                        tr.removeClass('shown');
                        $(this).find("i").removeClass("fa-caret-down")
                        $(this).find("i").addClass("fa-caret-right")
                    } else {
                        // Open this row
                        $(this).find("i").removeClass("fa-caret-right")
                        $(this).find("i").addClass("fa-caret-down")
                        row.child(renderTimeline(row.data())).show();
                        tr.addClass('shown');
                    }
                });
                // Setup the graphs
                $.each(campaign.timeline, function (i, event) {
                    if (event.message == "Campaign Created") {
                        return true
                    }
                    var event_date = moment.utc(event.time).local()
                    timeline_series_data.push({
                        email: event.email,
                        message: event.message,
                        x: event_date.valueOf(),
                        y: 1,
                        marker: {
                            fillColor: statuses[event.message].color
                        }
                    })
                })
                renderTimelineChart({
                    data: timeline_series_data
                })
                $.each(email_series_data, function (status, count) {
                    var email_data = []
                    if (!(status in statusMapping)) {
                        return true
                    }
                    email_data.push({
                        name: status,
                        y: Math.floor((count / campaign.results.length) * 100),
                        count: count
                    })
                    email_data.push({
                        name: '',
                        y: 100 - Math.floor((count / campaign.results.length) * 100)
                    })
                    var chart = renderPieChart({
                        elemId: statusMapping[status] + '_chart',
                        title: status,
                        name: status,
                        data: email_data,
                        colors: [statuses[status].color, '#dddddd']
                    })
                })

                // Render course completion chart if learning stats exist
                if (hasLearningStats) {
                    renderCourseCompletionChart(campaign.learning_stats);
                }

                // Render the combined stats chart
                renderCombinedStatsChart(email_series_data, campaign.results.length, hasLearningStats, campaign.learning_stats);

                if (use_map) {
                    $("#resultsMapContainer").show()
                    map = new Datamap({
                        element: document.getElementById("resultsMap"),
                        responsive: true,
                        fills: {
                            defaultFill: "#ffffff",
                            point: "#283F50"
                        },
                        geographyConfig: {
                            highlightFillColor: "#1abc9c",
                            borderColor: "#283F50"
                        },
                        bubblesConfig: {
                            borderColor: "#283F50"
                        }
                    });
                }
                updateMap(campaign.results)
            }
        })
        .error(function () {
            $("#loading").hide()
            errorFlash(" Campaign not found!")
        })
}

var setRefresh

function refresh() {
    if (!doPoll) {
        return;
    }
    $("#refresh_message").show()
    $("#refresh_btn").hide()
    poll()
    clearTimeout(setRefresh)
    setRefresh = setTimeout(refresh, 60000)
};

function report_mail(rid, cid) {
    Swal.fire({
        title: "Are you sure?",
        text: "This result will be flagged as reported (RID: " + rid + ")",
        type: "question",
        animation: false,
        showCancelButton: true,
        confirmButtonText: "Continue",
        confirmButtonColor: "#428bca",
        reverseButtons: true,
        allowOutsideClick: false,
        showLoaderOnConfirm: true
    }).then(function (result) {
        if (result.value){
            api.campaignId.get(cid).success((function(c) {
                report_url = new URL(c.url)
                report_url.pathname = '/report'
                report_url.search = "?rid=" + rid 
                fetch(report_url)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    refresh();
                })
                .catch(error => {
                    let errorMessage = error.message;
                    if (error.message === "Failed to fetch") {
                        errorMessage = "This might be due to Mixed Content issues or network problems.";
                    }
                    Swal.fire({
                        title: 'Error',
                        text: errorMessage,
                        type: 'error',
                        confirmButtonText: 'Close'
                    });
                });
            }));
        }
    })
}

$(document).ready(function () {
    Highcharts.setOptions({
        global: {
            useUTC: false
        }
    })
    load();

    // Start the polling loop
    setRefresh = setTimeout(refresh, 60000)
})
