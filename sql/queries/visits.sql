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
