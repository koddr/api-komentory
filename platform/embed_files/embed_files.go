package embed_files

import _ "embed"

// SQLQueryGetManyProjects string with query for getting all (many) projects.
//go:embed sql_queries/project_getMany.sql
var SQLQueryGetManyProjects string

// SQLQueryGetManyProjectsByUserID string with query for getting all (many) projects by user ID.
//go:embed sql_queries/project_getManyByUserID.sql
var SQLQueryGetManyProjectsByUserID string

// SQLQueryGetOneProjectByID string with query for getting one project by ID.
//go:embed sql_queries/project_getOneByID.sql
var SQLQueryGetOneProjectByID string
