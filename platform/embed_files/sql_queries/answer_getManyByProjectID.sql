--
-- Query to get all (many) answers by project ID.
-- Show only answers with answer_status == 1 (active).
-- Function signature:
--  func (q *AnswerQueries) GetAnswersByProjectID(project_id uuid.UUID) ([]models.GetAnswers, int, error)
-- 

SELECT
	a.id,
	a.created_at,
	a.updated_at,
	a.answer_attrs,
	jsonb_build_object(
		'user_id', u.id,
		'first_name', u.user_attrs->'first_name',
		'last_name', u.user_attrs->'last_name',
		'picture', u.user_attrs->'picture',
		'abilities', u.user_attrs->'abilities'
	) AS author
FROM
	answers AS a
	LEFT JOIN users AS u ON a.user_id = u.id
WHERE
	a.project_id = $1::uuid
	AND a.answer_status = 1
ORDER BY
	a.created_at DESC