package types

type Credentials struct {
	Basic Basic `json:"basic,omitempty" yaml:"basic,omitempty"`
}

type Basic struct {
	User     string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}
