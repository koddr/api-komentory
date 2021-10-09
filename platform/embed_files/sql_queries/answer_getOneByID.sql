--
-- Query to get one answer by ID.
-- Function signature:
--  func (q *AnswerQueries) GetAnswerByID(answer_id uuid.UUID) (models.Answer, int, error)
-- 

SELECT
	a.id,
	a.created_at,
	a.updated_at,
	a.project_id,
	a.task_id,
	a.answer_status,
	a.answer_attrs,
	jsonb_build_object(
		'user_id', u.id,
		'first_name', u.user_attrs->'first_name',
		'last_name', u.user_attrs->'last_name',
		'picture', u.user_attrs->'picture'
	) AS author
FROM
	answers AS a
	LEFT JOIN users AS u ON u.id = a.user_id
WHERE
	a.id = $1::uuid
GROUP BY
	a.id,
	u.id
LIMIT 1