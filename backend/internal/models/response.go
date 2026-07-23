package models

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   *ErrorBody  `json:"error"`
	Meta    *MetaBody   `json:"meta"`
}

type ErrorBody struct {
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type MetaBody struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Version   string `json:"version"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	SchoolID string   `json:"school_id"`
}
