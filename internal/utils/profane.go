package utils

import (
	"slices"
	"strings"
)

type Profanity struct {
	Word     []string
	Replacer string
}

// Replace the profanity in the chirp with the replacer
func (p *Profanity) RemoveProfanity(chirp string) string {
	splited := strings.Split(chirp, " ")

	for i, word := range splited {
		if slices.Contains(p.Word, strings.ToLower(word)) {
			splited[i] = p.Replacer
		}
	}

	return strings.Join(splited, " ")
}
