package types

import "github.com/Peripli/service-manager/pkg/query"

// LabelChanges wraps multiple labels change request body structure
type LabelChanges struct {
	LabelChanges []*query.LabelChange `json:"labels,omitempty"`
}
