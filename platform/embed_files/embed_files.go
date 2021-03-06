package embed_files

import _ "embed"

var (
	// SQLQueryGetManyProjects string with query for getting all (many) projects.
	//go:embed sql_queries/project_getMany.sql
	SQLQueryGetManyProjects string

	// SQLQueryGetManyProjectsByUserID string with query for getting all (many) projects by user ID.
	//go:embed sql_queries/project_getManyByUserID.sql
	SQLQueryGetManyProjectsByUserID string

	// SQLQueryGetOneProjectByID string with query for getting one project by ID.
	//go:embed sql_queries/project_getOneByID.sql
	SQLQueryGetOneProjectByID string

	// SQLQueryGetOneTaskByID string with query for getting one task by ID.
	//go:embed sql_queries/task_getOneByID.sql
	SQLQueryGetOneTaskByID string

	// SQLQueryGetManyProjectsByUserID string with query for getting all (many) tasks by project ID.
	//go:embed sql_queries/task_getManyByProjectID.sql
	SQLQueryGetManyTasksByProjectID string

	// SQLQueryGetOneAnswerByID string with query for getting one answer by ID.
	//go:embed sql_queries/answer_getOneByID.sql
	SQLQueryGetOneAnswerByID string

	// SQLQueryGetManyAnswersByTaskID string with query for getting all (many) answers by task ID.
	//go:embed sql_queries/answer_getManyByTaskID.sql
	SQLQueryGetManyAnswersByTaskID string

	// SQLQueryGetManyAnswersByProjectID string with query for getting all (many) answers by project ID.
	//go:embed sql_queries/answer_getManyByProjectID.sql
	SQLQueryGetManyAnswersByProjectID string
)
