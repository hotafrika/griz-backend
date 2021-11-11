package resources

// SelfUserResponse ...
type SelfUserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
}
