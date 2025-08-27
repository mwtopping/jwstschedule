-- name: CreateProgram :one
INSERT INTO program_info (id, created_at, updated_at, title, pi, eap, primetime, paralleltime, instrumentmode, programtype)
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
	?
)
RETURNING *;

-- name: GetProgramIDs :many
SELECT 
	DISTINCT id 
FROM 
	program_info;
