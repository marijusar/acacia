package schemas

type CreateIssueInput struct {
	Name                 string  `json:"name" validate:"required"`
	Description          *string `json:"description" validate:"required"`
	DescriptionSerialized *string `json:"description_serialized"`
	ColumnId             int64   `json:"column_id" validate:"required"`
}

type UpdateIssueInput struct {
	ID                    int64   `json:"id" validate:"required"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	DescriptionSerialized *string `json:"description_serialized"`
	ColumnId              int64   `json:"column_id"`
}

type ReassignIssuesInput struct {
	SourceColumnId int64 `json:"source_column" validate:"required"`
	TargetColumnId int64 `json:"target_column" validate:"required"`
}

type ReassignIssueColumn struct {
	IssueId        int64 `json:"issue_id" validate:"required"`
	TargetColumnId int64 `json:"target_column" validate:"required"`
}
