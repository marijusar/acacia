package schemas

type CreateIssueInput struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"required"`
	ColumnId    int64   `json:"column_id" validate:"required"`
}

type UpdateIssueInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type ReassignIssuesInput struct {
	SourceColumnId int64 `json:"source_column" validate:"required"`
	TargetColumnId int64 `json:"target_column" validate:"required"`
}

