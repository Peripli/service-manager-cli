package query

// GeneralParameter is the name of the parameter that is used to make general key-value query pairs
const GeneralParameter = "param"

// Parameters are the parameters used to make requests to Service Manager
type Parameters map[string]*[]string

// Add adds a new parameters
func (p Parameters) Add(key, value string) {
	val, exists := p[key]
	if !exists {
		p[key] = &[]string{value}
		return
	}
	*val = append(*val, value)
}

// Get returns the values for the provided parameter key. If no such exists, returns an empty slice
func (p Parameters) Get(key string) *[]string {
	val, exists := p[key]
	if !exists {
		p[key] = &[]string{}
		return p[key]
	}
	return val
}

// Copy returns a read-only copy of the parameters
func (p Parameters) Copy() map[string][]string {
	cpy := make(map[string][]string)
	for k, v := range p {
		cpy[k] = *v
	}
	return cpy
}
