package schemas

type CreateProjectStatusColumnInput struct {
	ProjectID     int32  `json:"project_id" validate:"required"`
	Name          string `json:"name" validate:"required,min=1,max=255"`
	PositionIndex int16  `json:"position_index" validate:"min=0"`
}

type UpdateProjectStatusColumnInput struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	PositionIndex int16  `json:"position_index" validate:"min=0"`
}

