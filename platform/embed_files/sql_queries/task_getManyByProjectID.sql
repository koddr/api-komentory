--
-- Query to get all (many) tasks by project ID.
-- Show only tasks with task_status == 1 (active).
-- Function signature:
--  func (q *TaskQueries) GetTasksByProjectID(project_id uuid.UUID) ([]models.GetTasks, int, error)
-- 

SELECT
	t.id,
	t.created_at,
	t.updated_at,
	t.task_attrs,
	COUNT(a.id) AS answers_count
FROM
	tasks AS t
	LEFT JOIN answers AS a ON a.task_id = t.id
WHERE
	t.project_id = $1::uuid
	AND t.task_status = 1
GROUP BY
	t.id
ORDER BY
	t.created_at DESC