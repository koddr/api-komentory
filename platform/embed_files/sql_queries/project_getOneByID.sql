--
-- Query to get one project by ID.
-- Function signature:
--  func (q *ProjectQueries) GetProjectByID(project_id uuid.UUID) (models.GetProject, int, error)
-- 

SELECT
	p.id,
	p.created_at,
	p.updated_at,
	p.project_status,
	p.project_attrs,
	jsonb_build_object(
		'user_id', u.id,
		'first_name', u.user_attrs->'first_name',
		'last_name', u.user_attrs->'last_name',
		'picture', u.user_attrs->'picture'
	) AS author,
	COUNT(t.id) AS tasks_count,
	COALESCE(
		jsonb_agg(
			jsonb_build_object(
				'id', t.id,
				'status', t.task_status,
				'name', t.task_attrs->'name',
				'description', t.task_attrs->'description',
				'steps_count', jsonb_array_length(t.task_attrs->'steps')
			) 
		)
		FILTER (WHERE t.project_id IS NOT NULL), '[]'
	) AS tasks
FROM
	projects AS p
	LEFT JOIN users AS u ON u.id = p.user_id
	LEFT JOIN tasks AS t ON t.project_id = p.id
WHERE
	p.id = $1::uuid
GROUP BY
	p.id,
	u.id
LIMIT 1