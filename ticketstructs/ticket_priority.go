package ticketstructs

import (
	"slices"
	"strings"
)

type Priority struct{}

func (p *Priority) GetPriorities() []string {
	return []string{
		"critical",
		"high",
		"medium",
		"low",
	}
}

func (p *Priority) IsValid(priority string) bool {
	return slices.Contains(p.GetPriorities(), strings.ToLower(priority))
}
