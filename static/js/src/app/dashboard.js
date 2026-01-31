var campaigns = []

// statuses is a helper map to point result statuses to ui classes
var statuses = {
    "Email Sent": {
        color: "#1abc9c",
        label: "label-success-modern",
        icon: "fa-envelope",
        point: "ct-point-sent"
    },
    "Emails Sent": {
        color: "#1abc9c",
        label: "label-success-modern",
        icon: "fa-envelope",
        point: "ct-point-sent"
    },
    "In progress": {
        label: "label-primary-modern"
    },
    "Queued": {
        label: "label-info-modern"
    },
    "Completed": {
        label: "label-success-modern"
    },
    "Email Opened": {
        color: "#f9bf3b",
        label: "label-warning-modern",
        icon: "fa-envelope-open",
        point: "ct-point-opened"
    },
    "Email Reported": {
        color: "#45d6ef",
        label: "label-info-modern",
        icon: "fa-flag",
        point: "ct-point-reported"
    },
    "Clicked Link": {
        color: "#F39C12",
        label: "label-warning-modern",
        icon: "fa-mouse-pointer",
        point: "ct-point-clicked"
    },
    "Success": {
        color: "#f05b4f",
        label: "label-danger-modern",
        icon: "fa-exclamation",
        point: "ct-point-clicked"
    },
    "Error": {
        color: "#6c7a89",
        label: "label-secondary-modern",
        icon: "fa-times",
        point: "ct-point-error"
    },
    "Error Sending Email": {
        color: "#6c7a89",
        label: "label-secondary-modern",
        icon: "fa-times",
        point: "ct-point-error"
    },
    "Submitted Data": {
        color: "#f05b4f",
        label: "label-danger-modern",
        icon: "fa-exclamation",
        point: "ct-point-clicked"
    },
    "Unknown": {
        color: "#6c7a89",
        label: "label-secondary-modern",
        icon: "fa-question",
        point: "ct-point-error"
    },
    "Sending": {
        color: "#428bca",
        label: "label-primary-modern",
        icon: "fa-spinner",
        point: "ct-point-sending"
    },
    "Campaign Created": {
        label: "label-success-modern",
        icon: "fa-rocket"
    }
}

var statsMapping = {
    "sent": "Email Sent",
    "opened": "Email Opened",
    "email_reported": "Email Reported",
    "clicked": "Clicked Link",
    "submitted_data": "Submitted Data",
}

// Color mapping for stat cards
var statColors = {
    "sent": "#1abc9c",
    "opened": "#f9bf3b",
    "clicked": "#F39C12",
    "submitted_data": "#f05b4f",
    "email_reported": "#45d6ef"
}

function deleteCampaign(campaignId) {
    var campaign = campaigns.find(function(c) { return c.id === campaignId; });
    if (!campaign) {
        errorFlash("Campaign not found");
        return;
    }

    Swal.fire({
        title: 'Delete Campaign?',
        text: "Are you sure you want to delete " + campaign.name + "?",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#667eea',
        cancelButtonColor: '#6c757d',
        confirmButtonText: 'Yes, delete it!'
    }).then((result) => {
        if (result.isConfirmed) {
            api.campaignId.delete(campaignId)
                .success(function (data) {
                    successFlash(data.message)
                    location.reload()
                })
                .error(function() {
                    errorFlash("Error deleting campaign")
                })
        }
    })
}

/* Renders a mini pie chart in stat cards */
function renderMiniPieChart(chartopts) {
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

function generateStatsPieCharts(campaigns) {
    var stats_series_data = {}
    var total = 0

    $.each(campaigns, function (i, campaign) {
        $.each(campaign.stats, function (status, count) {
            if (status == "total") {
                total += count
                return true
            }
            if (!stats_series_data[status]) {
                stats_series_data[status] = count;
            } else {
                stats_series_data[status] += count;
            }
        })
    })

    $.each(stats_series_data, function (status, count) {
        if (!(status in statsMapping)) {
            return true
        }

        // Update the count display
        $('#' + status + '-count').text(count)

        // Calculate percentage
        var percentage = total > 0 ? Math.floor((count / total) * 100) : 0

        // Render mini pie chart
        var stats_data = [
            { name: status, y: percentage },
            { name: '', y: 100 - percentage }
        ]

        renderMiniPieChart({
            elemId: status + '_chart',
            data: stats_data,
            colors: [statColors[status], "#e2e8f0"]
        })
    });

    // Render combined stats chart (will be called separately with learning stats)
    return { stats: stats_series_data, total: total };
}

/* Renders the big combined stats chart */
function renderCombinedStatsChart(statsData, total, learningStats) {
    // Always show all 6 bars, even with zero values
    var data = [
        { name: 'Sent', y: statsData['sent'] || 0, color: '#1abc9c' },
        { name: 'Opened', y: statsData['opened'] || 0, color: '#f9bf3b' },
        { name: 'Clicked', y: statsData['clicked'] || 0, color: '#F39C12' },
        { name: 'Submitted', y: statsData['submitted_data'] || 0, color: '#f05b4f' },
        { name: 'Reported', y: statsData['email_reported'] || 0, color: '#45d6ef' },
        { name: 'Completed', y: (learningStats && learningStats.completed_count) || 0, color: '#27ae60' }
    ];

    // Ensure yAxis max is at least 1 to show the grid
    var yMax = total > 0 ? total : 10;

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
            max: yMax,
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

function generateTimelineChart(campaigns) {
    if (!campaigns || campaigns.length === 0) {
        $('#overview_chart').html('<div class="empty-state" style="padding: 2rem;"><i class="fa fa-line-chart" style="font-size: 3rem; color: #cbd5e0; margin-bottom: 1rem;"></i><p style="color: #718096;">No campaign data to display yet</p></div>');
        return;
    }

    var overview_data = []
    $.each(campaigns, function (i, campaign) {
        var campaign_date = moment.utc(campaign.created_date).local()
        campaign.y = 0
        campaign.y += campaign.stats.clicked
        campaign.y = Math.floor((campaign.y / campaign.stats.total) * 100)
        overview_data.push({
            campaign_id: campaign.id,
            name: campaign.name,
            x: campaign_date.valueOf(),
            y: campaign.y
        })
    })

    Highcharts.chart('overview_chart', {
        chart: {
            zoomType: 'x',
            type: 'areaspline',
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
            max: 100,
            title: {
                text: "Success Rate (%)",
                style: {
                    color: '#718096'
                }
            },
            gridLineColor: '#e2e8f0',
            labels: {
                style: {
                    color: '#718096'
                }
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
                return '<b>' + this.point.name + '</b><br/>' +
                    Highcharts.dateFormat('%b %d, %Y', new Date(this.x)) +
                    '<br/>Success Rate: <b>' + this.y + '%</b>'
            }
        },
        legend: {
            enabled: false
        },
        plotOptions: {
            areaspline: {
                fillColor: {
                    linearGradient: {
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 1
                    },
                    stops: [
                        [0, 'rgba(102, 126, 234, 0.4)'],
                        [1, 'rgba(102, 126, 234, 0.05)']
                    ]
                },
                lineColor: '#667eea',
                lineWidth: 3,
                marker: {
                    enabled: true,
                    fillColor: '#667eea',
                    lineColor: '#ffffff',
                    lineWidth: 2,
                    radius: 5,
                    symbol: 'circle'
                }
            },
            series: {
                cursor: 'pointer',
                point: {
                    events: {
                        click: function (e) {
                            window.location.href = "/campaigns/" + this.campaign_id
                        }
                    }
                }
            }
        },
        credits: {
            enabled: false
        },
        series: [{
            data: overview_data,
            name: 'Success Rate'
        }]
    })
}

$(document).ready(function () {
    Highcharts.setOptions({
        global: {
            useUTC: false
        }
    })

    api.campaigns.summary()
        .success(function (data) {
            $("#loading").hide()
            campaigns = data.campaigns
            if (campaigns.length > 0) {
                $("#dashboard").show()

                // Initialize DataTable with modern styling
                var campaignTable = $("#campaignTable").DataTable({
                    columnDefs: [{
                            orderable: false,
                            targets: "no-sort"
                        },
                        {
                            className: "color-sent",
                            targets: [2]
                        },
                        {
                            className: "color-opened",
                            targets: [3]
                        },
                        {
                            className: "color-clicked",
                            targets: [4]
                        },
                        {
                            className: "color-success",
                            targets: [5]
                        },
                        {
                            className: "color-reported",
                            targets: [6]
                        }
                    ],
                    order: [
                        [1, "desc"]
                    ],
                    language: {
                        emptyTable: "No campaigns found",
                        info: "Showing _START_ to _END_ of _TOTAL_ campaigns",
                        infoEmpty: "No campaigns to show",
                        infoFiltered: "(filtered from _MAX_ total campaigns)",
                        lengthMenu: "Show _MENU_ campaigns",
                        search: "Search:",
                        paginate: {
                            first: "First",
                            last: "Last",
                            next: "Next",
                            previous: "Previous"
                        }
                    },
                    pageLength: 10,
                    dom: '<"row"<"col-sm-6"l><"col-sm-6"f>>rt<"row"<"col-sm-6"i><"col-sm-6"p>>'
                });

                var campaignRows = []
                $.each(campaigns, function (i, campaign) {
                    var campaign_date = moment(campaign.created_date).format('MMM D, YYYY h:mm A')
                    var label = statuses[campaign.status].label || "label-secondary-modern";

                    // Quick stats tooltip
                    var launchDate;
                    if (moment(campaign.launch_date).isAfter(moment())) {
                        launchDate = "Scheduled: " + moment(campaign.launch_date).format('MMM D, YYYY h:mm A')
                        var quickStats = launchDate + "<br>Recipients: " + campaign.stats.total
                    } else {
                        launchDate = "Launched: " + moment(campaign.launch_date).format('MMM D, YYYY h:mm A')
                        var quickStats = launchDate +
                            "<br>Recipients: " + campaign.stats.total +
                            "<br>Opened: " + campaign.stats.opened +
                            "<br>Clicked: " + campaign.stats.clicked +
                            "<br>Submitted: " + campaign.stats.submitted_data +
                            "<br>Errors: " + campaign.stats.error +
                            "<br>Reported: " + campaign.stats.email_reported
                    }

                    campaignRows.push([
                        escapeHtml(campaign.name),
                        campaign_date,
                        campaign.stats.sent,
                        campaign.stats.opened,
                        campaign.stats.clicked,
                        campaign.stats.submitted_data,
                        campaign.stats.email_reported,
                        '<span class="label-modern ' + label + '" data-toggle="tooltip" data-placement="top" data-html="true" title="' + quickStats + '">' + campaign.status + '</span>',
                        '<div class="action-buttons">' +
                            '<a class="btn-modern btn-primary-modern btn-sm-modern" href="/campaigns/' + campaign.id + '" aria-label="View results for ' + escapeHtml(campaign.name) + '" data-toggle="tooltip" data-placement="top" title="View Results">' +
                                '<i class="fa fa-bar-chart" aria-hidden="true"></i>' +
                            '</a>' +
                            '<button class="btn-modern btn-danger-modern btn-sm-modern" onclick="deleteCampaign(' + campaign.id + ')" aria-label="Delete ' + escapeHtml(campaign.name) + '" data-toggle="tooltip" data-placement="top" title="Delete Campaign">' +
                                '<i class="fa fa-trash" aria-hidden="true"></i>' +
                            '</button>' +
                        '</div>'
                    ])
                })

                campaignTable.rows.add(campaignRows).draw()

                // Initialize tooltips
                $('[data-toggle="tooltip"]').tooltip()

                // Build the charts
                var statsResult = generateStatsPieCharts(campaigns)
                generateTimelineChart(campaigns)

                // Update course completion card with learning stats
                var learningStats = data.learning_stats
                if (learningStats && learningStats.phished_count > 0) {
                    $('#course-completion-count').text(learningStats.completed_count + '/' + learningStats.phished_count)
                } else {
                    $('#course-completion-count').text('0/0')
                }

                // Render combined stats chart with learning data
                renderCombinedStatsChart(statsResult.stats, statsResult.total, learningStats)
            } else {
                $("#emptyMessage").show()
            }
        })
        .error(function () {
            $("#loading").hide()
            errorFlash("Error fetching campaigns")
        })
})
