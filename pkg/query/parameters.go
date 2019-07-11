package query

import (
	"net/url"
	"strings"

	smquery "github.com/Peripli/service-manager/pkg/query"
)

// Parameters holds common query parameters
type Parameters struct {
	FieldQuery    []string
	LabelQuery    []string
	GeneralParams []string
}

// Encode encodes the parameters as URL query parameters
func (p *Parameters) Encode() string {
	v := url.Values{}

	if len(p.FieldQuery) > 0 {
		v.Set(string(smquery.FieldQuery), strings.Join(p.FieldQuery, "|"))
	}

	if len(p.LabelQuery) > 0 {
		v.Set(string(smquery.LabelQuery), strings.Join(p.LabelQuery, "|"))
	}

	for _, param := range p.GeneralParams {
		s := strings.Split(param, "=")
		v.Add(s[0], s[1])
	}

	return v.Encode()
}
