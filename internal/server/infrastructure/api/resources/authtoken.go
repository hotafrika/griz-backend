package resources

import "errors"

// AuthTokenRequest ...
type AuthTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ...
func (tr AuthTokenRequest) Validate() error {
	if tr.Username == "" || tr.Password == "" {
		return errors.New("params missing")
	}
	return nil
}

// AuthTokenResponse ...
type AuthTokenResponse struct {
	Token string `json:"token"`
}
