let customMinDate = null;
let customMaxDate = null;
let currentTimeScale = 'days';

function generateTimelines() {
    const rows = document.querySelectorAll('#timeline-table tbody tr');

    let minDate, maxDate;

    if (customMinDate && customMaxDate) {
        minDate = customMinDate;
        maxDate = customMaxDate;
    } else {
        // Determine min and max date from data
        minDate = new Date('2099-12-31T23:59:59');
        maxDate = new Date('1900-01-01T00:00:00');

        rows.forEach(row => {
            const startDate = new Date(row.dataset.start);
            const endDate = new Date(row.dataset.end);
            if (startDate < minDate) minDate = startDate;
            if (endDate > maxDate) maxDate = endDate;
        });
    }

    const totalMilliseconds = maxDate - minDate;

    function generateHashMarks(showTimes) {
        let hashMarksHTML = '';
        let numMarks = 4;
        let formatOptions = {};
        formatOptions = { year: 'numeric', month: 'short', day: 'numeric' };

        for (let i = 0; i <= numMarks; i++) {
            const percentage = (i / numMarks) * 100;
            const markDate = new Date(minDate.getTime() + (totalMilliseconds * (i / numMarks)));

            const labelClass = 'hash-label';
            if (showTimes === true) {
             hashMarksHTML += `
                    <div class="hash-mark" style="left: ${percentage}%;">
                        <div class="hash-line"></div>
                        <div class="${labelClass}">${markDate.toLocaleDateString('en-US', formatOptions)}</div>
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
        const startDate = new Date(row.dataset.start);
        const endDate = new Date(row.dataset.end);
        const timelineCell = row.querySelector('.timeline-cell');

        const startOffset = Math.max(0, ((startDate - minDate) / totalMilliseconds) * 100);
        const endOffset = Math.min(100, ((endDate - minDate) / totalMilliseconds) * 100);
        const duration = endOffset - startOffset;

        if (ii%occurance ===0) {
        timelineCell.innerHTML = `
            <div class="timeline-track"></div>
            <div class="timeline-segment" style="left: ${startOffset}%; width: ${duration}%;"></div>
            ${generateHashMarks(true)}
        `;
        } else {
        timelineCell.innerHTML = `
            <div class="timeline-track"></div>
            <div class="timeline-segment" style="left: ${startOffset}%; width: ${duration}%;"></div>
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

function autoFitToData() {
    customMinDate = null;
    customMaxDate = null;

    const rows = document.querySelectorAll('#timeline-table tbody tr');
    let minDate = new Date('2099-12-31T23:59:59');
    let maxDate = new Date('1900-01-01T00:00:00');

    rows.forEach(row => {
        const startDate = new Date(row.dataset.start);
        const endDate = new Date(row.dataset.end);
        if (startDate < minDate) minDate = startDate;
        if (endDate > maxDate) maxDate = endDate;
    });

    document.getElementById('min-datetime').value = minDate.toISOString().slice(0, 16);
    document.getElementById('max-datetime').value = maxDate.toISOString().slice(0, 16);

    generateTimelines();
}

function updateTimeScale() {
    currentTimeScale = document.getElementById('time-scale').value;
    generateTimelines();
}

document.addEventListener('DOMContentLoaded', function() {
    autoFitToData();

    document.getElementById('auto-fit').addEventListener('click', autoFitToData);

    document.getElementById('min-datetime').addEventListener('change', updateTimelineRange);
    document.getElementById('max-datetime').addEventListener('change', updateTimelineRange);
});
