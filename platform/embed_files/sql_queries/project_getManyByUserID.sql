--
-- Query to get all (many) projects by user ID.
-- Show only projects with project_status == 1 (active).
-- Function signature:
--  func (q *ProjectQueries) GetProjectsByUserID(user_id uuid.UUID) ([]models.GetProjects, int, error)
-- 

SELECT
	p.id,
	p.created_at,
	p.updated_at,
	p.project_attrs,
	jsonb_build_object(
		'user_id', u.id,
		'first_name', u.user_attrs->'first_name',
		'last_name', u.user_attrs->'last_name',
		'picture', u.user_attrs->'picture'
	) AS author,
	COUNT(t.id) AS tasks_count
FROM
	projects AS p
	LEFT JOIN users AS u ON u.id = p.user_id
	LEFT JOIN tasks AS t ON t.project_id = p.id
WHERE
	u.id = $1::uuid
	AND p.project_status = 1
GROUP BY
	p.id,
	u.id
ORDER BY
	p.created_at DESC