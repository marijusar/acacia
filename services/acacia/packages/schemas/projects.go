package schemas

import "acacia/packages/db"

type CreateProjectInput struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateProjectInput struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type GetProjectDetailsResponse struct {
	db.Project
	Columns []db.ProjectStatusColumn `json:"columns"`
	Issues  []db.Issue               `json:"issues"`
}
