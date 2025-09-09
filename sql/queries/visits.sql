-- name: CreateVisit :one
INSERT INTO visits (id, created_at, updated_at, program_ID, observation, visit, Status, Target, Configuration, StartTime, EndTime)
VALUES (
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?
)
RETURNING *;


-- name: UpdateVisit :one
UPDATE visits
SET 
updated_at=?, 
Status=?, 
Target=?, 
Configuration=?, 
StartTime=?, 
EndTime=?
WHERE id = ?
RETURNING *;



-- name: GetAllVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
ORDER BY visits.StartTime;


-- name: GetWeekVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	visits.StartTime - ? BETWEEN 0 AND 60*60*24*7
ORDER BY visits.StartTime;

-- name: GetMonthVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	visits.StartTime - ? BETWEEN 0 AND 60*60*24*30
ORDER BY visits.StartTime;

-- name: GetYearVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	visits.StartTime - ? BETWEEN 0 AND 60*60*24*365
ORDER BY visits.StartTime;



-- name: GetPendingPrograms :many
SELECT 
	DISTINCT program_info.id
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	visits.Status NOT IN (
		"Archived",
		"Withdrawn",
		"Inactive",
		"Failed",
		"Skipped"
		)
ORDER BY program_info.id;


-- name: GetAllReleaseVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap,
	visits.StartTime + 30*60*60*24*program_info.eap as release_date
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
ORDER BY release_date;



-- name: GetWeekReleaseVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap,
	visits.StartTime + 30*60*60*24*program_info.eap - ? as release_date
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	release_date BETWEEN 0 AND 60*60*24*7
ORDER BY release_date;


-- name: GetMonthReleaseVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap,
	visits.StartTime + 30*60*60*24*program_info.eap - ? as release_date
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	release_date BETWEEN 0 AND 60*60*24*30
ORDER BY release_date;


-- name: GetYearReleaseVisits :many
SELECT 
	program_info.id, visits.observation, visits.visit, program_info.title, visits.Status, visits.StartTime, visits.EndTime, program_info.eap,
	visits.StartTime + 30*60*60*24*program_info.eap - ? as release_date
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
	AND
	release_date BETWEEN 0 AND 60*60*24*365
ORDER BY release_date;



