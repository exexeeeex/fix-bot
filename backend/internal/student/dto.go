package student

type createStudentRequest struct {
	Name     string `json:"name"`
	GroupId  string `json:"group_id"`
	IsActive bool   `json:"is_active"`
}
