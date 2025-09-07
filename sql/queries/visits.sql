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
