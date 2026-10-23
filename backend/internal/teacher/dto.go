package teacher

type createTeacherRequest struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization_id"`
}