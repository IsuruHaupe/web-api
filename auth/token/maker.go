package auth

import "time"

// This is a general interface for token based authentification/authorization.
type Maker interface {
	// CreateToken creates a new token for a specific user and duration.
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
