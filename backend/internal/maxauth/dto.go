package maxauth

type MaxUser struct {
	Id           int     `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Username     string  `json:"username"`
	LanguageCode *string `json:"language_code"`
	PhotoURL     *string `json:"photo_url"`
}

type AuthResponse struct {
	Success bool    `json:"success"`
	User    MaxUser `json:"user"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type AuthRequest struct {
	InitData string `json:"initData"`
}
