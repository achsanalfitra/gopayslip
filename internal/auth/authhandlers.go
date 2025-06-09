package auth

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"user"`
// 	Password string `json:"pass"`
// 	Role     string `json:"role"`
// }

type AuthHandler struct {
	AuthService AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: svc,
	}
}
