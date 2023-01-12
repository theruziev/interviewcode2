package auth

import "github.com/golang-jwt/jwt/v4"

type Scope string

const (
	TwoFACheckScope Scope = "2fa-check"
	UserScope       Scope = "user"
)

type Claim struct {
	jwt.RegisteredClaims
	PublicID string  `json:"pid,omitempty"`
	Email    string  `json:"email,omitempty"`
	Scopes   []Scope `json:"scp,omitempty"`
}

func (c *Claim) CheckScope(s Scope) bool {
	for _, scp := range c.Scopes {
		if s == scp {
			return true
		}
	}

	return false
}

type OtpGenerated struct {
	URL    string
	Secret string
}
