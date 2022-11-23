package types

// Credentials contains types of credentials
type Credentials struct {
	Basic *Basic `json:"basic,omitempty" yaml:"basic,omitempty"`
	TLS   *TLS   `json:"tls,omitempty" yaml:"tls,omitempty"`
}

// Basic wraps basic credentials
type Basic struct {
	User     string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type TLS struct {
	Certificate           string `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
	Key                   string `json:"client_key,omitempty" yaml:"client_key,omitempty"`
	SMProvidedCredentials bool   `json:"sm_provided_tls_credentials" yaml:"sm_provided_tls_credentials"`
}
