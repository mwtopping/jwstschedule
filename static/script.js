let customMinDate = null;
let customMaxDate = null;
let currentTimeScale = 'days';

function generateTimelines() {
    const rows = document.querySelectorAll('#timeline-table tbody tr');

    let minDate, maxDate;

    minDate = new Date(document.getElementById('min-datetime').value)
    maxDate = new Date(document.getElementById('max-datetime').value)
    console.log(minDate, maxDate);

    const totalMilliseconds = maxDate - minDate;

    function generateHashMarks(showTimes) {
        let hashMarksHTML = '';
        let numMarks = 4;
        let formatOptions = {};
        formatOptions = { year: 'numeric', month: 'short', day: 'numeric', hour: 'numeric', minute: 'numeric'};

        for (let i = 0; i <= numMarks; i++) {
            const percentage = (i / numMarks) * 100;
            const markDate = new Date(minDate.getTime() + (totalMilliseconds * (i / numMarks)));
            console.log(percentage, markDate, minDate);

            let translateAmount = 0;
            if (i > 0) {
                translateAmount = -50;
            }
            if (i === numMarks) {
                translateAmount = -100;
            }


            //*transform: translateX(-50%);*/
            const labelClass = 'hash-label';
            if (showTimes === true) {
             hashMarksHTML += `
                    <div class="hash-mark" style="left: ${percentage}%;">
                        <div class="hash-line"></div>
                        <div class="${labelClass}" style="transform: translateX(${translateAmount}%);">${markDate.toLocaleDateString('en-US', formatOptions)}</div>
                    </div>
                `;



            } else {
                hashMarksHTML += `
                    <div class="hash-mark" style="left: ${percentage}%;">
                        <div class="hash-line"></div>
                    </div>
                `;

            }
        }

        return hashMarksHTML;
    }

    const occurance = 10;
    ii = 0;

    rows.forEach(row => {

        row.style.display = "";
        const startDate = new Date(row.dataset.start);
        const endDate = new Date(row.dataset.end);
        const eapMonths = row.dataset.eap;
        const status = row.dataset.status;
        const timelineCell = row.querySelector('.timeline-cell');
        let color = "#000000";

        switch (status) {
            case "Executed":
            case "Completed":
            case "Archived":
                color = "#31a354";
                break;
            case "Collecting":
                color = "#addd8e";
                break;
            case "Implementation":
            case "Pending Submission":
                color = "#bdd7e7";
                break;
            case "Flight Ready":
            case "Scheduled":
                color = "#3182bd";
                break;
            case "Failed":
            case "Skipped":
            case "Withdrawn":
                color = "#fb6a4a";
                break;
            default:
                color = "#555555";
                break;
        }


        if (endDate - minDate < -1*24*60*60*1000) {
            row.style.display = "none";
        }



        const startOffset = Math.max(0, ((startDate - minDate) / totalMilliseconds) * 100);

        const startEap = new Date(startDate.getTime() + eapMonths*30*24*60*60*1000);
        const eapOffset = Math.max(0, ((startEap - minDate) / totalMilliseconds) * 100);

        const endOffset = Math.min(100, ((endDate - minDate) / totalMilliseconds) * 100);
        const duration = endOffset - startOffset;

        let eapDisplay = "";
        if (eapOffset > 100) {
            eapDisplay = "none"
        }

        if (eapOffset >= 0) {
            row.style.display = "";
        }




        if (ii%occurance ===0) {
        timelineCell.innerHTML = `
            <div class="timeline-track"></div>
            <div class="timeline-segment" style="background: ${color}; left: ${startOffset}%; width: ${duration}%;"></div>
            <div class="eap-segment" style="left: ${eapOffset}%;"></div>
            ${generateHashMarks(true)}
        `;
        } else {
        timelineCell.innerHTML = `
            <div class="timeline-track"></div>
            <div class="timeline-segment" style="background: ${color}; left: ${startOffset}%; width: ${duration}%;"></div>
            <div class="eap-segment" style="display: ${eapDisplay}; left: ${eapOffset}%;"></div>
            ${generateHashMarks(false)}
        `;

        }
        ii += 1;

    });
}

function updateTimelineRange() {
    const minDateInput = document.getElementById('min-datetime');
    const maxDateInput = document.getElementById('max-datetime');

    const minValue = minDateInput.value;
    const maxValue = maxDateInput.value;

    if (minValue && maxValue) {
        customMinDate = new Date(minValue);
        customMaxDate = new Date(maxValue);

        if (customMinDate >= customMaxDate) {
            alert('Start date/time must be before end date/time!');
            return;
        }

        generateTimelines();
    }
}

function autoFitToData(customMinDate = null, customMaxDate = null) {

    const rows = document.querySelectorAll('#timeline-table tbody tr');
    let minDate = new Date('2099-12-31T23:59:59');
    let maxDate = new Date('1900-01-01T00:00:00');

    rows.forEach(row => {
        const startDate = new Date(row.dataset.start);
        const endDate = new Date(row.dataset.end);
        if (startDate < minDate) minDate = startDate;
        if (endDate > maxDate) maxDate = endDate;
    });

    // (date1 > date2 ? date1 : date2)
    let finalMinDate = customMinDate === null ? minDate : customMinDate;
    let finalMaxDate = customMaxDate === null ? maxDate : customMaxDate;

    console.log(minDate, customMinDate, finalMinDate);
    document.getElementById('min-datetime').value = finalMinDate.toISOString().slice(0, 16);
    document.getElementById('max-datetime').value = finalMaxDate.toISOString().slice(0, 16);
    
    console.log(finalMinDate, finalMinDate.toISOString().slice(0, 16));
    generateTimelines();
}

function updateTimeScale() {
    currentTimeScale = document.getElementById('time-scale').value;
    generateTimelines();
}

document.addEventListener('DOMContentLoaded', function() {
    if (window.location.pathname.endsWith('/week')) {
        let today = new Date();
        let oneweek = new Date(today.getTime() + 7*24*60*60*1000);
        console.log("today:", today)
        autoFitToData(customMinDate = today, customMaxDate = oneweek);
        let pagebutton = document.getElementById('weeknav');
        pagebutton.className = "navselected";
    } else if (window.location.pathname.endsWith('/month')) {
        let today = new Date();
        let onemonth = new Date(today.getTime() + 30*24*60*60*1000);
        console.log("today:", today)
        autoFitToData(customMinDate = today, customMaxDate = onemonth);
        let pagebutton = document.getElementById('monthnav');
        pagebutton.className = "navselected";

    } else if (window.location.pathname.endsWith('/year')) {
        let today = new Date();
        let onemonth = new Date(today.getTime() + 365.25*24*60*60*1000);
        console.log("today:", today)
        autoFitToData(customMinDate = today, customMaxDate = onemonth);
        let pagebutton = document.getElementById('yearnav');
        pagebutton.className = "navselected";
    } else if (window.location.pathname.endsWith('/all')) {
        autoFitToData();
        let pagebutton = document.getElementById('allnav');
        pagebutton.className = "navselected";
    } else {
        let today = new Date();
        let oneweek = new Date(today.getTime() + 7*24*60*60*1000);
        console.log("today:", today)
        autoFitToData(customMinDate = today, customMaxDate = oneweek);
        let pagebutton = document.getElementById('weeknav');
        pagebutton.className = "navselected";
    }

    document.getElementById('auto-fit').addEventListener('click', autoFitToData);

    document.getElementById('min-datetime').addEventListener('change', updateTimelineRange);
    document.getElementById('max-datetime').addEventListener('change', updateTimelineRange);
});
