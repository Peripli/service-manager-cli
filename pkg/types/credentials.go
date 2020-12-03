package types

// Credentials contains types of credentials
type Credentials struct {
	Basic Basic `json:"basic,omitempty" yaml:"basic,omitempty"`
	Oauth Oauth `json:"oauth,omitempty" yaml:"oauth,omitempty"`
}

// Basic wraps basic credentials
type Basic struct {
	User     string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

// Oauth credentials
type Oauth struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}
