--
-- Query to get one task by ID.
-- Function signature:
--  func (q *TaskQueries) GetTaskByID(task_id uuid.UUID) (models.GetTask, int, error)
-- 

SELECT
	t.*,
	COUNT(a.id) AS answers_count
FROM
	tasks AS t
	LEFT JOIN answers AS a ON a.task_id = t.id
WHERE
	t.id = $1::uuid
GROUP BY
	t.id
LIMIT 1