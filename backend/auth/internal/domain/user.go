package domain

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RoleInfo struct {
	Role string `json:"role"`
}

type RegisterInfo struct {
	Status string `json:"status"`
}
