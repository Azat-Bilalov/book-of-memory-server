package ds

type LoginRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

type LoginResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Role        string `json:"role"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	User_id     string `json:"user_id"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}

type RegisterResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Role        string `json:"role"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	User_id     string `json:"user_id"`
}
