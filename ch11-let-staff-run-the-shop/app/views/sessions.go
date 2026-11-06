package views

// CredentialsRequest is the body of POST /users and POST /sessions.
type CredentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SessionResponse is a signed-in customer's token.
type SessionResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
