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
	program_info.id, visits.observation, visits.visit, visits.StartTime, visits.EndTime
FROM
	visits
	JOIN
		program_info
	ON visits.program_ID = program_info.id
WHERE
	visits.StartTime > 0
ORDER BY visits.StartTime;
