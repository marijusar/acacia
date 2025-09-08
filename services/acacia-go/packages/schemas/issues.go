package schemas

type CreateIssueInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type UpdateIssueInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}